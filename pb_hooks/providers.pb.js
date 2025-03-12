// @ts-check

/// <reference path="../pb_data/types.d.ts" />
/** @typedef {import('./utils.js')} Utils */

//

onRecordCreateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    if (utils.isAdminContext(e)) e.next();

    //

    const ownerId = utils.getUserFromContext(e)?.id;

    // Check if the user has already submitted a claim for this provider

    const providerClaims = $app.findRecordsByFilter(
        "provider_claims",
        [
            `provider = "${e.record?.get("provider")}"`,
            `owner = "${ownerId}"`,
        ].join(" && "),
        "created",
        1,
        0
    );

    if (providerClaims.length > 0) {
        throw new BadRequestError(
            "You have already submitted a claim for this provider."
        );
    }

    // Adding default values

    e.record?.set("status", "in_review");
    e.record?.set("owner", utils.getUserFromContext(e)?.id);

    //

    e.next();
}, "provider_claims");

//

onRecordUpdateRequest((e) => {
    /** @type {Utils} */
    const utils = require(`${__hooks}/utils.js`);

    const providerClaim = e.record;
    const originalProviderClaim = providerClaim?.original();

    // If the request is not from an admin, restore the original values
    if (!utils.isAdminContext(e)) {
        $app.runInTransaction(() => {
            ["status", "owner", "provider"].forEach((field) => {
                providerClaim?.set(field, originalProviderClaim?.get(field));
            });
        });
    }

    e.next();

    // If the request is from an admin, and the status is "approved":
    // - Update the provider
    // - Create an organization if missing
    // - Delete the claim
    // - Notify the user

    if (utils.isAdminContext(e) && providerClaim) {
        const logo = utils.getFirstFile(providerClaim, "logo");
        console.log("logoName", logo?.name);
    }

    if (
        utils.isAdminContext(e) &&
        providerClaim &&
        providerClaim.get("status") === "approved"
    ) {
        $app.runInTransaction((app) => {
            const ownerId = providerClaim.getString("owner");

            /**
             * Create an organization if missing
             */

            /** @type {string} */
            let organizationId;

            const organizations = app
                .findRecordsByFilter(
                    "organizations",
                    `claimed_by = "${ownerId}"`,
                    "",
                    0,
                    0
                )
                .filter((org) => org !== undefined);

            if (organizations.length === 1) {
                organizationId = organizations[0].id;
            } else if (organizations.length > 1) {
                throw new BadRequestError(
                    "Multiple organizations found for the same user."
                );
            } else {
                const orgCollection =
                    app.findCollectionByNameOrId("organizations");

                const newOrganization = new Record(orgCollection, {
                    name: providerClaim.get("name"),
                    claimed_by: ownerId,
                    created_by: ownerId,
                    country: providerClaim.get("country"),
                    // logo: logo,
                });

                app.save(newOrganization);
                organizationId = newOrganization.id;
            }

            /**
             * Update the provider
             */

            const provider = app.findRecordById(
                "services",
                providerClaim.get("provider")
            );

            [
                "name",
                "description",
                "country",
                "legal_entity",
                "external_links",
            ].forEach((field) => {
                provider.set(field, providerClaim.get(field));
            });
            // provider.set("logo", logo);
            provider.set("owner", organizationId);

            app.save(provider);

            //

            app.delete(providerClaim);

            console.log("Provider updated and claim deleted");
        });
    }
}, "provider_claims");
