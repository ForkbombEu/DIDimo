// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "oclpukflylnz4y7",
    "created": "2023-12-14 11:51:08.241Z",
    "updated": "2023-12-14 11:51:08.241Z",
    "name": "credential_issuers_scans",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "jd8vlqxf",
        "name": "credential_issuer",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "d86z46atrkozyt8",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "znwpduqx",
        "name": "result",
        "type": "json",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {}
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
  const collection = dao.findCollectionByNameOrId("oclpukflylnz4y7");

  return dao.deleteCollection(collection);
})
