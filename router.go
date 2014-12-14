package husky

import (
	"bytes"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"regexp"
)

type Router struct {
	Routes     []*RouterInfo
	WebSockets []*WebSocketInfo
}

type RouterInfo struct {
	Method  []byte
	Route   string
	Handler HuskyHandler
}

type WebSocketInfo struct {
	Route   string
	Handler websocket.Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Bind(method, route string, handler HuskyHandler) {
	info := &RouterInfo{
		Method:  bytes.ToUpper([]byte(method)),
		Route:   route,
		Handler: handler,
	}

	r.Routes = append(r.Routes, info)
}

func (r *Router) BindWebSocket(route string, handler HuskyWebSocketHandler) {
	info := &WebSocketInfo{
		Route:   route,
		Handler: NewWebSocketHandler(handler),
	}

	r.WebSockets = append(r.WebSockets, info)
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var dispatcher *RouterInfo
	var notFound bool

	requestURL := req.URL

	if dispatcher, notFound = r.FindDispatcher(req.Method, requestURL.Path); notFound {
		r.Send404(res, req)
		return
	}

	log.Printf("%s %s\n", req.Method, requestURL.Path)
	dispatch := NewDispatcher(res, req)
	dispatcher.Handler(dispatch)
}

func (r *Router) Send404(res http.ResponseWriter, req *http.Request) {
	dispatcher, notFound := r.FindDispatcher("404", "/")
	if notFound {
		r.SendDefault404(res, req)
		return
	}

	dispatch := NewDispatcher(res, req)
	dispatcher.Handler(dispatch)
}

func (r *Router) FindDispatcher(method, path string) (dispatcher *RouterInfo, notFound bool) {
	for i := 0; i < len(r.Routes); i++ {
		info := r.Routes[i]
		if string(info.Method) != method {
			continue
		}

		regex := regexp.MustCompile(info.Route)
		matched := regex.FindString(path)
		if matched == "" {
			continue
		}

		return info, false
	}

	return nil, true
}

func (r *Router) SendDefault404(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(404)
	res.Write([]byte("Not found."))
}
