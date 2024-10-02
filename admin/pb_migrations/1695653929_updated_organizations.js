// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("aako88kt3br4npt")

  collection.updateRule = null
  collection.deleteRule = null

  // remove
  collection.schema.removeField("rdd3ph1n")

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("aako88kt3br4npt")

  collection.updateRule = "owners.id ~ @request.auth.id"
  collection.deleteRule = "owners.id ~ @request.auth.id"

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "rdd3ph1n",
    "name": "owners",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "_pb_users_auth_",
      "cascadeDelete": false,
      "minSelect": 1,
      "maxSelect": null,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
})
