/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  collection.name = "org_invites"
  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_xyC4d4D` ON `org_invites` (`user_email`)",
    "CREATE UNIQUE INDEX `idx_cdEx89V` ON `org_invites` (`user`)"
  ]

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("v12xstc3xygpawh")

  collection.name = "org_invite"
  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_xyC4d4D` ON `org_invite` (`user_email`)",
    "CREATE UNIQUE INDEX `idx_cdEx89V` ON `org_invite` (`user`)"
  ]

  return dao.saveCollection(collection)
})
