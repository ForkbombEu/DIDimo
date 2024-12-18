// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/// <reference path="./ambient.d.ts" />

/** @typedef {import("../../webapp/src/modules/pocketbase/types").OrgRolesRecord} OrgRole */
/** @typedef {import("../../webapp/src/modules/pocketbase/types").OrgAuthorizationsRecord} OrgAuthorization */
/** @typedef {import("../../webapp/src/modules/pocketbase/types").UsersRecord} User */

/** @typedef {Omit<mail.Address, "string">} Address */

//

/* -- Error codes -- */

const errors = {
    not_authorized: "not_authorized",
    missing_data: "missing_data",
    user_not_logged: "user_not_logged",

    cant_create_an_authorization_for_yourself:
        "cant_create_an_authorization_for_yourself",

    cant_edit_last_owner_role: "cant_edit_last_owner_role",
    cant_delete_last_owner_role: "cant_delete_last_owner_role",

    cant_create_role_higher_than_or_equal_to_yours:
        "cant_create_role_higher_than_or_equal_to_yours",
    cant_edit_role_higher_than_or_equal_to_yours:
        "cant_edit_role_higher_than_or_equal_to_yours",
    cant_delete_role_higher_than_or_equal_to_yours:
        "cant_delete_role_higher_than_or_equal_to_yours",

    user_is_already_member: "user_is_already_member",
};

/* -- RBAC Utils -- */

/**
 * @param {string} name
 * @returns {RecordModel<OrgRole> | undefined}
 */
function getRoleByName(name) {
    try {
        return $app.dao().findFirstRecordByData("orgRoles", "name", name);
    } catch {
        return undefined;
    }
}

/**
 * @param {models.Record} role
 * @returns {number}
 */
function getRoleLevel(role) {
    return role.get("level");
}

/**
 * @param {models.Record} orgAuthorization
 */
function isLastOwnerAuthorization(orgAuthorization) {
    const organizationId = orgAuthorization.get("organization");
    const roleId = orgAuthorization.get("role");
    const ownerRoleId = getRoleByName("owner")?.getId();

    if (roleId !== ownerRoleId) return false;

    const ownerAuthorizations = findRecordsByFilter(
        "orgAuthorizations",
        `organization="${organizationId}" && role="${ownerRoleId}"`
    );

    return ownerAuthorizations.length == 1;
}

/**
 * @param {string} userId
 * @param {string} organizationId
 * @returns {RecordModel<OrgRole> | undefined}
 */
function getUserRole(userId, organizationId, dao = $app.dao()) {
    const authorization = findFirstRecordByFilter(
        "orgAuthorizations",
        `user = "${userId}" && organization = "${organizationId}"`,
        dao
    );
    if (!authorization) return undefined;
    return getExpanded(authorization, "role", dao);
}

/**
 *
 * @param {core.RecordUpdateEvent|core.RecordDeleteEvent|core.RecordViewEvent|core.RecordCreateEvent} e
 */
function getUserContextInOrgAuthorizationHookEvent(e) {
    const userId = getUserFromContext(e.httpContext)?.getId();

    /** @type {string | undefined} */
    const organizationId = e.record?.get("organization");

    if (!userId || !organizationId)
        throw createMissingDataError("requestingUserId", "organizationId");

    const isSelf = userId === e.record?.get("user");

    const userRole = getUserRole(userId, organizationId);
    if (!userRole) throw createMissingDataError("requestingUserRole");

    const userRoleLevel = getRoleLevel(userRole);

    return { userId, userRole, isSelf, userRoleLevel };
}

/* -- Pocketbase utils -- */

/**
 * @param {echo.Context} c
 * @returns {RecordModel<User> | undefined}
 */
function getUserFromContext(c) {
    return $apis.requestInfo(c).authRecord;
}

/**
 * @param {string} collection
 * @param {string} filter
 * @param {daos.Dao} [dao=$app.dao()]
 * @returns {Array<models.Record>}
 */
