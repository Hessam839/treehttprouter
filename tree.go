package treehttprouter

import (
	"errors"
	"net/http"
)

var (
	ErrorRouteNotFound = errors.New("route not found")
)

type MuxTree struct {
	root *node

	routeNotFound *node

	middlewares []*Handler
}

//var Tree *MuxTree

func NewMux() *MuxTree {
	var rnf Handler = func(r *http.Request) error {
		return ErrorRouteNotFound
	}

	n := newNode("*")
	_ = n.addRoute("*", &rnf)

	tree := &MuxTree{
		routeNotFound: n,
		middlewares:   make([]*Handler, 0),
	}
	return tree
}

func (t *MuxTree) AddHandler(method string, path string, h Handler) error {
	current, next := split(path)

	if current == "/" && next == "" {
		if t.root == nil {
			node := newNode(path)
			_ = node.addRoute(method, &h)
			t.root = node

			return nil
		}
		return errors.New("root route defined")
	}

	currentNode := t.root
	child, remain := split(next)
	for {
		childNode := currentNode.haveChild(child)

		//if leaf node
		if remain == "" {
			if childNode == nil {
				newNode := newNode(child)
				if err := newNode.addRoute(method, &h); err != nil {
					return err
				}
				currentNode.addChild(newNode)
				return nil
			}

			err := childNode.addRoute(method, &h)
			if err != nil {
				return err
			}
			return nil
		}

		if childNode == nil {
			newNode := newNode(child)
			currentNode.addChild(newNode)
			currentNode = newNode
			child, remain = split(remain)
			continue
		}

		currentNode = childNode
		child, remain = split(remain)
	}
}

func (t *MuxTree) match(r *http.Request) Handler {
	path := r.URL.Path
	node := t.root.search(path)
	if node == nil || !node.isAvailable() {
		return *t.routeNotFound.getHandler("*")
	}
	return *node.getHandler(r.Method)
}

func (t *MuxTree) DisablePath(path string) {
	node := t.root.search(path)
	if node != nil {
		node.unavailableNode()
	}
}

func (t *MuxTree) EnablePath(path string) {
	node := t.root.search(path)
	if node != nil {
		node.availableNode()
	}
}

func (t *MuxTree) Use(handler Handler) {
	t.middlewares = append(t.middlewares, &handler)
}

func (t *MuxTree) Serve(r *http.Request) error {
	handler := t.match(r)
	for _, middleware := range t.middlewares {
		h := *middleware
		if err := h(r); err != nil {
			return err
		}
	}
	return handler(r)
}
