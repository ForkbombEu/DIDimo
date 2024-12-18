/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia")

  // remove
  collection.schema.removeField("ckgk3hqn")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "xkpmxcbr",
    "name": "relation_field",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "ktjgpqf146ss2ia",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("ktjgpqf146ss2ia")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "ckgk3hqn",
    "name": "relation_field",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "nopzrf0n7mbfu58",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  // remove
  collection.schema.removeField("xkpmxcbr")

  return dao.saveCollection(collection)
})
