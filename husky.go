package husky

import (
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
}

type Response struct {
	http.ResponseWriter
}

type app struct {
	Router *Router
	Config *Config
	Server *Server
}

func NewApp() *app {
	return &app{
		Router: NewRouter(),
		Config: NewConfig(),
		Server: NewServer(),
	}
}

type HandlerFunc func(Response, *Request)

func (app *app) joinPath(path string) string {
	p, _ := app.Config.Get("path")
	joined := p.(string) + path

	return "/" + strings.TrimLeft(joined, "/")
}

func (app *app) Post(route string, handler HandlerFunc) {
	app.Router.Bind("POST", app.joinPath(route), handler)
}

func (app *app) Get(route string, handler HandlerFunc) {
	app.Router.Bind("GET", app.joinPath(route), handler)
}

func (app *app) Put(route string, handler HandlerFunc) {
	app.Router.Bind("PUT", app.joinPath(route), handler)
}

func (app *app) Delete(route string, handler HandlerFunc) {
	app.Router.Bind("DELETE", app.joinPath(route), handler)
}

func (app *app) Notfound(route string, handler HandlerFunc) {
	app.Router.Bind("404", app.joinPath(route), handler)
}

func (app *app) Serve() {
	app.Server.Listen(app.Config, app.Router)
}
