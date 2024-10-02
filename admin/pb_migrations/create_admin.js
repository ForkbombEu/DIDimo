// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

migrate(
    (db) => {
        const admin = new Admin();
        admin.email = "admin@example.org";
        admin.setPassword("adminadmin");
        return Dao(db).saveAdmin(admin);
    },
    (db) => {
        const dao = new Dao(db);
        const admin = dao.findAdminByEmail("admin@example.org");
        return dao.deleteAdmin(admin);
    }
);
