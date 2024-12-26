/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia")

  collection.listRule = "owner = null ||\n( owner != null && (\n    owner.id = @request.auth.id ||\n    ( @collection.authorizations.record_id = id &&\n      @collection.authorizations.users.id ?= @request.auth.id &&\n      @collection.authorizations.owner.id ?= owner.id )\n  )\n)"
  collection.viewRule = "owner = null ||\n( owner != null && (\n    owner.id = @request.auth.id ||\n    ( @collection.authorizations.record_id = id &&\n      @collection.authorizations.users.id ?= @request.auth.id &&\n      @collection.authorizations.owner.id ?= owner.id )\n  )\n)"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia")

  collection.listRule = "(owner != null && owner.id = @request.auth.id) || owner = null"
  collection.viewRule = "(owner != null && owner.id = @request.auth.id) || owner = null"

  return dao.saveCollection(collection)
})
