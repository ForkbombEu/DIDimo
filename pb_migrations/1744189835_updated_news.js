// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate(
    (app) => {
        const collection = app.findCollectionByNameOrId("pbc_987692768");

        // update collection data
        unmarshal(
            {
                listRule: "published=true",
                viewRule: "published=true",
            },
            collection
        );

        return app.save(collection);
    },
    (app) => {
        const collection = app.findCollectionByNameOrId("pbc_987692768");

        // update collection data
        unmarshal(
            {
                listRule: "",
                viewRule: null,
            },
            collection
        );

        return app.save(collection);
    }
);
