package treehttprouter

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewRouter(t *testing.T) {
	r := newRoute()
	assert.NotNilf(t, r, "return nil router")
}

func TestAddHandler(t *testing.T) {
	r := newRoute()

	var h1 Handler = func(r *http.Request) {}
	var h2 Handler = func(r *http.Request) {}

	r.add("GET", &h1)

	r.add("POST", &h2)

	tr := &Route{
		handler: map[string]*Handler{
			"GET":  &h1,
			"POST": &h2,
		},
	}

	assert.Equal(t, tr, r)
}
