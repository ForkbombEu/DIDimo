/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "zxnrlx8sdr589dx",
    "created": "2024-12-23 10:53:02.535Z",
    "updated": "2024-12-23 10:53:02.535Z",
    "name": "services",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "u5njqdn2",
        "name": "name",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": 3,
          "max": null,
          "pattern": ""
        }
      },
      {
        "system": false,
        "id": "cv21hi1u",
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
        "id": "kju0juul",
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
    "indexes": [],
    "listRule": "",
    "viewRule": "",
    "createRule": "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= organization.id &&\n(@collection.orgAuthorizations.role.name ?= \"admin\" || \n @collection.orgAuthorizations.role.name ?= \"owner\") ",
    "updateRule": "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= organization.id &&\n(@collection.orgAuthorizations.role.name ?= \"admin\" || \n @collection.orgAuthorizations.role.name ?= \"owner\") ",
    "deleteRule": "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= organization.id &&\n(@collection.orgAuthorizations.role.name ?= \"admin\" || \n @collection.orgAuthorizations.role.name ?= \"owner\") ",
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("zxnrlx8sdr589dx");

  return dao.deleteCollection(collection);
})
