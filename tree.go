package treehttprouter

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrorRouteNotFound = errors.New("route not found")
)

type tree struct {
	root *node

	routeNotFound *node

	middlewares []*Handler
}

//var Tree *tree

func newTree() *tree {
	var rnf Handler = func(r *http.Request) error {
		return ErrorRouteNotFound
	}

	n := newNode("*")
	_ = n.addRoute("*", &rnf)

	tree := &tree{
		routeNotFound: n,
	}
	return tree
}

func (t *tree) AddHandler(method string, path string, h Handler) error {
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

func (t *tree) Match(r *http.Request) Handler {
	path := r.URL.Path
	node := t.root.search(path)
	if node == nil || !node.isAvailable() {
		return *t.routeNotFound.getHandler("*")
	}
	return *node.getHandler(r.Method)
}

func (t *tree) DisablePath(path string) {
	node := t.root.search(path)
	if node != nil {
		node.unavailableNode()
	}
}

func (t *tree) EnablePath(path string) {
	node := t.root.search(path)
	if node != nil {
		node.availableNode()
	}
}
func split(path string) (string, string) {
	p := strings.Split(path, "/")
	if len(p) == 0 {
		return "/", ""
	}
	if p[0] == "" {
		return "/", strings.Join(p[1:], "/")
	}
	return p[0], strings.Join(p[1:], "/")
}
