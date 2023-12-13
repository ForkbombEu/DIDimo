migrate((db) => {
    const dao = new Dao(db);

    const settings = dao.findSettings();
    settings.meta.appName = "didimo";

    dao.saveSettings(settings);
});
