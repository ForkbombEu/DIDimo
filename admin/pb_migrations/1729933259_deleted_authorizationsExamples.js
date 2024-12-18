/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("sijd38964ht83q5");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "sijd38964ht83q5",
    "created": "2023-07-11 13:26:18.946Z",
    "updated": "2024-10-24 07:03:42.098Z",
    "name": "authorizationsExamples",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "s6z3yyrm",
        "name": "name",
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
        "id": "3l6g88rp",
        "name": "owner",
        "type": "relation",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "_pb_users_auth_",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": []
        }
      }
    ],
    "indexes": [],
    "listRule": "owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id ?= id)",
    "viewRule": "owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id ?= id)",
    "createRule": "@request.auth.id != ''",
    "updateRule": "owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id ?= id)",
    "deleteRule": "owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id ?= id)",
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})
