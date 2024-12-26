/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "v12xstc3xygpawh",
    "created": "2024-09-30 18:02:08.746Z",
    "updated": "2024-09-30 18:02:08.746Z",
    "name": "org_invite",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "7yfkzy39",
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
        "id": "32axfwku",
        "name": "user",
        "type": "relation",
        "required": false,
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
        "id": "gbqnjiwb",
        "name": "user_email",
        "type": "email",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "exceptDomains": [],
          "onlyDomains": []
        }
      }
    ],
    "indexes": [],
    "listRule": "user.id = @request.auth.id ||\n( @request.auth.id = @collection.orgAuthorizations.user.id &&\n  organization.id = @collection.orgAuthorizations.organization.id &&\n  ( @collection.orgAuthorizations.role.name = 'admin' ||\n    @collection.orgAuthorizations.role.name = 'owner' )\n)",
    "viewRule": "user.id = @request.auth.id ||\n( @request.auth.id = @collection.orgAuthorizations.user.id &&\n  organization.id = @collection.orgAuthorizations.organization.id &&\n  ( @collection.orgAuthorizations.role.name = 'admin' ||\n    @collection.orgAuthorizations.role.name = 'owner' )\n)",
    "createRule": "@request.auth.id = @collection.orgAuthorizations.user.id &&\norganization.id = @collection.orgAuthorizations.organization.id &&\n( @collection.orgAuthorizations.role.name = 'admin' ||\n  @collection.orgAuthorizations.role.name = 'owner' )",
    "updateRule": null,
    "deleteRule": "user.id = @request.auth.id ||\n( @request.auth.id = @collection.orgAuthorizations.user.id &&\n  organization.id = @collection.orgAuthorizations.organization.id &&\n  ( @collection.orgAuthorizations.role.name = 'admin' ||\n    @collection.orgAuthorizations.role.name = 'owner' )\n)",
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh");

  return dao.deleteCollection(collection);
})
