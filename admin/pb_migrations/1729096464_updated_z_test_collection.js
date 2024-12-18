/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "wlqgr8z3",
    "name": "select_multi_field",
    "type": "select",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSelect": 2,
      "values": [
        "x",
        "y",
        "z"
      ]
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "wlqgr8z3",
    "name": "select_multi_field",
    "type": "select",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSelect": 3,
      "values": [
        "x",
        "y",
        "z"
      ]
    }
  }))

  return dao.saveCollection(collection)
})
