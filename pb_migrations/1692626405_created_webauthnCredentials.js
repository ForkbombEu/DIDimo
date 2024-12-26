migrate((db) => {
  const collection = new Collection({
    "id": "nopzrf0n7mbfu58",
    "created": "2023-08-21 14:00:05.210Z",
    "updated": "2023-08-21 14:00:05.210Z",
    "name": "webauthnCredentials",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "xootznbs",
        "name": "user",
        "type": "relation",
        "required": false,
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
        "id": "of3px3ud",
        "name": "credential",
        "type": "json",
        "required": false,
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

  collection.indexes = [
    "CREATE INDEX `idx_p4OsklQ` ON `webauthnCredentials` (`user`)"
  ]
  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("nopzrf0n7mbfu58");

  collection.indexes = []

  return dao.deleteCollection(collection);
})
