/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_xyC4d4D` ON `org_invite` (`user_email`)",
    "CREATE UNIQUE INDEX `idx_cdEx89V` ON `org_invite` (`user`)"
  ]

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  collection.indexes = []

  return dao.saveCollection(collection)
})
