// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// @ts-check
/// <reference path="../pb_data/types.d.ts" />

/** @type {{name:string, envVariables?:Record<string,unknown>, active?:boolean}[]} */
const features = [
    {
        name: "keypairoom",
        envVariables: {
            SALT: "bWltbW8K",
        },
    },
    {
        name: "DID",
        envVariables: {
            DID_URL: "url",
            DID_SPEC: "string",
            DID_SIGNER_SPEC: "string",
            DID_IDENTITY: "string",
            DID_KEYRING: "json, currently passed base64 encoded",
        },
    },
    {
        name: "auth",
    },
    {
        name: "maintenance",
        active: false,
    },
    {
        name: "organizations",
    },
    {
        name: "webauthn",
        envVariables: {
            DISPLAY_NAME: "{{cookiecutter.project_name}}",
            RPID: "localhost",
            RPORIGINS: "http://localhost:5173",
        },
    },
    {
        name: "oauth",
    },
    {
        name: "demo",
        active: false,
    },
];

//

migrate((app) => {
    const featuresCollection = app.findCollectionByNameOrId("features");

    features
        .map(
            (feature) =>
                new Record(featuresCollection, {
                    ...feature,
                    active: feature.active ?? true,
                })
        )
        .forEach((featureRecord) => app.save(featureRecord));
});
