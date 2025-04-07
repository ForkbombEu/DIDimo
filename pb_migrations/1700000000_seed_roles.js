// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// @ts-check

/// <reference path="../pb_data/types.d.ts" />

/** @type {Array<{name:string, level:number}>} */
const roles = [
    { name: "owner", level: 0 },
    { name: "admin", level: 1 },
    { name: "member", level: 9 },
];

//

migrate((app) => {
    const rolesCollection = app.findCollectionByNameOrId("orgRoles");

    roles
        .map((role) => new Record(rolesCollection, role))
        .forEach((roleRecord) => app.save(roleRecord));
});
