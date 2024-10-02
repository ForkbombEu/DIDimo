// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("oclpukflylnz4y7")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "veypc5zf",
    "name": "error",
    "type": "select",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSelect": 1,
      "values": [
        "FILE_NOT_FOUND",
        "CONNECTION_ERROR",
        "VALIDATION_ERROR"
      ]
    }
  }))

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "znwpduqx",
    "name": "validation_result",
    "type": "json",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("oclpukflylnz4y7")

  // remove
  collection.schema.removeField("veypc5zf")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "znwpduqx",
    "name": "result",
    "type": "json",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
})
