package husky

import (
	"strings"
)

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

type HuskyHandler func(dispatcher *Dispatcher)

func (app *app) joinPath(path string) string {
	p, _ := app.Config.Get("path")
	joined := p.(string) + path

	return "/" + strings.TrimLeft(joined, "/")
}

func (app *app) Post(route string, handler HuskyHandler) {
	app.Router.Bind("POST", app.joinPath(route), handler)
}

func (app *app) Get(route string, handler HuskyHandler) {
	app.Router.Bind("GET", app.joinPath(route), handler)
}

func (app *app) Put(route string, handler HuskyHandler) {
	app.Router.Bind("PUT", app.joinPath(route), handler)
}

func (app *app) Delete(route string, handler HuskyHandler) {
	app.Router.Bind("DELETE", app.joinPath(route), handler)
}

func (app *app) Notfound(route string, handler HuskyHandler) {
	app.Router.Bind("404", app.joinPath(route), handler)
}

func (app *app) AcceptCORS() {
	app.Router.Bind("OPTIONS", "/", func(d *Dispatcher) {
		d.Output.SetHeader("Access-Control-Allow-Origin", "*")
		d.Output.SetHeader("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		d.Output.SetHeader("Access-Control-Allow-Headers", "X-Requested-With")
		d.Output.SetStatus(204)
	})
}

func (app *app) Serve() {
	app.Server.Listen(app.Config, app.Router)
}
