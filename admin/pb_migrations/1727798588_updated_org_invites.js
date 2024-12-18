/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  collection.listRule = "user.id = @request.auth.id ||\n( @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id &&\n  ( @collection.orgAuthorizations.role.name ?= \"admin\" || \n   @collection.orgAuthorizations.role.name ?= \"owner\") )"
  collection.viewRule = "user.id = @request.auth.id ||\n( @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id &&\n  ( @collection.orgAuthorizations.role.name ?= \"admin\" || \n   @collection.orgAuthorizations.role.name ?= \"owner\") )"
  collection.createRule = "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= organization.id &&\n(@collection.orgAuthorizations.role.name ?= \"admin\" || \n @collection.orgAuthorizations.role.name ?= \"owner\") "
  collection.deleteRule = "user.id = @request.auth.id ||\n( @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id &&\n  ( @collection.orgAuthorizations.role.name ?= \"admin\" || \n   @collection.orgAuthorizations.role.name ?= \"owner\") )"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  collection.listRule = ""
  collection.viewRule = ""
  collection.createRule = "@request.auth.id = @collection.orgAuthorizations.user.id &&\norganization.id = @collection.orgAuthorizations.organization.id &&\n( @collection.orgAuthorizations.role.name = 'admin' ||\n  @collection.orgAuthorizations.role.name = 'owner' )"
  collection.deleteRule = "user.id = @request.auth.id ||\n( @request.auth.id = @collection.orgAuthorizations.user.id &&\n  organization.id = @collection.orgAuthorizations.organization.id &&\n  ( @collection.orgAuthorizations.role.name = 'admin' ||\n    @collection.orgAuthorizations.role.name = 'owner' )\n)"

  return dao.saveCollection(collection)
})
