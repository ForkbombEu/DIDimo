/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "ktjgpqf146ss2ia",
    "created": "2024-07-19 20:54:49.363Z",
    "updated": "2024-07-19 20:54:49.363Z",
    "name": "z_test_collection",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "dyxsckke",
        "name": "text_field",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": 2,
          "max": 10,
          "pattern": "^\\w+$"
        }
      },
      {
        "system": false,
        "id": "5nakh5br",
        "name": "url_field",
        "type": "url",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "exceptDomains": [],
          "onlyDomains": [
            "uld.com"
          ]
        }
      },
      {
        "system": false,
        "id": "zjaakeak",
        "name": "number_field",
        "type": "number",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": 2,
          "max": 70,
          "noDecimal": true
        }
      },
      {
        "system": false,
        "id": "plf9bqlo",
        "name": "boolean_field",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "o9i7jdyx",
        "name": "richtext_field",
        "type": "editor",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "convertUrls": true
        }
      },
      {
        "system": false,
        "id": "euflpigq",
        "name": "email_field",
        "type": "email",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "exceptDomains": [
            "mat.cos"
          ],
          "onlyDomains": []
        }
      },
      {
        "system": false,
        "id": "r6bk98mc",
        "name": "date_field",
        "type": "date",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "min": "2024-07-18 12:00:00.000Z",
          "max": "2024-07-28 12:00:00.000Z"
        }
      },
      {
        "system": false,
        "id": "jnvj44gc",
        "name": "select_field",
        "type": "select",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "a",
            "b",
            "c"
          ]
        }
      },
      {
        "system": false,
        "id": "wlqgr8z3",
        "name": "select_multi_field",
        "type": "select",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 3,
          "values": [
            "x",
            "y",
            "z"
          ]
        }
      },
      {
        "system": false,
        "id": "ldafagmk",
        "name": "file_field",
        "type": "file",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "mimeTypes": [
            "application/zip"
          ],
          "thumbs": [],
          "maxSelect": 1,
          "maxSize": 5242880,
          "protected": true
        }
      },
      {
        "system": false,
        "id": "ctqisupe",
        "name": "file_multi_field",
        "type": "file",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "mimeTypes": [],
          "thumbs": [],
          "maxSelect": 99,
          "maxSize": 5242880,
          "protected": false
        }
      },
      {
        "system": false,
        "id": "ckgk3hqn",
        "name": "relation_field",
        "type": "relation",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "nopzrf0n7mbfu58",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "dmehmjlp",
        "name": "relation_multi_field",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "_pb_users_auth_",
          "cascadeDelete": false,
          "minSelect": 1,
          "maxSelect": 4,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "wj9awedr",
        "name": "json_field",
        "type": "json",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSize": 2000000
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
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia");

  return dao.deleteCollection(collection);
})
