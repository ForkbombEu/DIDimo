/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("oclpukflylnz4y7")

  collection.name = "credential_issuers_reports"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("oclpukflylnz4y7")

  collection.name = "reports"

  return dao.saveCollection(collection)
})
