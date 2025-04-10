/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_120182150")

  // update field
  collection.fields.addAt(8, new Field({
    "hidden": false,
    "id": "json3351773640",
    "maxSize": 0,
    "name": "conformance_checks",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "json"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_120182150")

  // update field
  collection.fields.addAt(8, new Field({
    "hidden": false,
    "id": "json3351773640",
    "maxSize": 0,
    "name": "conformace_checks",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "json"
  }))

  return app.save(collection)
})
