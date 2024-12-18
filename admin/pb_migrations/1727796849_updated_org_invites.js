/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "6stxg7b7",
    "name": "failed_email_send",
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
  collection.schema.removeField("6stxg7b7")

  return dao.saveCollection(collection)
})
