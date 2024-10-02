// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate(
    (db) => {
        const dao = new Dao(db);

        const usersPublicKeysCollection =
            dao.findCollectionByNameOrId("users_public_keys");

        const users = dao.findRecordsByFilter(
            "users",
            `email != ""`,
            "-created",
            0,
            0
        );

        for (const user of users) {
            if (user) {
                const userPublicKeys = new Record(usersPublicKeysCollection, {
                    owner: user.getId(),
                    bitcoin_public_key: user.get("bitcoin_public_key"),
                    ecdh_public_key: user.get("ecdh_public_key"),
                    eddsa_public_key: user.get("eddsa_public_key"),
                    ethereum_address: user.get("ethereum_address"),
                    reflow_public_key: user.get("reflow_public_key"),
                    es256_public_key: user.get("es256_public_key"),
                });

                dao.saveRecord(userPublicKeys);
            }
        }
    },
    (db) => {
        const dao = new Dao(db);

        const usersPublicKeys = dao.findRecordsByFilter(
            "users_public_keys",
            `owner != ""`,
            "-created",
            0,
            0
        );

        for (const keys of usersPublicKeys) {
            dao.deleteRecord(keys);
        }
    }
);
