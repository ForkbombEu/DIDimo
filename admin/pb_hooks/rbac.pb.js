// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
/**
 * @typedef {import('./utils.js')} Utils
 */

onRecordAfterCreateRequest((e) => {
    console.log("Hook - Creating owner role for new organization");

    // Don't create auth if organization is created from admin panel
    if ($apis.requestInfo(e.httpContext).admin) return;

    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const userId = utils.getUserFromContext(e.httpContext).id;
    const organizationId = e.record.id;

    const ownerRole = utils.getOwnerRole();
    const ownerRoleId = ownerRole.id;

    const collection = $app.dao().findCollectionByNameOrId("orgAuthorizations");
    const record = new Record(collection, {
        organization: organizationId,
        role: ownerRoleId,
        user: userId,
    });
    $app.dao().saveRecord(record);
}, "organizations");

onRecordBeforeDeleteRequest((e) => {
    console.log("Hook - Checking if deleting owner role is possible");

    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const organizationId = e.record.get("organization");
    const roleId = e.record.get("role");
    const ownerRoleId = utils.getOwnerRole().id;

    if (roleId !== ownerRoleId) return;

    const adminAuthorizations = $app
        .dao()
        .findRecordsByFilter(
            "orgAuthorizations",
            `organization="${organizationId}" && role="${ownerRoleId}"`
        );

    if (adminAuthorizations.length > 1) return;

    throw new Error("Can't remove the last owner role!");
}, "orgAuthorizations");

routerAdd("POST", "/verify-org-authorization", (c) => {
    console.log("Route - Checking if user has the correct org authorization");

    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const userId = utils.getUserFromContext(c).id;
    if (!userId) throw new Error("User must be logged!");

    /**  @type {organizationId: string, url: string}*/
    const { organizationId, url } = $apis.requestInfo(c).data;
    if (!organizationId || !url) throw new Error("Missing data in request");

    const userAuthorization = $app
        .dao()
        .findFirstRecordByFilter(
            "orgAuthorizations",
            `organization="${organizationId}" && user="${userId}"`
        );
    // Here we assume that there is only one role for each organization
    // Also enforced by API rules
    if (!userAuthorization) throw new Error("Not authorized");
    const userRole = userAuthorization.get("role");

    const protectedPaths = $app
        .dao()
        .findRecordsByFilter("orgProtectedPaths", "pathRegex != ''");

    const matchingPaths = protectedPaths.filter((p) => {
        const regex = new RegExp(p.get("pathRegex"));
        return regex.test(url);
    });

    const isAllowed = matchingPaths
        .map((p) => {
            /** @type {string[]} */
            const roles = p.get("roles");
            return roles;
        })
        .every((roles) => roles.includes(userRole));

    if (!isAllowed) throw new Error("Not authorized");
});

routerAdd("POST", "/verify-user-role", (c) => {
    console.log("Route - Checking if user has the required role");

    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const userId = utils.getUserFromContext(c).id;
    if (!userId) throw new Error("User must be logged!");

    /**  @type {organizationId: string, roles: string[]}*/
    const { organizationId, roles } = $apis.requestInfo(c).data;
    if (!organizationId || !roles || roles.length === 0)
        throw new Error("Missing data in request");

    const roleFilter = `( ${roles
        .map((r) => `role.name="${r}"`)
        .join(" || ")} )`;

    const userAuthorization = $app
        .dao()
        .findFirstRecordByFilter(
            "orgAuthorizations",
            `organization="${organizationId}" && user="${userId}" && ${roleFilter}`
        );
    // Here we assume that there is only one role for each organization
    // Also enforced by API rules
    if (!userAuthorization) throw new Error("Not authorized");
});
