// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/** @typedef {import('./utils.js')} Utils */
/** @typedef {import('./auditLogger.js')} AuditLogger */

/**
 * INDEX
 * - Routes
 * - Business logic hooks
 * - Audit hooks
 * - Email hooks
 */

/* Routes */

routerAdd("POST", "/organizations/verify-user-membership", (e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);
    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    const userId = utils.getUserFromContext(e)?.id;

    /** @type {string | undefined} */
    const organizationId = e.requestInfo().body["organizationId"];
    if (!organizationId)
        throw utils.createMissingDataError("organizationId", "roles");

    try {
        $app.findFirstRecordByFilter(
            "orgAuthorizations",
            `organization="${organizationId}" && user="${userId}"`
        );
        return e.json(200, { isMember: true });
    } catch {
        auditLogger(e).info(
            "request_from_user_not_member",
            "organizationId",
            organizationId
        );
        return e.json(200, { isMember: false });
    }
});

routerAdd("POST", "/organizations/verify-user-role", (e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const userId = utils.getUserFromContext(e)?.id;

    /** @type {{organizationId: string, roles: string[]}}*/
    // @ts-ignore
    const { organizationId, roles } = e.requestInfo().body;
    if (!organizationId || !roles || roles.length === 0)
        throw utils.createMissingDataError("organizationId", "roles");

    const roleFilter = `( ${roles
        .map((r) => `role.name="${r}"`)
        .join(" || ")} )`;

    try {
        $app.findFirstRecordByFilter(
            "orgAuthorizations",
            `organization="${organizationId}" && user="${userId}" && ${roleFilter}`
        );
        return e.json(200, { hasRole: true });
    } catch {
        return e.json(200, { hasRole: false });
    }
});

/* Business logic hooks */

// On Organization Create – Creating owner authorization

onRecordCreateRequest((e) => {
    e.next();

    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);
    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    // Don't create orgAuthorization if organization is created from admin panel
    if (utils.isAdminContext(e)) return;

    const userId = utils.getUserFromContext(e)?.id;
    const organizationId = e.record?.id;

    const ownerRole = utils.getRoleByName("owner");
    const ownerRoleId = ownerRole?.id;

    const collection = $app.findCollectionByNameOrId("orgAuthorizations");
    const record = new Record(collection, {
        organization: organizationId,
        role: ownerRoleId,
        user: userId,
    });
    $app.Save(record);

    //

    auditLogger(e).info(
        "Created owner role for organization",
        "organizationId",
        e.record?.id,
        "organizationName",
        e.record?.get("name"),
        "userId",
        userId
    );
}, "organizations");

/* Audit hooks */

//  Log organization creation

onRecordCreateRequest((e) => {
    e.next();

    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    auditLogger(e).info(
        "Created organization",
        "organizationId",
        e.record?.id,
        "organizationName",
        e.record?.get("name")
    );
}, "organizations");

/* Email hooks */

// Send email after organization creation

onRecordCreateRequest((e) => {
    e.next();

    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (!e.record) throw utils.createMissingDataError("organization");

    const user = utils.getUserFromContext(e);
    if (!user) throw utils.createMissingDataError("user creating organization");

    const userAddress = utils.getUserEmailAddressData(user);
    const organizationName = e.record.get("name");

    const emailData = utils.renderEmail("new-organization", {
        OrganizationName: organizationName,
        UserName: user.get("name") ?? "User",
        DashboardLink: utils.getOrganizationPageUrl(e.record.id),
        AppName: utils.getAppName(),
    });

    const res = utils.sendEmail({
        to: userAddress,
        ...emailData,
    });
    if (res instanceof Error) {
        console.error(res);
    }
}, "organizations");
