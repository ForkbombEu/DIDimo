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

routerAdd("POST", "/organizations/verify-user-membership", (c) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);
    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    const userId = utils.getUserFromContext(c)?.getId();

    /** @type {string | undefined} */
    const organizationId = $apis.requestInfo(c).data["organizationId"];
    if (!organizationId)
        throw utils.createMissingDataError("organizationId", "roles");

    try {
        $app.dao().findFirstRecordByFilter(
            "orgAuthorizations",
            `organization="${organizationId}" && user="${userId}"`
        );
        return c.json(200, { isMember: true });
    } catch {
        auditLogger(c).info(
            "request_from_user_not_member",
            "organizationId",
            organizationId
        );
        return c.json(200, { isMember: false });
    }
});

routerAdd("POST", "/organizations/verify-user-role", (c) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const userId = utils.getUserFromContext(c)?.getId();

    /** @type {{organizationId: string, roles: string[]}}*/
    // @ts-ignore
    const { organizationId, roles } = $apis.requestInfo(c).data;
    if (!organizationId || !roles || roles.length === 0)
        throw utils.createMissingDataError("organizationId", "roles");

    const roleFilter = `( ${roles
        .map((r) => `role.name="${r}"`)
        .join(" || ")} )`;

    try {
        $app.dao().findFirstRecordByFilter(
            "orgAuthorizations",
            `organization="${organizationId}" && user="${userId}" && ${roleFilter}`
        );
        return c.json(200, { hasRole: true });
    } catch {
        return c.json(200, { hasRole: false });
    }
});

/* Business logic hooks */

// On Organization Create – Creating owner authorization

onRecordAfterCreateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);
    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    // Don't create orgAuthorization if organization is created from admin panel
    if (utils.isAdminContext(e.httpContext)) return;

    const userId = utils.getUserFromContext(e.httpContext)?.getId();
    const organizationId = e.record?.getId();

    const ownerRole = utils.getRoleByName("owner");
    const ownerRoleId = ownerRole?.getId();

    const collection = $app.dao().findCollectionByNameOrId("orgAuthorizations");
    const record = new Record(collection, {
        organization: organizationId,
        role: ownerRoleId,
        user: userId,
    });
    $app.dao().saveRecord(record);

    //

    auditLogger(e.httpContext).info(
        "Created owner role for organization",
        "organizationId",
        e.record?.getId(),
        "organizationName",
        e.record?.get("name"),
        "userId",
        userId
    );
}, "organizations");

/* Audit hooks */

//  Log organization creation

onRecordAfterCreateRequest((e) => {
    /** @type {AuditLogger} */
    const auditLogger = require(`${__hooks}/auditLogger.js`);

    auditLogger(e.httpContext).info(
        "Created organization",
        "organizationId",
        e.record?.getId(),
        "organizationName",
        e.record?.get("name")
    );
}, "organizations");

/* Email hooks */

// Send email after organization creation

onRecordAfterCreateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (!e.record) throw utils.createMissingDataError("organization");

    const user = utils.getUserFromContext(e.httpContext);
    if (!user) throw utils.createMissingDataError("user creating organization");

    const userAddress = utils.getUserEmailAddressData(user);
    const organizationName = e.record.get("name");

    const emailData = utils.renderEmail("new-organization", {
        OrganizationName: organizationName,
        UserName: user.get("name") ?? "User",
        DashboardLink: utils.getOrganizationPageUrl(e.record.getId()),
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
