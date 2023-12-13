migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("sijd38964ht83q5")

  collection.listRule = "owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id = id)"
  collection.viewRule = "owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id = id)"
  collection.updateRule = "owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id = id)"
  collection.deleteRule = "owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id = id)"

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "3l6g88rp",
    "name": "owner",
    "type": "relation",
    "required": false,
    "unique": false,
    "options": {
      "collectionId": "_pb_users_auth_",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": []
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("sijd38964ht83q5")

  collection.listRule = null
  collection.viewRule = null
  collection.updateRule = null
  collection.deleteRule = null

  // remove
  collection.schema.removeField("3l6g88rp")

  return dao.saveCollection(collection)
})
