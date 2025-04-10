// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate(
    (app) => {
        const collection = app.findCollectionByNameOrId("pbc_987692768");

        // add field
        collection.fields.addAt(
            7,
            new Field({
                autogeneratePattern: "",
                hidden: false,
                id: "text1874629670",
                max: 0,
                min: 0,
                name: "tags",
                pattern: "#(\\w+)",
                presentable: false,
                primaryKey: false,
                required: false,
                system: false,
                type: "text",
            })
        );

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

        // remove field
        collection.fields.removeById("text1874629670");

        return app.save(collection);
    }
);
