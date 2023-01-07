package treehttprouter

import "net/http"

type Handler func(r *http.Request)

type Route struct {
	handler map[string]*Handler
}

func newRoute() *Route {
	return &Route{
		handler: make(map[string]*Handler),
	}
}

func (r *Route) add(method string, h *Handler) {
	r.handler[method] = h
}
