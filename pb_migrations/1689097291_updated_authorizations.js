migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("285guwyxvr46lsu")

  collection.updateRule = "@request.auth.id = owner.id"
  collection.deleteRule = "@request.auth.id = owner.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("285guwyxvr46lsu")

  collection.updateRule = null
  collection.deleteRule = null

  return dao.saveCollection(collection)
})
