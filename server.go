package husky

import (
	"fmt"
	"net/http"
)

type Server struct {
	Address string
	Port    int
	Path    string
}

func NewServer() *Server {
	return &Server{
		Address: "127.0.0.1",
		Port:    80,
		Path:    "/",
	}
}

func (s *Server) Listen(config *Config, router *Router) {
	host, _ := config.Get("host")
	port, _ := config.Get("port")
	path, _ := config.Get("path")
	handlePath, _ := path.(string)

	bind := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Server Listen: %s\n", bind)
	http.Handle(handlePath, router)
	for _, ws := range router.WebSockets {
		http.Handle(ws.Route, ws.Handler)
	}
	if err := http.ListenAndServe(bind, nil); err != nil {
		panic(fmt.Sprintf("%v\n", err))
	}
}
