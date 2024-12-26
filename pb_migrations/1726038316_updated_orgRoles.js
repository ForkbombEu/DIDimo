/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("pgsh9x4x20kdgjd")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "xcrt6mly",
    "name": "level",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("pgsh9x4x20kdgjd")

  // remove
  collection.schema.removeField("xcrt6mly")

  return dao.saveCollection(collection)
})
