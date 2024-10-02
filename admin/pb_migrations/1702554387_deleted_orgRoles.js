// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("pgsh9x4x20kdgjd");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "pgsh9x4x20kdgjd",
    "created": "2023-09-22 08:25:04.234Z",
    "updated": "2023-12-13 11:29:08.767Z",
    "name": "orgRoles",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "ewo9sxda",
        "name": "name",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": "^[a-z_]+$"
        }
      }
    ],
    "indexes": [
      "CREATE INDEX `idx_tuBFjhq` ON `orgRoles` (`name`)"
    ],
    "listRule": "",
    "viewRule": "",
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})
