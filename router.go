package husky

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
)

type Router struct {
	Routes []*RouterInfo
}

type RouterInfo struct {
	Method  []byte
	Route   string
	Handler HandlerFunc
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Bind(method, route string, handler HandlerFunc) {
	info := &RouterInfo{
		Method:  bytes.ToUpper([]byte(method)),
		Route:   route,
		Handler: handler,
	}

	r.Routes = append(r.Routes, info)
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
	response := Response{res}
	request := &Request{req}
	dispatcher.Handler(response, request)
}

func (r *Router) Send404(res http.ResponseWriter, req *http.Request) {
	dispatcher, notFound := r.FindDispatcher("404", "/")
	if notFound {
		r.SendDefault404(res, req)
		return
	}

	response := Response{res}
	request := &Request{req}
	dispatcher.Handler(response, request)
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
	res.WriteHeader(404)
	res.Header().Set("Content-Type", "text/plain")
	res.Write([]byte("Not found."))
}

func (r *Router) Dispatch(handlePath string) {
	if handlePath == "" {
		handlePath = "/"
	}
	http.HandleFunc(handlePath, r.ServeHTTP)
}
