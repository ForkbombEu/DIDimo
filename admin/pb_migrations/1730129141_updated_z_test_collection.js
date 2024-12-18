/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "moejq2ra",
    "name": "text_with_regex",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": "^\\w+$"
    }
  }))

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "dyxsckke",
    "name": "text_field",
    "type": "text",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 2,
      "max": 10,
      "pattern": ""
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia")

  // remove
  collection.schema.removeField("moejq2ra")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "dyxsckke",
    "name": "text_field",
    "type": "text",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 2,
      "max": 10,
      "pattern": "^\\w+$"
    }
  }))

  return dao.saveCollection(collection)
})
