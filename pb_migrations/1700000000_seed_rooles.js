/// <reference path="../pb_data/types.d.ts" />

const roles = ["owner", "admin", "member"];

migrate((db) => {
    const dao = new Dao(db);
    const rolesCollection = dao.findCollectionByNameOrId("orgRoles");

    for (const role of roles) {
        const record = new Record(rolesCollection, { name: role });
        dao.saveRecord(record);
    }

    // -- Seeding protected paths -- //

    const ownerRole = dao.findFirstRecordByData("orgRoles", "name", "owner");
    // const adminRole = dao.findFirstRecordByData("orgRoles", "name", "admin");

    const pathsCollection = dao.findCollectionByNameOrId("orgProtectedPaths");
    dao.saveRecord(
        new Record(pathsCollection, {
            pathRegex: "/organizations/(.*)/settings",
            roles: [ownerRole.id],
        })
    );
});
