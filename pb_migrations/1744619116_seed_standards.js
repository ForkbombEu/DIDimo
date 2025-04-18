// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />

/** @type {Array<object>} */
const standards = [
  {
    id: "la6i602ibti2len",
    name: "OpenID4VP",
    version: "Draft 23",
    description: "<p>OpenID for Verifiable Presentations</p>",
    standard_url: "https://openid.net/specs/openid-4-verifiable-presentations-1_0-23.html",
    latest_update: "2024-12-02 12:00:00.000Z",
    siblings: ["la6i602ibti2len", "n6n50a3y07bzyr9", "6wo0xf5s413a9ry"],
    test_suites: [],
    standard_uid: "",
    external_links: null,
  },
  {
    id: "n6n50a3y07bzyr9",
    name: "OpenID4VP",
    version: "Draft 24",
    description: "<p>OpenID for Verifiable Presentations</p>",
    standard_url: "https://openid.net/specs/openid-4-verifiable-presentations-1_0-24.html",
    latest_update: "2025-01-27 12:00:00.000Z",
    siblings: [],
    test_suites: [],
    standard_uid: "",
    external_links: null,
  },
  {
    id: "6wo0xf5s413a9ry",
    name: "OpenID4VP",
    version: "Draft 25",
    description: "<p>OpenID for Verifiable Presentations</p>",
    standard_url: "https://openid.net/specs/openid-4-verifiable-presentations-1_0-25.html",
    latest_update: "2025-04-03 12:00:00.000Z",
    siblings: [],
    test_suites: [],
    standard_uid: "",
    external_links: null,
  },
  {
    id: "sx61644z8v6o462",
    name: "OpenID4VCI",
    version: "Draft 15",
    description: "<p>OpenID for Verifiable Credential Issuance</p>",
    standard_url: "https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0-15.html",
    latest_update: "2024-12-19 12:00:00.000Z",
    siblings: [],
    test_suites: [],
    standard_uid: "",
    external_links: null,
  },
  {
    id: "gc2i346e1rcuet8",
    name: "OpenID4VCI",
    version: "Draft 14",
    description: "<p>OpenID for Verifiable Credential Issuance&nbsp;</p>",
    standard_url: "https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0-14.html",
    latest_update: "2024-08-21 12:00:00.000Z",
    siblings: [],
    test_suites: [],
    standard_uid: "",
    external_links: null,
  },
  {
    id: "978155059bw1nss",
    name: "OpenID4VCI",
    version: "Draft 13",
    description: "<p>OpenID for Verifiable Credential Issuance</p>",
    standard_url: "https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0-13.html",
    latest_update: "2024-02-08 12:00:00.000Z",
    siblings: ["978155059bw1nss", "gc2i346e1rcuet8", "sx61644z8v6o462"],
    test_suites: [],
    standard_uid: "",
    external_links: null,
  },
];

//

migrate((app) => {
  const collection = app.findCollectionByNameOrId("standards");

  standards
    .map((item) => new Record(collection, item))
    .forEach((record) => app.save(record));
});
