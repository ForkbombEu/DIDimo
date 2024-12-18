/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("g7w0g7iqidynhim")

  collection.name = "orgJoinRequests"
  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_H6SZMtG` ON `orgJoinRequests` (\n  `user`,\n  `organization`\n)"
  ]

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("g7w0g7iqidynhim")

  collection.name = "orgJoinRequest"
  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_H6SZMtG` ON `orgJoinRequest` (\n  `user`,\n  `organization`\n)"
  ]

  return dao.saveCollection(collection)
})
