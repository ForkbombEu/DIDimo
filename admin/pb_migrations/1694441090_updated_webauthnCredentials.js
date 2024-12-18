migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nopzrf0n7mbfu58")

  collection.updateRule = "user.id = @request.auth.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nopzrf0n7mbfu58")

  collection.updateRule = null

  return dao.saveCollection(collection)
})
