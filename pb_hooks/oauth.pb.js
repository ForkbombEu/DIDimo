// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/** @typedef {import('./utils.js')} Utils */

/* Updating user info on first register */

onRecordAuthWithOAuth2Request((e) => {
    e.next();

    if (e.isNewRecord) {
        /** @type {Utils} */
        const utils = require(`${__hooks}/utils.js`);

        const { record: user, oAuth2User } = e;
        if (!user || !oAuth2User)
            throw utils.createMissingDataError("user", "oAuth2User");

        user.set("name", oAuth2User.name);
        user.markAsNotNew();
        $app.Save(user);
    }
}, "users");
