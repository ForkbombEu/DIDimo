/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("k1vlx34o1x8tzno")

  collection.createRule = "@collection.organizationAuthorizations.user.id ?= @request.auth.id &&\n@collection.organizationAuthorizations.organization.id ?= organization.id &&\n@collection.organizationAuthorizations.role.name ?= \"admin\""
  collection.updateRule = "@collection.organizationAuthorizations.user.id ?= @request.auth.id &&\n@collection.organizationAuthorizations.organization.id ?= organization.id &&\n@collection.organizationAuthorizations.role.name ?= \"admin\""
  collection.deleteRule = "@collection.organizationAuthorizations.user.id ?= @request.auth.id &&\n@collection.organizationAuthorizations.organization.id ?= organization.id &&\n@collection.organizationAuthorizations.role.name ?= \"admin\""

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("k1vlx34o1x8tzno")

  collection.createRule = ""
  collection.updateRule = ""
  collection.deleteRule = ""

  return dao.saveCollection(collection)
})
