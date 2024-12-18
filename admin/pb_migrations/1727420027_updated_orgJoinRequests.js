/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("g7w0g7iqidynhim")

  collection.listRule = "user.id = @request.auth.id ||\n( @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id &&\n  ( @collection.orgAuthorizations.role.name ?= \"admin\" ||\n    @collection.orgAuthorizations.role.name ?= \"owner\" )\n)"
  collection.viewRule = "user.id = @request.auth.id ||\n( @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id &&\n  ( @collection.orgAuthorizations.role.name ?= \"admin\" ||\n    @collection.orgAuthorizations.role.name ?= \"owner\" )\n)"
  collection.updateRule = "@request.data.organization:isset = false &&\n@request.data.user:isset = false &&\n@request.data.reminders:isset = false &&\n( @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id &&\n  ( @collection.orgAuthorizations.role.name ?= \"admin\" ||\n    @collection.orgAuthorizations.role.name ?= \"owner\" )\n)"
  collection.deleteRule = "user.id = @request.auth.id ||\n( @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id &&\n  ( @collection.orgAuthorizations.role.name ?= \"admin\" ||\n    @collection.orgAuthorizations.role.name ?= \"owner\" )\n)"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("g7w0g7iqidynhim")

  collection.listRule = "user.id = @request.auth.id ||\n(\n  @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id\n  &&\n  (\n    @collection.orgAuthorizations.role.name ?= \"admin\" ||\n    @collection.orgAuthorizations.role.name ?= \"owner\"\n  )\n)"
  collection.viewRule = "user.id = @request.auth.id ||\n(\n  @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id\n  &&\n  (\n    @collection.orgAuthorizations.role.name ?= \"admin\" ||\n    @collection.orgAuthorizations.role.name ?= \"owner\"\n  )\n)"
  collection.updateRule = "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= organization.id &&\n(\n  @collection.orgAuthorizations.role.name ?= \"admin\" ||\n  @collection.orgAuthorizations.role.name ?= \"owner\"\n)"
  collection.deleteRule = "user.id = @request.auth.id ||\n(\n  @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id\n  &&\n  (\n    @collection.orgAuthorizations.role.name ?= \"admin\" ||\n    @collection.orgAuthorizations.role.name ?= \"owner\"\n  )\n)"

  return dao.saveCollection(collection)
})
