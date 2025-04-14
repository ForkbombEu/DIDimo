// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// @ts-check

/** @typedef {import('./utils.js')} Utils */

/** @param {core.RequestEvent} e */
function auditLogger(e) {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    /** @type {unknown[]} */
    const args = ["actorIp", e.realIP()];

    if (utils.isAdminContext(e)) {
        args.push("actorId", "ADMIN");
    } else {
        args.push(
            "actorId",
            e.auth?.id,
            "actorCollection",
            e.auth?.collection().name
        );
    }

    return $app.logger().with("audit", true, ...args);
}

module.exports = auditLogger;
