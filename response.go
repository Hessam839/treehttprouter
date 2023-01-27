package treehttprouter

import "bytes"

type Response struct {
	headers    map[string]string
	body       *bytes.Buffer
	statusCode int
}
