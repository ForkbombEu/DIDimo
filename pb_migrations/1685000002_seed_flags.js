// @ts-check
/// <reference path="../pb_data/types.d.ts" />

/** @type {{name:string, envVariables?:Record<string,unknown>, active?:boolean}[]} */
const flags = [
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
            DISPLAY_NAME: "DIDimo",
            RPID: "localhost",
            RPORIGINS: "http://localhost:5173",
        },
    },
    {
        name: "oauth",
    },
];

//

migrate((app) => {
    const featuresCollection = app.findCollectionByNameOrId("flags");

    flags
        .map(
            (flag) =>
                new Record(featuresCollection, {
                    ...flag,
                    active: flag.active ?? true,
                }),
        )
        .forEach((flagRecord) => app.save(flagRecord));
});
