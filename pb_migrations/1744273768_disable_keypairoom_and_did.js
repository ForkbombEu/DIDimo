// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// @ts-check
/// <reference path="../pb_data/types.d.ts" />

//

migrate((app) => {
    const keypairoom = app.findFirstRecordByFilter(
        "features",
        `name="keypairoom"`
    );
    keypairoom.set("active", false);
    app.save(keypairoom);

    const DID = app.findFirstRecordByFilter("features", `name="DID"`);
    DID.set("active", false);
    app.save(DID);
});
