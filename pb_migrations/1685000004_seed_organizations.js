// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// @ts-check
/// <reference path="../pb_data/types.d.ts" />

const ORGANIZATIONS_NAME = "organizations";

/**
 *
 * @param {excludeHooks<core.App>} app
 * @returns
 */
function addDefaultOrganization(app) {
    const name = "default";
    const collection = app.findCollectionByNameOrId(ORGANIZATIONS_NAME);
    const record = new Record(collection);
    record.set("name", name);
    record.set("description", "Credimi public organization");
    app.save(record);
}

migrate((app) => {
    addDefaultOrganization(app);
});
