// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nopzrf0n7mbfu58")

  collection.listRule = "user.id = @request.auth.id"
  collection.viewRule = "user.id = @request.auth.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nopzrf0n7mbfu58")

  collection.listRule = null
  collection.viewRule = null

  return dao.saveCollection(collection)
})
