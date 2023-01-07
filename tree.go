package treehttprouter

import (
	"errors"
	"net/http"
	"strings"
)

type node struct {
	path     string
	body     *Route
	children []*node
}

func newNode(path string) *node {
	return &node{
		path:     path,
		body:     nil,
		children: nil,
	}
}

func (n *node) addRoute(method string, h *Handler) error {
	if n.body != nil {
		handler := n.body.handler[method]
		if handler != nil {
			return errors.New("method defined")
		}
		n.body.add(method, h)
		return nil
	}

	r := newRoute()
	r.add(method, h)
	n.body = r
	return nil
}

func (n *node) addChild(child *node) *node {
	n.children = append(n.children, child)
	return n
}

func (n *node) haveChild(path string) *node {
	for _, child := range n.children {
		if child.path == path {
			return child
		}
	}
	return nil
}

func (n *node) search(path string) *node {
	current, next := split(path)
	if current == n.path && next == "" {
		return n
	}
	child, remain := split(next)
	for idx := 0; idx < len(n.children); idx++ {
		if child == n.children[idx].path {
			if remain == "" {
				return n.children[idx]
			}
			return n.children[idx].search(next)
		}
	}
	return nil
}

func (n node) getHandler(method string) *Handler {
	return n.body.handler[method]
}

type tree struct {
	root *node

	routeNotFound *node
}

//var Tree *tree

func newTree() *tree {
	r := newRoute()
	var rnf Handler = func(r *http.Request) {}
	r.add("*", &rnf)
	tree := &tree{
		routeNotFound: &node{
			path:     "*",
			body:     r,
			children: nil,
		},
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
	if node == nil {
		return *t.routeNotFound.getHandler("*")
	}
	return *node.getHandler(r.Method)
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
