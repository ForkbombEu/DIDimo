// @ts-check

/** @typedef {import('./utils.js')} Utils */

/** @param {echo.Context} c */
function auditLogger(c) {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    /** @type {unknown[]} */
    const args = ["actorIp", c.realIP()];

    if (utils.isAdminContext(c)) {
        args.push("actorId", "ADMIN");
    } else {
        args.push(
            "actorId",
            $apis.requestInfo(c).authRecord?.getId(),
            "actorCollection",
            $apis.requestInfo(c).authRecord?.collection().name
        );
    }

    return $app.logger().with("audit", true, ...args);
}

module.exports = auditLogger;
