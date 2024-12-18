/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "g7w0g7iqidynhim",
    "created": "2023-10-25 10:19:41.008Z",
    "updated": "2023-10-25 10:19:41.008Z",
    "name": "orgJoinRequest",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "kj3nhvxx",
        "name": "user",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "_pb_users_auth_",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "dfgtbvkp",
        "name": "organization",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "aako88kt3br4npt",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "l5dq93vm",
        "name": "status",
        "type": "select",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "pending",
            "accepted",
            "rejected"
          ]
        }
      }
    ],
    "indexes": [
      "CREATE UNIQUE INDEX `idx_H6SZMtG` ON `orgJoinRequest` (\n  `user`,\n  `organization`\n)"
    ],
    "listRule": "user.id = @request.auth.id",
    "viewRule": "",
    "createRule": "",
    "updateRule": "",
    "deleteRule": "",
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("g7w0g7iqidynhim");

  return dao.deleteCollection(collection);
})
