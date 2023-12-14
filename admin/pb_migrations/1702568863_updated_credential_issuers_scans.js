/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("oclpukflylnz4y7")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "ojaexdhy",
    "name": "valid",
    "type": "bool",
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
  collection.schema.removeField("ojaexdhy")

  return dao.saveCollection(collection)
})
