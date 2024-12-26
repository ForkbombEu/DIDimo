/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("aako88kt3br4npt")

  collection.updateRule = "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= id &&\n@collection.orgAuthorizations.role.name ?= \"owner\""
  collection.deleteRule = "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= id &&\n@collection.orgAuthorizations.role.name ?= \"owner\""

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("aako88kt3br4npt")

  collection.updateRule = "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= id &&\n@collection.orgAuthorizations.role.name ?= \"admin\""
  collection.deleteRule = "@collection.orgAuthorizations.user.id ?= @request.auth.id &&\n@collection.orgAuthorizations.organization.id ?= id &&\n@collection.orgAuthorizations.role.name ?= \"admin\""

  return dao.saveCollection(collection)
})
