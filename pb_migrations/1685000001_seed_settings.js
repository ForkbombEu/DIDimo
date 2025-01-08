migrate((app) => {
    const settings = app.settings();
    settings.meta.appName = "{{cookiecutter.project_name}}";
    app.save(settings);
});
