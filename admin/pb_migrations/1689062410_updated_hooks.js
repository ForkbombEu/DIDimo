// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("3fhw2mfr9zrgodj")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "u3bmgjpb",
    "name": "action_type",
    "type": "select",
    "required": true,
    "unique": false,
    "options": {
      "maxSelect": 1,
      "values": [
        "command",
        "post",
        "sendmail"
      ]
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("3fhw2mfr9zrgodj")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "u3bmgjpb",
    "name": "action_type",
    "type": "select",
    "required": true,
    "unique": false,
    "options": {
      "maxSelect": 1,
      "values": [
        "command",
        "post"
      ]
    }
  }))

  return dao.saveCollection(collection)
})
