package treehttprouter

import (
	"bytes"
	"net/http"
)

var (
	readBuff = 1 << 12
)

type Context struct {
	Request  *http.Request
	Response *Response
}

func NewCtx(req *http.Request) (*Context, error) {
	var b []byte
	r := &Response{
		headers:    make(map[string]string),
		body:       bytes.NewBuffer(b),
		statusCode: 0,
	}
	return &Context{
		Request:  req,
		Response: r,
	}, nil
}
