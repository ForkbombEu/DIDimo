/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("k1vlx34o1x8tzno")

  collection.name = "orgAuthorizations"
  collection.listRule = ""
  collection.viewRule = ""
  collection.createRule = ""
  collection.updateRule = ""
  collection.deleteRule = ""
  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_ik2InMV` ON `orgAuthorizations` (\n  `organization`,\n  `user`\n)"
  ]

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("k1vlx34o1x8tzno")

  collection.name = "organizationAuthorizations"
  collection.listRule = "user.id = @request.auth.id || organization.owners.id ~ @request.auth.id"
  collection.viewRule = "user.id = @request.auth.id || organization.owners.id ~ @request.auth.id"
  collection.createRule = "@collection.organizationAuthorizations.user.id ?= @request.auth.id &&\n@collection.organizationAuthorizations.organization.id ?= organization.id &&\n@collection.organizationAuthorizations.role.name ?= \"admin\""
  collection.updateRule = "@collection.organizationAuthorizations.user.id ?= @request.auth.id &&\n@collection.organizationAuthorizations.organization.id ?= organization.id &&\n@collection.organizationAuthorizations.role.name ?= \"admin\""
  collection.deleteRule = "@collection.organizationAuthorizations.user.id ?= @request.auth.id &&\n@collection.organizationAuthorizations.organization.id ?= organization.id &&\n@collection.organizationAuthorizations.role.name ?= \"admin\""
  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_ik2InMV` ON `organizationAuthorizations` (\n  `organization`,\n  `user`\n)"
  ]

  return dao.saveCollection(collection)
})
