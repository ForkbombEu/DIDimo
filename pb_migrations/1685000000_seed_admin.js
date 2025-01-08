/// <reference path="../pb_data/types.d.ts" />
// @ts-check

const ADMIN_EMAIL = "admin@example.org";
const SUPERUSERS = "_superusers";

migrate(
    (app) => {
        try {
            app.findAuthRecordByEmail(SUPERUSERS, ADMIN_EMAIL);
        } catch {
            const superusers = app.findCollectionByNameOrId(SUPERUSERS);
            const admin = new Record(superusers);
            admin.setEmail(ADMIN_EMAIL);
            admin.setPassword("adminadmin");
            app.save(admin);
        }
    },
    (app) => {
        const admin = app.findAuthRecordByEmail(SUPERUSERS, ADMIN_EMAIL);
        app.delete(admin);
    }
);