function findRecordsByFilter(collection, filter, dao = $app.dao()) {
    return dao
        .findRecordsByFilter(collection, filter, "", 0, 0)
        .filter((v) => v != undefined);
}

/**
 * @param {string} collection
 * @param {string} filter
 */
function findFirstRecordByFilter(collection, filter, dao = $app.dao()) {
    try {
        return dao.findFirstRecordByFilter(collection, filter);
    } catch {
        return undefined;
    }
}

/**
 * @param {models.Record} record
 * @param {string} key
 * @returns {models.Record | undefined}
 */
function getExpanded(record, key, dao = $app.dao()) {
    try {
        // @ts-ignore
        dao.expandRecord(record, [key], null);
        return record.expandedOne(key);
    } catch (e) {
        return undefined;
    }
}

/**
 * @param {echo.Context} c
 */
function isAdminContext(c) {
    return Boolean($apis.requestInfo(c).admin);
}

/**
 * @param {string[]} args
 */
function createMissingDataError(...args) {
    return new BadRequestError(errors.missing_data + ": " + args.join(", "));
}

/**
 * @param {models.Record} user
 * @returns {Address}
 */
function getUserEmailAddressData(user) {
    /** @type {string} */
    const name = user.get("name");
    /** @type {string} */
    const address = user.get("email");
    if (!name || !address)
        throw createMissingDataError("userEmail", "userName");
    return {
        name,
        address,
    };
}

/**
 * @param {{from?: Address, to: Address[] | Address, subject: string, html: string}} data
 * @return {Error | undefined}
 */
function sendEmail(data) {
    try {
        const message = new MailerMessage({
            // @ts-expect-error Missing string() object
            from: {
                address:
                    data.from?.address ?? $app.settings().meta.senderAddress,
                name: data.from?.name ?? $app.settings().meta.senderName,
            },
            // @ts-expect-error Missing string() object
            to: Array.isArray(data.to) ? data.to : [data.to],
            subject: data.subject,
            html: data.html,
        });

        $app.newMailClient().send(message);
    } catch (e) {
        return e;
    }
}

/**
 * @param {string} string
 * @returns {string}
 */
function removeTrailingSlash(string) {
    if (string.endsWith("/")) return string.slice(0, -1);
    else return string;
}

/**
 * @param {string} organizationId
 * @param {daos.Dao} [dao=$app.dao()]
 * @returns {Address[]}
 */
function getOrganizationAdminsAddresses(organizationId, dao = $app.dao()) {
    const recipients = findRecordsByFilter(
        "orgAuthorizations",
        `organization.id = "${organizationId}" && ( role.name = "admin" || role.name = "owner" )`,
        dao
    );

    return recipients
        .map((r) => getExpanded(r, "user"))
        .filter((u) => u != undefined)
        .map((u) => getUserEmailAddressData(u));
}

/** @returns {string} */
function getAppUrl() {
    return removeTrailingSlash($app.settings().meta.appUrl);
}

/** @returns {string} */
function getAppName() {
    return $app.settings().meta.appName;
}

/**
 * @param {string} organizationId
 * @returns {string}
 */
function getOrganizationPageUrl(organizationId) {
    return `${getAppUrl()}/my/organizations/${organizationId}`;
}

/**
 * @param {string} organizationId
 * @returns {string}
 */
function getOrganizationMembersPageUrl(organizationId) {
    return `${getOrganizationPageUrl(organizationId)}/members`;
}

/**
 * @param {echo.Context} c
 */
function runOrganizationInviteEndpointChecks(c) {
    /** @type {{inviteId: string | undefined}} */
    // @ts-ignore
    const data = $apis.requestInfo(c).data;
    const { inviteId } = data;
    if (!inviteId || typeof inviteId != "string")
        throw createMissingDataError("inviteId");

    const userId = getUserFromContext(c)?.getId();
    if (!userId) throw createMissingDataError("userId");

    const invite = findFirstRecordByFilter("org_invites", `id = "${inviteId}"`);
    if (!invite) throw createMissingDataError("organization invite");

    const isOwner = invite.get("user") == userId;
    if (!isOwner) throw new ForbiddenError();

    return { userId, invite, isOwner };
}

