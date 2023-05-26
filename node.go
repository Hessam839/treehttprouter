package treehttprouter

import "errors"

var (
	ErrorMethodDuplicate  = errors.New("method defined")
	ErrorMethodNotAllowed = errors.New("method not allowed")
)

type node struct {
	path string

	body *Route

	children []*node

	available bool

	Middlewares []*Handler
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

func (n *node) addMiddleware(h *Handler) error {
	if n.Middlewares != nil {
		n.Middlewares = append(n.Middlewares, h)
		return nil
	}

	m := make([]*Handler, 0)
	m = append(m, h)
	n.Middlewares = m
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

func (n *node) searchWithMiddleware(path string) (*node, *[]*Handler) {
	Mlist := []*Handler{}

	node := n.populateMiddlewareList(&Mlist, path)

	return node, &Mlist
}

func (n *node) populateMiddlewareList(mList *[]*Handler, path string) *node {
	current, next := split(path)
	if current == n.path && next == "" {
		//fmt.Print("current == n.path && next == '',",current,next)
		if n.Middlewares != nil {
			*mList = append(*mList, n.Middlewares...)
		}

		return n
	}
	child, remain := split(next)
	for idx := 0; idx < len(n.children); idx++ {
		if child == n.children[idx].path {

			//fmt.Print("child == n.children[idx].path,",child,n.children[idx].path)
			if n.Middlewares != nil {
				*mList = append(*mList, n.Middlewares...)
			}
			if remain == "" {

				//fmt.Print("remain == '',",remain)
				return n.children[idx]
			}

			return n.children[idx].populateMiddlewareList(mList, next)
		}
	}
	return nil
}

func (n *node) getHandler(method string) (*Handler, error) {
	h, ok := n.body.handler[method]
	if !ok {
		return nil, ErrorMethodNotAllowed
	}
	return h, nil
}

func (n *node) availableNode() {
	n.available = true
}

func (n *node) unavailableNode() {
	n.available = false
}
