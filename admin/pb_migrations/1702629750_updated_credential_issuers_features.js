// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("51eupx0oybwd9ce")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "82akkqw2",
    "name": "report",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "oclpukflylnz4y7",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("51eupx0oybwd9ce")

  // remove
  collection.schema.removeField("82akkqw2")

  return dao.saveCollection(collection)
})
