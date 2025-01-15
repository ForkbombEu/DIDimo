/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_863811952")

  // update collection data
  unmarshal({
    "name": "solutions"
  }, collection)

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_863811952")

  // update collection data
  unmarshal({
    "name": "services"
  }, collection)

  return app.save(collection)
})
