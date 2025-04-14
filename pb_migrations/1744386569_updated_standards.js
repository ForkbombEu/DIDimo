// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate(
    (app) => {
        const collection = app.findCollectionByNameOrId("pbc_216647879");

        // update field
        collection.fields.addAt(
            8,
            new Field({
                autogeneratePattern: "",
                hidden: false,
                id: "text2104180948",
                max: 0,
                min: 0,
                name: "standard_uid",
                pattern: "",
                presentable: false,
                primaryKey: false,
                required: false,
                system: false,
                type: "text",
            })
        );

        return app.save(collection);
    },
    (app) => {
        const collection = app.findCollectionByNameOrId("pbc_216647879");

        // update field
        collection.fields.addAt(
            8,
            new Field({
                autogeneratePattern: "",
                hidden: false,
                id: "text2104180948",
                max: 0,
                min: 0,
                name: "standar_uid",
                pattern: "",
                presentable: false,
                primaryKey: false,
                required: false,
                system: false,
                type: "text",
            })
        );

        return app.save(collection);
    }
);
