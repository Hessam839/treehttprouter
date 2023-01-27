package treehttprouter

import (
	"bytes"
	"fmt"
)

type Response struct {
	Headers    map[string]string
	Body       *bytes.Buffer
	StatusCode int
}

func (r Response) String() string {
	return fmt.Sprintf("code=%d\nheaders: %v\nbody: %v", r.StatusCode, r.Headers, r.Body)
}
