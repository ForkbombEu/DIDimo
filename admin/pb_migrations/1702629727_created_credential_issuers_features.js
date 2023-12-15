/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "51eupx0oybwd9ce",
    "created": "2023-12-15 08:42:07.752Z",
    "updated": "2023-12-15 08:42:07.752Z",
    "name": "credential_issuers_features",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "mjadbjkc",
        "name": "type",
        "type": "select",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "FILE_EXISTS",
            "VALID_JSON",
            "SCHEMA_COMPLIANT"
          ]
        }
      },
      {
        "system": false,
        "id": "vpteavtu",
        "name": "metadata",
        "type": "json",
        "required": false,
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
  const collection = dao.findCollectionByNameOrId("51eupx0oybwd9ce");

  return dao.deleteCollection(collection);
})
