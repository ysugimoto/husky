package husky

import (
	"encoding/json"
	"net/http"
)

type Output struct {
	resp       http.ResponseWriter
	statusCode int
	headers    map[string]string
	headerSent bool
}

func NewOutput(res http.ResponseWriter) *Output {
	return &Output{
		resp:       res,
		statusCode: 200,
		headerSent: false,
		headers:    map[string]string{},
	}
}

func (o *Output) IsHeadersSent() bool {
	return o.headerSent
}

func (o *Output) SetStatus(code int) *Output {
	o.statusCode = code

	return o
}

func (o *Output) SetHeader(key, value string) *Output {
	o.headers[key] = value

	return o
	o.resp.Header().Set(key, value)
}

func (o *Output) SendHeader() *Output {
	if o.IsHeadersSent() == false {
		header := o.resp.Header()

		for key, value := range o.headers {
			header.Set(key, value)
		}

		o.resp.WriteHeader(o.statusCode)
	}

	return o
}

func (o *Output) Send(out interface{}) {
	o.SendHeader()

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
