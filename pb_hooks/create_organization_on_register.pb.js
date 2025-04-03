// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/** @typedef {import('./utils.js')} Utils */

onRecordCreateRequest((e) => {
    e.next();

    const userId = e.record?.id;
    if (!userId) return;

    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const organizationCollection =
        $app.findCollectionByNameOrId("organizations");

    const newOrg = new Record(organizationCollection, {
        name: userId,
    });

    $app.save(newOrg);

    utils.createOwnerRoleForOrganization(e, newOrg.id, userId);
}, "users");
