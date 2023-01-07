package treehttprouter

import (
	"net/http"
	"testing"
)

func TestSplit(t *testing.T) {
	p1 := "/api/v2/users"
	p2 := ""
	for {
		p1, p2 = split(p1)
		t.Logf("p1= %v, p2= %v", p1, p2)
		if p2 == "" {
			break
		}
		p1 = p2
		p2 = ""
	}
}

func TestTree(t *testing.T) {
	newTree()

	var v1 Handler = func(r *http.Request) {}
	Tree.add(&node{
		path: "/",
		body: &Route{map[string]*Handler{
			"GET": &v1,
		}},
		children: nil,
	})
	Tree.add(&node{
		path: "/api",
		body: &Route{map[string]*Handler{
			"POST": &v1,
		}},
		children: nil,
	})
}

func TestTraverse(t *testing.T) {
	var v1 Handler = func(r *http.Request) {}

	n2 := &node{
		path: "admin",
		body: &Route{map[string]*Handler{
			"GET": &v1,
		}},
		children: nil,
	}

	n3 := &node{
		path: "api",
		body: &Route{map[string]*Handler{
			"GET": &v1,
		}},
		children: nil,
	}

	n1 := &node{
		path: "/",
		body: &Route{map[string]*Handler{
			"POST": &v1,
		}},
	}
	n1.children = append(n1.children, n2)
	n1.children = append(n1.children, n3)

	node := traverse(n1, "/api")
	_ = node
}
