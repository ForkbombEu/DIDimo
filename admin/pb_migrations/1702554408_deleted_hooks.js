// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("3fhw2mfr9zrgodj");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "3fhw2mfr9zrgodj",
    "created": "2022-10-03 21:50:44.238Z",
    "updated": "2023-12-13 11:29:08.739Z",
    "name": "hooks",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "j8mewfur",
        "name": "collection",
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
        "id": "4xcxcfuv",
        "name": "event",
        "type": "select",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "insert",
            "update",
            "delete"
          ]
        }
      },
      {
        "system": false,
        "id": "u3bmgjpb",
        "name": "action_type",
        "type": "select",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "command",
            "post",
            "sendmail",
            "restroom-mw"
          ]
        }
      },
      {
        "system": false,
        "id": "kayyu1l3",
        "name": "action",
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
        "id": "zkengev8",
        "name": "action_params",
        "type": "text",
        "required": false,
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
        "id": "balsaeka",
        "name": "expands",
        "type": "text",
        "required": false,
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
        "id": "emgxgcok",
        "name": "disabled",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      }
    ],
    "indexes": [
      "CREATE INDEX `_3fhw2mfr9zrgodj_created_idx` ON `hooks` (`created`)"
    ],
    "listRule": null,
    "viewRule": null,
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})
