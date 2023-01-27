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
		Headers:    make(map[string]string),
		Body:       bytes.NewBuffer(b),
		StatusCode: 0,
	}
	return &Context{
		Request:  req,
		Response: r,
	}, nil
}
