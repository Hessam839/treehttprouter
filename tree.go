package treehttprouter

import (
	"net/http"
	"strings"
)

type node struct {
	path     string
	body     *Route
	children []*node
}

type tree struct {
	root *node

	routeNotFound *node
}

var Tree *tree

func newTree() {
	r := newRoute()
	var rnf Handler = func(r *http.Request) {}
	r.add("*", &rnf)
	Tree = &tree{
		routeNotFound: &node{
			path:     "*",
			body:     r,
			children: nil,
		},
	}
}

func (t *tree) add(n *node) error {
	//var lNode *node = nil
	//
	//if t.root == nil {
	//	t.root = n
	//	return nil
	//}
	//if n == nil {
	//	return errors.New("node must ne not empty")
	//}
	//p1, p2 := split(n.path)

	return nil
}

// path /api/v2
// - root /
// |	  /api
// |		  /v2
//

func traverse(n *node, path string) *node {
	lNode := n
	if lNode == nil {
		return nil
	}
	p1, p2 := split(path)
	if p1 != lNode.path {
		return nil
	}
	if p2 == "" {
		return lNode
	}
	p1, p2 = split(p2)
	var nde *node = nil
	for idx := 0; idx < len(lNode.children); idx++ {
		nde = traverse(lNode.children[idx], p1)
	}
	if nde != nil {
		return nde
	}
	return Tree.routeNotFound
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
