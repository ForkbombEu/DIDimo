/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_2784107861")

  // update field
  collection.fields.addAt(5, new Field({
    "cascadeDelete": false,
    "collectionId": "pbc_2923319880",
    "hidden": false,
    "id": "relation518144557",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "test_suite",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "relation"
  }))

  // update field
  collection.fields.addAt(6, new Field({
    "cascadeDelete": false,
    "collectionId": "pbc_863811952",
    "hidden": false,
    "id": "relation2462348188",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "provider",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "relation"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_2784107861")

  // update field
  collection.fields.addAt(5, new Field({
    "cascadeDelete": false,
    "collectionId": "pbc_2923319880",
    "hidden": false,
    "id": "relation518144557",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "test_suite",
    "presentable": false,
    "required": true,
    "system": false,
    "type": "relation"
  }))

  // update field
  collection.fields.addAt(6, new Field({
    "cascadeDelete": false,
    "collectionId": "pbc_863811952",
    "hidden": false,
    "id": "relation2462348188",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "provider",
    "presentable": false,
    "required": true,
    "system": false,
    "type": "relation"
  }))

  return app.save(collection)
})
