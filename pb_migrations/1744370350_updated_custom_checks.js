// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_1108732172")

  // add field
  collection.fields.addAt(3, new Field({
    "hidden": false,
    "id": "select2044034195",
    "maxSelect": 1,
    "name": "tooling",
    "presentable": false,
    "required": true,
    "system": false,
    "type": "select",
    "values": [
      "stepci"
    ]
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_1108732172")

  // remove field
  collection.fields.removeById("select2044034195")

  return app.save(collection)
})
