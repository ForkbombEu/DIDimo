/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("g7w0g7iqidynhim")

  collection.deleteRule = "user.id = @request.auth.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("g7w0g7iqidynhim")

  collection.deleteRule = "user.id = @request.auth.id ||\n( @collection.orgAuthorizations.user.id ?= @request.auth.id &&\n  @collection.orgAuthorizations.organization.id ?= organization.id &&\n  ( @collection.orgAuthorizations.role.name ?= \"admin\" ||\n    @collection.orgAuthorizations.role.name ?= \"owner\" )\n)"

  return dao.saveCollection(collection)
})
