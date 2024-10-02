// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("pgsh9x4x20kdgjd")

  collection.indexes = [
    "CREATE INDEX `idx_tuBFjhq` ON `orgRoles` (`name`)"
  ]

  // remove
  collection.schema.removeField("iz3yszre")

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("pgsh9x4x20kdgjd")

  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_IXVrGSd` ON `orgRoles` (\n  `name`,\n  `label`\n)"
  ]

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "iz3yszre",
    "name": "label",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  return dao.saveCollection(collection)
})
