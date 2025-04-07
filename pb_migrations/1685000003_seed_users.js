// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// @ts-check
/// <reference path="../pb_data/types.d.ts" />

const USERS_COLLECTION_NAME = "users";

/**
 * @param {string} letter
 */
function createSampleUserData(letter) {
    const name = `user${letter}`;
    return {
        email: `user${letter}@example.org`,
        password: `user${letter}user${letter}`,
        username: name,
        name: name,
    };
}

/**
 *
 * @param {excludeHooks<core.App>} app
 * @param {string} letter
 * @returns
 */
function addSampleUser(app, letter) {
    const { email, password, name } = createSampleUserData(letter);
    const collection = app.findCollectionByNameOrId(USERS_COLLECTION_NAME);

    const record = new Record(collection);
    record.setEmail(email);
    record.setPassword(password);
    record.setVerified(true);
    record.set("name", name);
    record.set("emailVisibility", true);

    app.save(record);
}

migrate((app) => {
    addSampleUser(app, "A");
    addSampleUser(app, "B");
    addSampleUser(app, "C");
});
