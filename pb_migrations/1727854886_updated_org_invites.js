/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "gtmrpmp3",
    "name": "declined",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  // remove
  collection.schema.removeField("gtmrpmp3")

  return dao.saveCollection(collection)
})
