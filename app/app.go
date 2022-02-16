package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"blog/app/db"
	"blog/app/handler"

	//"blog/app/handler"
	"blog/config"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

// App has the mongo database and router instances
type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

// ConfigAndRunApp will create and initialize App structure. App factory function.
func ConfigAndRunApp(config *config.Config) {
	app := new(App)
	app.Initialize(config)
	app.Run(config.ServerHost)
}

// Initialize initialize the app with
func (app *App) Initialize(config *config.Config) {
	app.DB = db.InitialConnection("golang", config.MongoURI())
	app.Router = mux.NewRouter()
	app.UseMiddleware(handler.JSONContentTypeMiddleware)
	app.setRouters()
}

// UseMiddleware will add global middleware in router
func (app *App) UseMiddleware(middleware mux.MiddlewareFunc) {
	app.Router.Use(middleware)
}

// Run will start the http server on host that you pass in. host:<ip:port>
func (app *App) Run(host string) {
	// use signals for shutdown server gracefully.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)
	go func() {
		log.Fatal(http.ListenAndServe(host, app.Router))
	}()
	log.Printf("Server is listning on http://%s\n", host)
	sig := <-sigs
	log.Println("Signal: ", sig)

	log.Println("Stoping MongoDB Connection...")
	app.DB.Client().Disconnect(context.Background())
}

// RequestHandlerFunction is a custome type that help us to pass db arg to all endpoints
type RequestHandlerFunction func(db *mongo.Database, w http.ResponseWriter, r *http.Request)

// handleRequest is a middleware we create for pass in db connection to endpoints.
func (app *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(app.DB, w, r)
	}
}

// Get will register Get method for an endpoint
func (app *App) Get(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("GET").Queries(queries...)
}

// Post will register Post method for an endpoint
func (app *App) Post(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("POST").Queries(queries...)
}

// Put will register Put method for an endpoint
func (app *App) Put(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("PUT").Queries(queries...)
}

// Patch will register Patch method for an endpoint
func (app *App) Patch(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("PATCH").Queries(queries...)
}

// Delete will register Delete method for an endpoint
func (app *App) Delete(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("DELETE").Queries(queries...)
}
