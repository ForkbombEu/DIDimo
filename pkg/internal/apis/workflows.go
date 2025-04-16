// All the related calls to the workflows

app.OnServe().BindFunc(func(se *core.ServeEvent) error {
	g := se.Router.Group("/api/compliance/check")

	// RICOPIA la roba dei workflows
	// {workflow}/{run}/hhistory
}


