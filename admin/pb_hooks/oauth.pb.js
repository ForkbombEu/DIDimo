// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/** @typedef {import('./utils.js')} Utils */

/* Updating user info on first register */
onRecordAfterAuthWithOAuth2Request((e) => {
    if (!e.isNewRecord) return;

    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const { record: user, oAuth2User } = e;
    if (!user || !oAuth2User)
        throw utils.createMissingDataError("user", "oAuth2User");

    user.set("name", oAuth2User.name);
    user.markAsNotNew();
    $app.dao().saveRecord(user);
}, "users");
