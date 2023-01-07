package treehttprouter

import (
	"github.com/stretchr/testify/assert"
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
	tree := newTree()

	var v1 Handler = func(r *http.Request) {}
	if err := tree.AddHandler("GET", "/", v1); err != nil {
		t.Fatalf("cant create handler: %v", err)
	}

	if err := tree.AddHandler("GET", "/api/v1/users", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	if err := tree.AddHandler("POST", "/api/v1/users", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	if err := tree.AddHandler("PUT", "/api/v1/users", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	if err := tree.AddHandler("GET", "/api/v1/products", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	if err := tree.AddHandler("POST", "/api/v1/products", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}
}

func TestNodeSearch(t *testing.T) {
	var v1 Handler = func(r *http.Request) {}

	n3 := newNode("admin")
	err := n3.addRoute("GET", &v1)
	err = n3.addRoute("POST", &v1)
	n4 := newNode("users")
	err = n4.addRoute("GET", &v1)
	n5 := newNode("v1")
	err = n5.addRoute("GET", &v1)
	n2 := newNode("api")
	err = n2.addRoute("GET", &v1)
	n1 := newNode("/")
	err = n1.addRoute("POST", &v1)
	if err != nil {
		t.Fatalf("")
	}
	// /api
	n1.addChild(n2.addChild(n5)).addChild(n3.addChild(n4))

	nde := n1.search("/api/v1")
	assert.NotNil(t, nde)
	nde = n1.search("/admin/users")
	assert.NotNil(t, nde)
	nde = n1.search("/api/v1/user")
	assert.Nil(t, nde)
}

func TestMatch(t *testing.T) {
	tree := newTree()

	var v1 Handler = func(r *http.Request) {}
	if err := tree.AddHandler("GET", "/", v1); err != nil {
		t.Fatalf("cant create handler: %v", err)
	}

	if err := tree.AddHandler("GET", "/api/v1/users", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	if err := tree.AddHandler("POST", "/api/v1/users", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	if err := tree.AddHandler("PUT", "/api/v1/users", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	if err := tree.AddHandler("GET", "/api/v1/products", v1); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	if err := tree.AddHandler("POST", "/api/v1/products", func(r *http.Request) {
		t.Logf("incoming req path: %v", r.URL.Path)
	}); err != nil {
		t.Fatalf("cant create handler %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/products", nil)
	if err != nil {
		t.Fatalf("with error: %v", err)
	}

	handler := tree.Match(req)
	handler(req)
}
