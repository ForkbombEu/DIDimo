/// <reference path="../pb_data/types.d.ts" />

const ADMIN_EMAIL = "admin@example.org";

migrate(
    (db) => {
        const dao = new Dao(db);

        try {
            dao.findAdminByEmail(ADMIN_EMAIL);
        } catch {
            const admin = new Admin();
            admin.email = ADMIN_EMAIL;
            admin.setPassword("adminadmin");
            return dao.saveAdmin(admin);
        }
    },
    (db) => {
        const dao = new Dao(db);
        const admin = dao.findAdminByEmail(ADMIN_EMAIL);
        return dao.deleteAdmin(admin);
    }
);
