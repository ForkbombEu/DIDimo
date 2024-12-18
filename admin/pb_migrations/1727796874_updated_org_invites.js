/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_xyC4d4D` ON `org_invites` (\n  `user_email`,\n  `organization`\n)"
  ]

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_xyC4d4D` ON `org_invites` (`user_email`)"
  ]

  return dao.saveCollection(collection)
})
