// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("aako88kt3br4npt")

  // remove field
  collection.fields.removeById("text4106974746")

  // remove field
  collection.fields.removeById("select1400097126")

  // remove field
  collection.fields.removeById("text2111657159")

  // remove field
  collection.fields.removeById("relation3725765462")

  // remove field
  collection.fields.removeById("relation694952315")

  // remove field
  collection.fields.removeById("text3401084027")

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("aako88kt3br4npt")

  // add field
  collection.fields.addAt(4, new Field({
    "autogeneratePattern": "",
    "hidden": false,
    "id": "text4106974746",
    "max": 0,
    "min": 0,
    "name": "external_website_url",
    "pattern": "",
    "presentable": false,
    "primaryKey": false,
    "required": false,
    "system": false,
    "type": "text"
  }))

  // add field
  collection.fields.addAt(5, new Field({
    "hidden": false,
    "id": "select1400097126",
    "maxSelect": 1,
    "name": "country",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "select",
    "values": [
      "AL",
      "AD",
      "AM",
      "AT",
      "AZ",
      "BY",
      "BE",
      "BA",
      "BG",
      "HR",
      "CY",
      "CZ",
      "DK",
      "EE",
      "FI",
      "FR",
      "GE",
      "DE",
      "GR",
      "HU",
      "IS",
      "IE",
      "IT",
      "KZ",
      "XK",
      "LV",
      "LI",
      "LT",
      "LU",
      "MT",
      "MD",
      "MC",
      "ME",
      "NL",
      "MK",
      "NO",
      "PL",
      "PT",
      "RO",
      "RU",
      "SM",
      "RS",
      "SK",
      "SI",
      "ES",
      "SE",
      "CH",
      "TR",
      "UA",
      "GB",
      "VA",
      "Other"
    ]
  }))

  // add field
  collection.fields.addAt(6, new Field({
    "autogeneratePattern": "",
    "hidden": false,
    "id": "text2111657159",
    "max": 0,
    "min": 0,
    "name": "documentation_url",
    "pattern": "",
    "presentable": false,
    "primaryKey": false,
    "required": false,
    "system": false,
    "type": "text"
  }))

  // add field
  collection.fields.addAt(7, new Field({
    "cascadeDelete": false,
    "collectionId": "_pb_users_auth_",
    "hidden": false,
    "id": "relation3725765462",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "created_by",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "relation"
  }))

  // add field
  collection.fields.addAt(8, new Field({
    "cascadeDelete": false,
    "collectionId": "_pb_users_auth_",
    "hidden": false,
    "id": "relation694952315",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "claimed_by",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "relation"
  }))

  // add field
  collection.fields.addAt(9, new Field({
    "autogeneratePattern": "",
    "hidden": false,
    "id": "text3401084027",
    "max": 0,
    "min": 0,
    "name": "contact_email",
    "pattern": "",
    "presentable": false,
    "primaryKey": false,
    "required": false,
    "system": false,
    "type": "text"
  }))

  return app.save(collection)
})
