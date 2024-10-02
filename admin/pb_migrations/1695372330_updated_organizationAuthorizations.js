// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("k1vlx34o1x8tzno")

  collection.listRule = "user.id = @request.auth.id || organization.owners.id ~ @request.auth.id"
  collection.viewRule = "user.id = @request.auth.id || organization.owners.id ~ @request.auth.id"
  collection.createRule = "organization.owners.id ~ @request.auth.id"
  collection.updateRule = "organization.owners.id ~ @request.auth.id"
  collection.deleteRule = "organization.owners.id ~ @request.auth.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("k1vlx34o1x8tzno")

  collection.listRule = null
  collection.viewRule = null
  collection.createRule = null
  collection.updateRule = null
  collection.deleteRule = null

  return dao.saveCollection(collection)
})
