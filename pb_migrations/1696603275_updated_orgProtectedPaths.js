/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("2y35btjamfzcr9o")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "4q4vbdck",
    "name": "roles",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "pgsh9x4x20kdgjd",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": null,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("2y35btjamfzcr9o")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "4q4vbdck",
    "name": "field",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "pgsh9x4x20kdgjd",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": null,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
})
