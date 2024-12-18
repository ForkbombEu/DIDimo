migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nopzrf0n7mbfu58")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "gfynehdb",
    "name": "description",
    "type": "text",
    "required": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nopzrf0n7mbfu58")

  // remove
  collection.schema.removeField("gfynehdb")

  return dao.saveCollection(collection)
})
