package treehttprouter

import (
	"context"
	"errors"
	"strings"
)

var (
	MethodError = errors.New("wrong method")
)

type Handler func(ctx context.Context) error

type Route struct {
	handler map[string]*Handler
}

func newRoute() *Route {
	return &Route{
		handler: make(map[string]*Handler),
	}
}

func (r *Route) haveMethod(method string) bool {
	return r.handler[method] != nil
}

func (r *Route) addHandler(method string, h *Handler) error {
	m := strings.ToUpper(method)

	switch m {
	case "GET":
		r.get(h)
	case "POST":
		r.post(h)
	case "PUT":
		r.put(h)
	case "DELETE":
		r.delete(h)
	case "HEAD":
		r.head(h)
	case "OPTION":
		r.option(h)
	case "*":
		r.global(h)
	default:
		return MethodError
	}

	return nil
}

func (r *Route) get(h *Handler) {
	r.handler["GET"] = h
}

func (r *Route) post(h *Handler) {
	r.handler["POST"] = h
}

func (r *Route) put(h *Handler) {
	r.handler["PUT"] = h
}

func (r *Route) delete(h *Handler) {
	r.handler["DELETE"] = h
}

func (r *Route) head(h *Handler) {
	r.handler["HEAD"] = h
}

func (r *Route) option(h *Handler) {
	r.handler["OPTION"] = h
}

func (r *Route) global(h *Handler) {
	r.handler["*"] = h
}
