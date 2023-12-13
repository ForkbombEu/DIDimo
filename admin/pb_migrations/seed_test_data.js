const USERS_COLLECTION_NAME = "users";

const users = {
    A: {
        email: "userA@example.org",
        password: "userAuserA",
    },
    B: {
        email: "userB@example.org",
        password: "userBuserB",
    },
    C: {
        email: "userC@example.org",
        password: "userCuserC",
    },
};

function createSampleUserData(letter) {
    return {
        email: `user${letter}@example.org`,
        username: `user${letter}`,
        password: `user${letter}user${letter}`,
    };
}

function addSampleUser(dao, letter) {
    const { email, username, password } = createSampleUserData(letter);
    const collection = dao.findCollectionByNameOrId(USERS_COLLECTION_NAME);

    const record = new Record(collection);
    record.setUsername(username);
    record.setEmail(email);
    record.setPassword(password);
    record.setVerified(true);
    record.set("emailVisibility", true);

    dao.saveRecord(record);

    return record;
}

function removeUser(dao, email) {
    const record = dao.findAuthRecordByEmail(USERS_COLLECTION_NAME, email);
    dao.deleteRecord(record);
}

migrate(
    (db) => {
        const dao = new Dao(db);

        /* Users */
        const userA = addSampleUser(dao, "A");
        const userB = addSampleUser(dao, "B");
        const userC = addSampleUser(dao, "C");

        /* AuthorizationExample */
        const authorizationsExamples = dao.findCollectionByNameOrId(
            "authorizationsExamples"
        );
        const authorizationExample = new Record(authorizationsExamples);
        authorizationExample.set("name", "authorizationExample");
        authorizationExample.set("owner", userA.id);
        dao.saveRecord(authorizationExample);

        /* Authorization */
        const authorizations = dao.findCollectionByNameOrId("authorizations");
        const authorization = new Record(authorizations);
        authorization.set("owner", userA.id);
        authorization.set("users", [userB.id]);
        authorization.set("collection_id", authorizationsExamples.id);
        authorization.set("record_id", authorizationExample.id);
        dao.saveRecord(authorization);

        return;
    },
    (db) => {
        return;
    }
);
