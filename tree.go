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

	methodNotAllowed *node

	middlewares []*Handler
}

//var Tree *MuxTree

func NewMux() *MuxTree {
	var rnf Handler = func(r *http.Request) error {
		return ErrorRouteNotFound
	}

	n := newNode("*")
	_ = n.addRoute("*", &rnf)

	var mnf Handler = func(r *http.Request) error {
		return ErrorMethodNotAllowed
	}

	n1 := newNode("*")
	_ = n1.addRoute("*", &mnf)

	tree := &MuxTree{
		routeNotFound:    n,
		methodNotAllowed: n1,
		middlewares:      make([]*Handler, 0),
	}
	return tree
}

func (t *MuxTree) AddHandler(method string, path string, h Handler) error {
	current, next := split(path)

	if t.root == nil {
		node := newNode(current)
		t.root = node
		if next == "" {
			if err := node.addRoute(method, &h); err != nil {
				return err
			}
			return nil
		}
	}
	//if current == "/" && next == "" {
	//	if t.root == nil {
	//		node := newNode(path)
	//		err := node.addRoute(method, &h)
	//		if err != nil {
	//			return err
	//		}
	//		t.root = node
	//
	//		return nil
	//	}
	//	return errors.New("root route defined")
	//}
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
		handler, _ := t.routeNotFound.getHandler("*")
		return *handler
	}
	handler, err := node.getHandler(r.Method)
	if err != nil {
		h, _ := t.methodNotAllowed.getHandler("*")
		return *h
	}
	return *handler
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

func (t *MuxTree) Mount(path string, tree *MuxTree) {
	node := t.root.search(path)
	node.addChild(tree.root)
}
