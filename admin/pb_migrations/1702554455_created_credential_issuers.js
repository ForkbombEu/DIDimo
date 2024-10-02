// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "d86z46atrkozyt8",
    "created": "2023-12-14 11:47:35.363Z",
    "updated": "2023-12-14 11:47:35.363Z",
    "name": "credential_issuers",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "f4fce31w",
        "name": "url",
        "type": "url",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "exceptDomains": [],
          "onlyDomains": []
        }
      }
    ],
    "indexes": [],
    "listRule": null,
    "viewRule": null,
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("d86z46atrkozyt8");

  return dao.deleteCollection(collection);
})
