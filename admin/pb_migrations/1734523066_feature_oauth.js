// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/**
 * @typedef {import('../../webapp/src/modules/pocketbase/types').FeaturesRecord} Feature
 */

const FEATURES_COLLECTION_NAME = "features";

/**
 * @type {Feature}
 */
const feature = {
    name: "oauth",
    active: false,
};

migrate(
    (db) => {
        const dao = new Dao(db);
        const collection = dao.findCollectionByNameOrId(
            FEATURES_COLLECTION_NAME
        );

        /** @type {models.Record} */
        let record;
        try {
            record = dao.findFirstRecordByData(
                FEATURES_COLLECTION_NAME,
                "name",
                feature.name
            );
        } catch (e) {
            record = new Record(collection, {
                name: feature.name,
            });
        }
        record.set("envVariables", feature.envVariables);
        record.set("active", feature.active ?? true);
        dao.saveRecord(record);
    },
    (db) => {
        const dao = new Dao(db);

        try {
            const record = dao.findFirstRecordByData(
                FEATURES_COLLECTION_NAME,
                "name",
                feature.name
            );

            dao.deleteRecord(record);
        } catch {}
    }
);
