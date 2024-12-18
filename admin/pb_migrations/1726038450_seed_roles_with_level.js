/// <reference path="../pb_data/types.d.ts" />
// @ts-check

/** @type {Array<{name:string, level:number}>} */
const roles = [
    { name: "owner", level: 0 },
    { name: "admin", level: 1 },
    { name: "member", level: 9 },
];

migrate(
    (db) => {
        const dao = new Dao(db);
        const rolesCollection = dao.findCollectionByNameOrId("orgRoles");

        roles
            // Taking into account the roles created with a previous migration
            .map((role) => {
                try {
                    const record = dao.findFirstRecordByData(
                        "orgRoles",
                        "name",
                        role.name
                    );
                    record.set("level", role.level);
                    return record;
                } catch {
                    return new Record(rolesCollection, role);
                }
            })
            .forEach((roleRecord) => dao.saveRecord(roleRecord));
    },
    (db) => {
        // add down queries...
    }
);
