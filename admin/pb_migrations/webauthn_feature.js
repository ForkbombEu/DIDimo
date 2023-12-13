const FEATURES_COLLECTION_NAME = "features";

/**
    type Feature = {
        name: string;
        envVariables: Record<string, unknown>;
    }
*/
const features = [
    {
        name: "webauthn",
        envVariables: {
            DISPLAY_NAME: "didimo",
            RPID: "localhost",
            RPORIGINS: "http://localhost:5173",
        },
    },
];

migrate(
    (db) => {
        const dao = new Dao(db);
        const collection = dao.findCollectionByNameOrId(
            FEATURES_COLLECTION_NAME
        );

        for (const feature of features) {
            const record = new Record(collection, {
                name: feature.name,
                envVariables: feature.envVariables,
                active: true,
            });

            dao.saveRecord(record);
        }
    },
    (db) => {
        const dao = new Dao(db);

        // delete the previously created record (if exists)
        for (const feature of features) {
            try {
                const record = dao.findFirstRecordByData(
                    FEATURES_COLLECTION_NAME,
                    "name",
                    feature.name
                );

                dao.deleteRecord(record);
            } catch {}
        }
    }
);
