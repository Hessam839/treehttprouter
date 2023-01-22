package treehttprouter

import (
	"net/http"
)

var (
	readBuff = 1 << 12
)

type Context struct {
	//Connection net.Conn
	Request *http.Request
}

func NewCtx(req *http.Request) (*Context, error) {
	//buff := make([]byte, readBuff)
	//readLen, rer := c.Read(buff)
	//if rer != nil {
	//	return nil, fmt.Errorf("read from connection failed: %v", rer)
	//}
	//
	//req, qer := http.ReadRequest(bufio.NewReader(bytes.NewReader(buff[:readLen])))
	//if qer != nil {
	//	return nil, fmt.Errorf("read http 1.1 request failed: %v", qer)
	//}

	return &Context{
		//Connection: c,
		Request: req,
	}, nil
}
