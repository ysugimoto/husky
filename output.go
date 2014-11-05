package husky

import (
	"encoding/json"
	"net/http"
)

type Output struct {
	resp http.ResponseWriter
}

func NewOutput(res http.ResponseWriter) *Output {
	return &Output{
		resp: res,
	}
}

func (o *Output) SetStatus(code int) {
	o.resp.WriteHeader(code)
}

func (o *Output) SetHeader(key, value string) {
	o.resp.Header().Set(key, value)
}

func (o *Output) Send(out interface{}) {
	if buf, ok := toBuffer(out); ok {
		o.resp.Write(buf)
		return
	}

	panic("Invalid output data")
}

func (o *Output) SendJSON(out interface{}) {
	if buf, ok := toJSONString(out); ok {
		o.SetHeader("Content-Type", "application/json")
		o.resp.Write(buf)
		return
	}

	panic("Invalid output data for JSON")
}

func toBuffer(data interface{}) ([]byte, bool) {
	if str, ok := data.(string); ok {
		return []byte(str), true
	} else if b, ok := data.([]byte); ok {
		return b, true
	} else {
		return nil, false
	}
}

func toJSONString(data interface{}) ([]byte, bool) {
	if str, ok := data.(string); ok {
		return []byte(str), true
	} else if json, err := json.Marshal(data); err != nil {
		return json, true
	} else {
		return nil, false
	}
}
