// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/** @typedef {import('./utils.js')} Utils */
/** @typedef {import('./auditLogger.js')} AuditLogger */

/**
 * INDEX
 * - Base hooks
 * - Email hooks
 */

/* Base hooks */

onRecordCreateRequest((e) => {
    e.record?.set("status", "pending");
    e.record?.set("reminders", 0);
}, "orgJoinRequests");

// Cannot create join request if user is already member

onRecordCreateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    /** @type {string | undefined} */
    const organizationId = e.record?.get("organization");
    /** @type {string | undefined} */
    const userId = e.record?.get("user");

    const authorization = utils.findFirstRecordByFilter(
        "orgAuthorizations",
        `organization.id = "${organizationId}" && user.id = "${userId}"`
    );

    if (authorization)
        throw new BadRequestError(utils.errors.user_is_already_member);
}, "orgJoinRequests");

// Create orgAuthorization after accepting membership request

onRecordUpdateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (!e.record) throw utils.createMissingDataError("orgJoinRequest");

    const status = e.record.get("status");
    if (status != "accepted") return;

    const orgAuthorizationsCollection = $app.findCollectionByNameOrId("orgAuthorizations");
    if (!orgAuthorizationsCollection)
        throw utils.createMissingDataError("orgAuthorizationsCollection");

    const organizationId = e.record.get("organization");
    const userId = e.record.get("user");

    const memberRole = utils.getRoleByName("member");
    const roleId = memberRole?.id;

    const record = new Record(orgAuthorizationsCollection, {
        user: userId,
        organization: organizationId,
        role: roleId,
    });

    $app.saveRecord(record);
}, "orgJoinRequests");

/* Email hooks - Notifications to Admins */

onRecordCreateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (!e.record) throw utils.createMissingDataError("orgJoinRequest");

    const organization = utils.getExpanded(e.record, "organization");
    if (!organization)
        throw utils.createMissingDataError("organization of orgJoinRequest");
    const user = utils.getExpanded(e.record, "user");
    if (!user) throw utils.createMissingDataError("user of orgJoinRequest");

    const organizationId = organization.id;
    const recipients = utils.getOrganizationAdminsAddresses(organizationId);

    for (const adminAddress of recipients) {
        const email = utils.renderEmail("membership-request-new", {
            OrganizationName: organization.getString("name"),
            Admin: adminAddress.name,
            UserName: user.getString("name"),
            DashboardLink: utils.getOrganizationMembersPageUrl(organizationId),
            AppName: utils.getAppName(),
        });

        const res = utils.sendEmail({
            to: adminAddress,
            ...email,
        });

        if (res instanceof Error) {
            console.error("Email send error");
        }
    }
}, "orgJoinRequests");

cronAdd("remind admins about join requests", "0 9 * * 1", () => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const organizations = utils.findRecordsByFilter(
        "organizations",
        "id != null"
    );

    organizations
        .map((organization) => ({
            organization,
            requests: utils
                .findRecordsByFilter(
                    "orgJoinRequests",
                    `organization.id = "${organization.id}"`
                )
                .filter((r) => r.get("status") == "pending"),
        }))
        .filter(({ requests }) => requests.length > 0)
        .forEach(({ organization, requests }) => {
            const organizationId = organization.id;
            const OrganizationName = organization.get("name");

            const recipients =
                utils.getOrganizationAdminsAddresses(organizationId);

            for (const recipient of recipients) {
                const email = utils.renderEmail("membership-request-pending", {
                    OrganizationName,
                    DashboardLink:
                        utils.getOrganizationMembersPageUrl(organizationId),
                    Admin: recipient.name,
                    PendingNumber: requests.length.toString(),
                    AppName: utils.getAppName(),
                });

                const res = utils.sendEmail({
                    to: recipient,
                    ...email,
                });
                if (res instanceof Error) {
                    console.error("Email send error");
                }
            }
        });
});

/* Email hooks - Notifications to Users */

onRecordUpdateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (!e.record) throw utils.createMissingDataError("orgJoinRequest");

    /** @type {string} */
    const status = e.record.get("status");

    const isRelevantChange = status == "accepted" || status == "rejected";
    if (!isRelevantChange) return;

    const organization = utils.getExpanded(e.record, "organization");
    if (!organization)
        throw utils.createMissingDataError("organization of orgJoinRequest");
    const user = utils.getExpanded(e.record, "user");
    if (!user) throw utils.createMissingDataError("user of orgJoinRequest");

    /** @type {string} */
    const OrganizationName = organization.get("name");
    const userAddress = utils.getUserEmailAddressData(user);

    /**
     * @typedef {Object} EmailContent
     * @property {string} subject
     * @property {string} html
     */

    /** @type {EmailContent} */
    const successEmail = utils.renderEmail("membership-request-accepted", {
        OrganizationName,
        UserName: userAddress.name,
        DashboardLink: utils.getOrganizationPageUrl(organization.id),
        AppName: utils.getAppName(),
    });

    /** @type {EmailContent} */
    const rejectEmail = utils.renderEmail("membership-request-declined", {
        OrganizationName,
        UserName: userAddress.name,
        DashboardLink: utils.getAppUrl() + "/my/organizations",
        AppName: utils.getAppName(),
    });

    const emailContent = status == "accepted" ? successEmail : rejectEmail;
    utils.sendEmail({ to: [userAddress], ...emailContent });
}, "orgJoinRequests");

/* Audit logs */

onRecordCreateRequest((e) => {
    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    auditLogger(e).info(
        "Created membership request",
        "organizationId",
        e.record?.get("organization"),
        "requestId",
        e.record?.id
    );
}, "orgJoinRequests");

onRecordUpdateRequest((e) => {
    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    auditLogger(e).info(
        "Updated membership request",
        "organizationId",
        e.record?.get("organization"),
        "status",
        e.record?.get("status"),
        "requestId",
        e.record?.id
    );
}, "orgJoinRequests");

onRecordDeleteRequest((e) => {
    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    auditLogger(e).info(
        "Deleted membership request",
        "organizationId",
        e.record?.get("organization"),
        "status",
        e.record?.get("status"),
        "requestId",
        e.record?.id
    );
}, "orgJoinRequests");

/* IMPORTANT: This hook must be registered last */

onRecordUpdateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (!e.record) throw utils.createMissingDataError("orgJoinRequest");

    /** @type {string} */
    const status = e.record.get("status");

    const isRelevantChange = status == "accepted" || status == "rejected";
    if (!isRelevantChange) return;

    $app.Delete(e.record);
}, "orgJoinRequests");
