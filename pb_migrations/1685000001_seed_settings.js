migrate((app) => {
    const settings = app.settings();
    settings.meta.appName = "DIDimo";
    app.save(settings);
});
