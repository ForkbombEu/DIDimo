// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("aako88kt3br4npt")

  collection.updateRule = "@collection.organizationAuthorizations.user.id ?= @request.auth.id &&\n@collection.organizationAuthorizations.organization.id ?= id &&\n@collection.organizationAuthorizations.role.name ?= \"admin\""
  collection.deleteRule = "@collection.organizationAuthorizations.user.id ?= @request.auth.id &&\n@collection.organizationAuthorizations.organization.id ?= id &&\n@collection.organizationAuthorizations.role.name ?= \"admin\""

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("aako88kt3br4npt")

  collection.updateRule = null
  collection.deleteRule = null

  return dao.saveCollection(collection)
})
