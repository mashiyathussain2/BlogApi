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

	app.Post("/comment", app.handleRequest(CreateComment))
	//app.Patch("/comment/{id}", app.handleRequest(UpdateComment))
	//app.Put("/comment/{id}", app.handleRequest(UpdateComment))
	app.Get("/comment/{id}", app.handleRequest(GetComment))
	app.Get("/comment", app.handleRequest(GetComments))

	app.Post("/login", app.handleRequest(Login))

}
