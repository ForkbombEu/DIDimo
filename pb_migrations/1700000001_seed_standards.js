// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// @ts-check

/// <reference path="../pb_data/types.d.ts" />

/** @type {Array<{name:string, level:number}>} */

const standards = [
    {
        created: "2025-04-11 11:34:58.303Z",
        description: "<p>OpenID for Verifiable Presentations</p>",
        external_links: null,
        id: "la6i602ibti2len",
        latest_update: "2024-12-02 12:00:00.000Z",
        name: "OpenID4VP",
        siblings: [],
        standard_uid: "",
        standard_url:
            "https://openid.net/specs/openid-4-verifiable-presentations-1_0-23.html",
        test_suites: [],
        updated: "2025-04-11 11:43:11.743Z",
        version: "Draft 23",
    },
    {
        created: "2025-04-11 11:36:18.410Z",
        description: "<p>OpenID for Verifiable Presentations</p>",
        external_links: null,
        id: "n6n50a3y07bzyr9",
        latest_update: "2025-01-27 12:00:00.000Z",
        name: "OpenID4VP",
        siblings: [],
        standard_uid: "",
        standard_url:
            "https://openid.net/specs/openid-4-verifiable-presentations-1_0-24.html",
        test_suites: [],
        updated: "2025-04-11 11:38:36.097Z",
        version: "Draft 24",
    },
    {
        created: "2025-04-11 11:37:09.054Z",
        description: "<p>OpenID for Verifiable Presentations</p>",
        external_links: null,
        id: "6wo0xf5s413a9ry",
        latest_update: "2025-04-03 12:00:00.000Z",
        name: "OpenID4VP",
        siblings: [],
        standard_uid: "",
        standard_url:
            "https://openid.net/specs/openid-4-verifiable-presentations-1_0-25.html",
        test_suites: [],
        updated: "2025-04-11 11:38:44.909Z",
        version: "Draft 25",
    },
    {
        created: "2025-04-11 11:40:23.318Z",
        description: "<p>OpenID for Verifiable Credential Issuance</p>",
        external_links: null,
        id: "sx61644z8v6o462",
        latest_update: "2024-12-19 12:00:00.000Z",
        name: "OpenID4VCI",
        siblings: [],
        standard_uid: "",
        standard_url:
            "https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0-15.html",
        test_suites: [],
        updated: "2025-04-11 11:40:23.318Z",
        version: "Draft 15",
    },
    {
        created: "2025-04-11 11:41:06.006Z",
        description: "<p>OpenID for Verifiable Credential Issuance&nbsp;</p>",
        external_links: null,
        id: "gc2i346e1rcuet8",
        latest_update: "2024-08-21 12:00:00.000Z",
        name: "OpenID4VCI",
        siblings: [],
        standard_uid: "",
        standard_url:
            "https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0-14.html",
        test_suites: [],
        updated: "2025-04-11 11:44:26.778Z",
        version: "Draft 14",
    },
    {
        created: "2025-04-11 11:41:57.136Z",
        description: "<p>OpenID for Verifiable Credential Issuance</p>",
        external_links: null,
        id: "978155059bw1nss",
        latest_update: "2024-02-08 12:00:00.000Z",
        name: "OpenID4VCI",
        siblings: [],
        standard_uid: "",
        standard_url:
            "https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0-13.html",
        test_suites: [],
        updated: "2025-04-11 11:43:24.886Z",
        version: "Draft 13",
    },
];

migrate((app) => {
    const collection = app.findCollectionByNameOrId("standards");

    standards.map((s) => new Record(collection, s)).forEach((r) => app.save(r));
});
