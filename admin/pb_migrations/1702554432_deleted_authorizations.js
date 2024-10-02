// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("285guwyxvr46lsu");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "285guwyxvr46lsu",
    "created": "2023-07-11 13:23:15.614Z",
    "updated": "2023-12-13 11:29:08.737Z",
    "name": "authorizations",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "5jd5bhu8",
        "name": "owner",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "_pb_users_auth_",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": []
        }
      },
      {
        "system": false,
        "id": "w4xrqdgs",
        "name": "users",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "_pb_users_auth_",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": null,
          "displayFields": []
        }
      },
      {
        "system": false,
        "id": "g1t9kpqo",
        "name": "collection_id",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
        }
      },
      {
        "system": false,
        "id": "fqvzbaze",
        "name": "record_id",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
        }
      }
    ],
    "indexes": [
      "CREATE UNIQUE INDEX `idx_w4uoK0u` ON `authorizations` (\n  `owner`,\n  `collection_id`,\n  `record_id`\n)"
    ],
    "listRule": "@request.auth.id = owner.id || users.id ?= @request.auth.id",
    "viewRule": "@request.auth.id = owner.id || users.id ?= @request.auth.id",
    "createRule": "@request.auth.id != ''",
    "updateRule": "@request.auth.id = owner.id",
    "deleteRule": "@request.auth.id = owner.id",
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})
