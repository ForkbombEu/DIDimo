/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_4151506089")

  // add field
  collection.fields.addAt(7, new Field({
    "hidden": false,
    "id": "json1338962426",
    "maxSize": 0,
    "name": "i18n_label",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "json"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_4151506089")

  // remove field
  collection.fields.removeById("json1338962426")

  return app.save(collection)
})
