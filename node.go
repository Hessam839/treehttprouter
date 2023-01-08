package treehttprouter

import "errors"

var (
	ErrorMethodDuplicate = errors.New("method defined")
)

type node struct {
	path string

	body *Route

	children []*node

	available bool
}

func newNode(path string) *node {
	return &node{
		path:      path,
		body:      nil,
		children:  nil,
		available: true,
	}
}

func (n *node) haveBody() bool {
	return n.body != nil
}

func (n *node) isAvailable() bool {
	return n.available
}

func (n *node) addRoute(method string, h *Handler) error {
	if n.haveBody() {
		if n.body.haveMethod(method) {
			return ErrorMethodDuplicate
		}
		err := n.body.addHandler(method, h)
		if err != nil {
			return err
		}
		return nil
	}

	r := newRoute()
	err := r.addHandler(method, h)
	if err != nil {
		return err
	}
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

func (n *node) getHandler(method string) *Handler {
	return n.body.handler[method]
}

func (n *node) availableNode() {
	n.available = true
}

func (n *node) unavailableNode() {
	n.available = false
}
