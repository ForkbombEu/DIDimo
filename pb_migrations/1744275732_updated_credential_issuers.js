/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_678514665")

  // update collection data
  unmarshal({
    "createRule": null,
    "deleteRule": "owner.id = @request.auth.id",
    "listRule": "published = true || owner.id = @request.auth.id",
    "updateRule": "@request.body.owner:isset = false &&\n@request.body.url:isset = false",
    "viewRule": "published = true || owner.id = @request.auth.id"
  }, collection)

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_678514665")

  // update collection data
  unmarshal({
    "createRule": "",
    "deleteRule": "",
    "listRule": "",
    "updateRule": "",
    "viewRule": ""
  }, collection)

  return app.save(collection)
})
