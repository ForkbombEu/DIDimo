migrate((db) => {
  const snapshot = [
    {
      "id": "3fhw2mfr9zrgodj",
      "created": "2022-10-03 21:50:44.238Z",
      "updated": "2023-07-11 15:40:56.394Z",
      "name": "hooks",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "j8mewfur",
          "name": "collection",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "4xcxcfuv",
          "name": "event",
          "type": "select",
          "required": true,
          "unique": false,
          "options": {
            "maxSelect": 1,
            "values": [
              "insert",
              "update",
              "delete"
            ]
          }
        },
        {
          "system": false,
          "id": "u3bmgjpb",
          "name": "action_type",
          "type": "select",
          "required": true,
          "unique": false,
          "options": {
            "maxSelect": 1,
            "values": [
              "command",
              "post"
            ]
          }
        },
        {
          "system": false,
          "id": "kayyu1l3",
          "name": "action",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "zkengev8",
          "name": "action_params",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "balsaeka",
          "name": "expands",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "emgxgcok",
          "name": "disabled",
          "type": "bool",
          "required": false,
          "unique": false,
          "options": {}
        }
      ],
      "indexes": [
        "CREATE INDEX `_3fhw2mfr9zrgodj_created_idx` ON `hooks` (`created`)"
      ],
      "listRule": null,
      "viewRule": null,
      "createRule": null,
      "updateRule": null,
      "deleteRule": null,
      "options": {}
    },
    {
      "id": "z4cc0g76ciqx13v",
      "created": "2023-05-17 16:48:21.255Z",
      "updated": "2023-07-11 12:40:55.518Z",
      "name": "features",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "lj04bczn",
          "name": "name",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "yctafkxl",
          "name": "active",
          "type": "bool",
          "required": false,
          "unique": false,
          "options": {}
        },
        {
          "system": false,
          "id": "yo3lch0p",
          "name": "envVariables",
          "type": "json",
          "required": false,
          "unique": false,
          "options": {}
        }
      ],
      "indexes": [
        "CREATE UNIQUE INDEX `idx_T3xOIQy` ON `features` (`name`)"
      ],
      "listRule": "",
      "viewRule": "",
      "createRule": null,
      "updateRule": null,
      "deleteRule": null,
      "options": {}
    },
    {
      "id": "p9dhra382rts3bg",
      "created": "2023-06-07 10:44:19.476Z",
      "updated": "2023-07-11 13:16:04.708Z",
      "name": "crudExample",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "o0hwn3yh",
          "name": "text",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": 6,
            "max": 10,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "3hmes6kr",
          "name": "file_only_pdf_json",
          "type": "file",
          "required": true,
          "unique": false,
          "options": {
            "maxSelect": 99,
            "maxSize": 5242880,
            "mimeTypes": [
              "application/json",
              "application/pdf"
            ],
            "thumbs": [],
            "protected": false
          }
        },
        {
          "system": false,
          "id": "o72s557k",
          "name": "boolean",
          "type": "bool",
          "required": false,
          "unique": false,
          "options": {}
        },
        {
          "system": false,
          "id": "dzneawfh",
          "name": "select",
          "type": "select",
          "required": true,
          "unique": false,
          "options": {
            "maxSelect": 1,
            "values": [
              "opt1",
              "opt2",
              "opt3",
              "opt4",
              "opt5"
            ]
          }
        },
        {
          "system": false,
          "id": "t1i6phfc",
          "name": "textarea",
          "type": "editor",
          "required": false,
          "unique": false,
          "options": {}
        },
        {
          "system": false,
          "id": "lybqqrme",
          "name": "text_with_regex",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": "^w+S"
          }
        },
        {
          "system": false,
          "id": "vkmf5bzv",
          "name": "multiselect",
          "type": "select",
          "required": false,
          "unique": false,
          "options": {
            "maxSelect": 2,
            "values": [
              "A",
              "B",
              "C",
              "D"
            ]
          }
        },
        {
          "system": false,
          "id": "kx4ozipk",
          "name": "relation",
          "type": "relation",
          "required": false,
          "unique": false,
          "options": {
            "collectionId": "z4cc0g76ciqx13v",
            "cascadeDelete": false,
            "minSelect": null,
            "maxSelect": null,
            "displayFields": []
          }
        },
        {
          "system": false,
          "id": "nfghsb83",
          "name": "image",
          "type": "file",
          "required": false,
          "unique": false,
          "options": {
            "maxSelect": 1,
            "maxSize": 5242880,
            "mimeTypes": [
              "image/png",
              "image/jpeg"
            ],
            "thumbs": [],
            "protected": false
          }
        },
        {
          "system": false,
          "id": "uy849pd2",
          "name": "relation_single",
          "type": "relation",
          "required": true,
          "unique": false,
          "options": {
            "collectionId": "z4cc0g76ciqx13v",
            "cascadeDelete": false,
            "minSelect": null,
            "maxSelect": 1,
            "displayFields": []
          }
        },
        {
          "system": false,
          "id": "61lxlzl2",
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
        }
      ],
      "indexes": [],
      "listRule": "",
      "viewRule": "",
      "createRule": "",
      "updateRule": "",
      "deleteRule": "",
      "options": {}
    },
    {
      "id": "285guwyxvr46lsu",
      "created": "2023-07-11 13:23:15.614Z",
      "updated": "2023-07-11 15:42:24.490Z",
      "name": "authorizations",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "5jd5bhu8",
          "name": "owner",
          "type": "relation",
          "required": true,
          "unique": false,
          "options": {
            "collectionId": "_pb_users_auth_",
            "cascadeDelete": false,
            "minSelect": null,
            "maxSelect": 1,
            "displayFields": []
          }
        },
        {
          "system": false,
          "id": "w4xrqdgs",
          "name": "users",
          "type": "relation",
          "required": true,
          "unique": false,
          "options": {
            "collectionId": "_pb_users_auth_",
            "cascadeDelete": false,
            "minSelect": null,
            "maxSelect": null,
            "displayFields": []
          }
        },
        {
          "system": false,
          "id": "g1t9kpqo",
          "name": "collection_id",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "fqvzbaze",
          "name": "record_id",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        }
      ],
      "indexes": [
        "CREATE UNIQUE INDEX `idx_w4uoK0u` ON `authorizations` (\n  `owner`,\n  `collection_id`,\n  `record_id`\n)"
      ],
      "listRule": "@request.auth.id = owner.id || users.id ?= @request.auth.id",
      "viewRule": "@request.auth.id = owner.id || users.id ?= @request.auth.id",
      "createRule": "@request.auth.id != ''",
      "updateRule": "@request.auth.id != owner.id",
      "deleteRule": "@request.auth.id != owner.id",
      "options": {}
    },
    {
      "id": "sijd38964ht83q5",
      "created": "2023-07-11 13:26:18.946Z",
      "updated": "2023-07-11 13:28:21.425Z",
      "name": "authorizationsExamples",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "s6z3yyrm",
          "name": "name",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        }
      ],
      "indexes": [],
      "listRule": "(@collection.authorizations.users.id ?= @request.auth.id || @collection.authorizations.owner.id = @request.auth.id) && @collection.authorizations.record_id = id",
      "viewRule": "(@collection.authorizations.users.id ?= @request.auth.id || @collection.authorizations.owner.id = @request.auth.id) && @collection.authorizations.record_id = id",
      "createRule": "@request.auth.id != ''",
      "updateRule": "(@collection.authorizations.users.id ?= @request.auth.id || @collection.authorizations.owner.id = @request.auth.id) && @collection.authorizations.record_id = id",
      "deleteRule": "(@collection.authorizations.users.id ?= @request.auth.id || @collection.authorizations.owner.id = @request.auth.id) && @collection.authorizations.record_id = id",
      "options": {}
    },
    {
      "id": "_pb_users_auth_",
      "created": "2023-07-11 15:40:56.391Z",
      "updated": "2023-07-11 15:40:56.394Z",
      "name": "users",
      "type": "auth",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "users_name",
          "name": "name",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "users_avatar",
          "name": "avatar",
          "type": "file",
          "required": false,
          "unique": false,
          "options": {
            "maxSelect": 1,
            "maxSize": 5242880,
            "mimeTypes": [
              "image/jpeg",
              "image/png",
              "image/svg+xml",
              "image/gif",
              "image/webp"
            ],
            "thumbs": null,
            "protected": false
          }
        },
        {
          "system": false,
          "id": "nf3jan6h",
          "name": "bitcoin_public_key",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "zjvzw9wa",
          "name": "ecdh_public_key",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "yg9yz2sc",
          "name": "eddsa_public_key",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "ouslgk2m",
          "name": "ethereum_address",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "vqu3zutd",
          "name": "reflow_public_key",
          "type": "text",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        }
      ],
      "indexes": [
        "CREATE INDEX `__pb_users_auth__created_idx` ON `users` (`created`)"
      ],
      "listRule": "id = @request.auth.id",
      "viewRule": "id = @request.auth.id",
      "createRule": "",
      "updateRule": "id = @request.auth.id",
      "deleteRule": "id = @request.auth.id",
      "options": {
        "allowEmailAuth": true,
        "allowOAuth2Auth": true,
        "allowUsernameAuth": true,
        "exceptEmailDomains": null,
        "manageRule": null,
        "minPasswordLength": 8,
        "onlyEmailDomains": null,
        "requireEmail": false
      }
    }
  ];

  const collections = snapshot.map((item) => new Collection(item));

  return Dao(db).importCollections(collections, true, null);
}, (db) => {
  return null;
})
