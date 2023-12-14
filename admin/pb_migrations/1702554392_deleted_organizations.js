/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("aako88kt3br4npt");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "aako88kt3br4npt",
    "created": "2023-09-20 08:05:20.415Z",
    "updated": "2023-12-13 11:29:08.769Z",
    "name": "organizations",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "de5ifbee",
        "name": "name",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": 2,
          "max": null,
          "pattern": ""
        }
      },
      {
        "system": false,
        "id": "zhuxbrib",
        "name": "avatar",
        "type": "file",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "maxSize": 5242880,
          "mimeTypes": [
            "image/png",
            "image/jpeg",
            "image/webp",
            "image/svg+xml"
          ],
          "thumbs": [],
          "protected": false
        }
      },
      {
        "system": false,
        "id": "pjjpq1r4",
        "name": "description",
        "type": "editor",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "convertUrls": false
        }
      }
    ],
    "indexes": [
      "CREATE UNIQUE INDEX `idx_PHN81EZ` ON `organizations` (`name`)"
    ],
    "listRule": "",
    "viewRule": "",
    "createRule": "@request.auth.id != ''",
    "updateRule": "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= id &&\n@collection.orgAuthorizations.role.name ?= \"owner\"",
    "deleteRule": "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= id &&\n@collection.orgAuthorizations.role.name ?= \"owner\"",
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})
