package husky

import "net/http"

type Dispatcher struct {
	Input  *Input
	Output *Output
}

func NewDispatcher(res http.ResponseWriter, req *http.Request) *Dispatcher {
	return &Dispatcher{
		Input:  NewInput(req),
		Output: NewOutput(res),
	}
}
