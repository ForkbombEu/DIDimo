migrate((db) => {
  const collection = new Collection({
    "id": "55roakey1q8vk77",
    "created": "2023-08-16 14:08:01.186Z",
    "updated": "2023-08-16 14:08:01.186Z",
    "name": "sessionDataWebauthn",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "xdfemwi3",
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
        "id": "kds5s6oc",
        "name": "session",
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

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("55roakey1q8vk77");

  return dao.deleteCollection(collection);
})
