// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/** @typedef {import('./utils.js')} Utils */

onRecordViewRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (!utils.isAdminContext(e)) {
        e.record?.set("envVariables", null);
    }

    e.next();
}, "features");

onRecordsListRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (!utils.isAdminContext(e)) {
        e.records.forEach((r) => {
            r?.set("envVariables", null);
        });
    }

    e.next();
}, "features");
