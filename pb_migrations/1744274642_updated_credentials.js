/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_183765882")

  // update collection data
  unmarshal({
    "listRule": "(published = true && credential_issuer.published = true) ||\n@request.auth.id = credential_issuer.owner.id",
    "updateRule": "@request.auth.id = credential_issuer.owner.id",
    "viewRule": "(published = true && credential_issuer.published = true) ||\n@request.auth.id = credential_issuer.owner.id"
  }, collection)

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_183765882")

  // update collection data
  unmarshal({
    "listRule": "",
    "updateRule": null,
    "viewRule": ""
  }, collection)

  return app.save(collection)
})
