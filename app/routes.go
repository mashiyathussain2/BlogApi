package app

// SetupRouters will register routes in router
func (app *App) setRouters() {
	app.Post("/person", app.handleRequest(CreatePerson))
	app.Patch("/person/{id}", app.handleRequest(UpdatePerson))
	app.Put("/person/{id}", app.handleRequest(UpdatePerson))
	app.Get("/person/{id}", app.handleRequest(GetPerson))
	app.Get("/person", app.handleRequest(GetPersons))

	app.Post("/blogpage", app.handleRequest(CreateBlog))
	app.Patch("/blogpage/{id}", app.handleRequest(UpdateBlog))
	app.Put("/blogpage/{id}", app.handleRequest(UpdateBlog))
	app.Get("/blogpage/{id}", app.handleRequest(GetBlog))
	app.Get("/blogpage", app.handleRequest(GetBlogs))
}
