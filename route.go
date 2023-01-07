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

func (r *Route) add(method string, h *Handler) *Route {
	r.handler[method] = h
	return r
}

func (r Route) get(h *Handler) {
	r.add("GET", h)
}

func (r Route) post(h *Handler) {
	r.add("POST", h)
}

func (r Route) put(h *Handler) {
	r.add("PUT", h)
}

func (r Route) delete(h *Handler) {
	r.add("DELETE", h)
}
