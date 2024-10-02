// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("285guwyxvr46lsu")

  collection.updateRule = "@request.auth.id = owner.id"
  collection.deleteRule = "@request.auth.id = owner.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("285guwyxvr46lsu")

  collection.updateRule = null
  collection.deleteRule = null

  return dao.saveCollection(collection)
})
