package app

// SetupRouters will register routes in router
func (app *App) setRouters() {
	app.Post("/person", app.handleRequest(CreatePerson))
	app.Patch("/person/{id}", app.handleRequest(UpdatePerson))
	app.Put("/person/{id}", app.handleRequest(UpdatePerson))
	app.Get("/person/{id}", app.handleRequest(GetPerson))
	app.Get("/person", app.handleRequest(GetPersons))
	app.Get("/person", app.handleRequest(GetPersons), "page", "{page}")
}