/**
 *
 * @param {core.RecordUpdateEvent} event
 * @param {string[]} fields
 */
function getRecordUpdateEventDiff(event, fields = []) {
    const updatedRecord = event.record;
    const originalRecord = event.record?.originalCopy();
    if (!updatedRecord || !originalRecord)
        throw createMissingDataError("updated record");

    if (fields.length == 0)
        fields = getCollectionFields(updatedRecord.collection());

    return fields
        .map((f) => ({
            field: f,
            newValue: updatedRecord.get(f),
            oldValue: originalRecord.get(f),
        }))
        .filter((d) => d.newValue != d.oldValue);
}

/**
 * @param {models.Collection} collection
 */
function getCollectionFields(collection) {
    return collection.schema
        .fields()
        .map((f) => f?.name)
        .filter((n) => n != undefined);
}

//

/** @type {EmailRenderer} */
const renderEmail = (name, data) => {
    const emailPath = $filepath.join(
        __hooks,
        "..",
        "pb_public",
        "emails",
        `${name}.html`
    );
    const html = $template
        .loadFiles(emailPath)
        .render(data)
        .replace("@WL", "{{Weblink}}")
        .replace("@US", "{{{unsubscribe}}}");
    const subject = html.match(/<title>(.*?)<\/title>/)?.at(1) ?? "";
    return {
        html,
        subject,
    };
};

/**
 *
 * @param {core.RecordUpdateEvent} event
 * @param {string[]} fields
 */
function getRecordUpdateEventDiff(event, fields = []) {
    const updatedRecord = event.record;
    const originalRecord = event.record?.originalCopy();
    if (!updatedRecord || !originalRecord)
        throw createMissingDataError("updated record");

    if (fields.length == 0)
        fields = getCollectionFields(updatedRecord.collection());

    return fields
        .map((f) => ({
            field: f,
            newValue: updatedRecord.get(f),
            oldValue: originalRecord.get(f),
        }))
        .filter((d) => d.newValue != d.oldValue);
}

/**
 * @param {models.Collection} collection
 */
function getCollectionFields(collection) {
    return collection.schema
        .fields()
        .map((f) => f?.name)
        .filter((n) => n != undefined);
}

/**
 *
 * @param { echo.Context } httpContext
 * @returns { RecordModel<User> | models.Admin | undefined }
 */
function getRequestAgent(httpContext) {
    /** @type {RecordModel<User> | models.Admin | undefined} */
    let agent = undefined;

    const adminContext = isAdminContext(httpContext);
    const user = getUserFromContext(httpContext);

    if (adminContext) {
        agent = $apis.requestInfo(httpContext).admin;
    } else if (user) {
        agent = user;
    }

    return agent;
}

/**
 *
 * @param { echo.Context } httpContext
 * @returns { string | undefined }
 */
function getRequestAgentName(httpContext) {
    const agent = getRequestAgent(httpContext);
    if (!agent) return undefined;

    if ("getString" in agent) {
        return agent.getString("name");
    } else {
        return `System Admin ${agent.getId()}`;
    }
}

//

module.exports = {
    getUserFromContext,
    getRoleByName,
    findRecordsByFilter,
    isLastOwnerAuthorization,
    findFirstRecordByFilter,
    getUserRole,
    getExpanded,
    isAdminContext,
    getRoleLevel,
    createMissingDataError,
    getUserContextInOrgAuthorizationHookEvent,
    getUserEmailAddressData,
    sendEmail,
    removeTrailingSlash,
    getOrganizationAdminsAddresses,
    getOrganizationMembersPageUrl,
    getAppUrl,
    getAppName,
    runOrganizationInviteEndpointChecks,
    renderEmail,
    getOrganizationPageUrl,
    getRecordUpdateEventDiff,
    getRequestAgent,
    getRequestAgentName,
    errors,
};
