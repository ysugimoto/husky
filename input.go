package husky

import (
	"net/http"
)

type Input struct {
	gets    map[string]string
	posts   map[string]string
	cookies map[string]string

	req *http.Request
}

func NewInput(req *http.Request) *Input {
	i := &Input{
		req: req,
	}

	i.initialize()

	return i
}

func (i *Input) initialize() {

}

func (i *input) GetRequest() *http.Request {
	return i.req
}

func (i *Input) Get(key string) string {
	var val string

	if value, ok := i.gets[key]; ok {
		val = value
	}

	return val
}
