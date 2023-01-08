package treehttprouter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
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
	tree := NewMux()

	var v1 Handler = func(r *http.Request) error { return nil }
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
	var v1 Handler = func(r *http.Request) error { return nil }

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
	nde = n1.search("/api/v1")
}

func TestMatch(t *testing.T) {
	tree, _ := CreateTree()

	req, err := http.NewRequest("POST", "/api/v1/products", nil)
	if err != nil {
		t.Fatalf("with error: %v", err)
	}

	handler := tree.match(req)
	assert.NotNil(t, handler)

	err = handler(req)
	if err != nil {
		t.Fatalf("handler match error: %v", err)
	}
}

func TestDisableRoute(t *testing.T) {
	tree, _ := CreateTree()

	tree.DisablePath("/api/v1/users")

	req, rer := http.NewRequest("POST", "/api/v1/users", nil)
	if rer != nil {
		t.Fatalf("with error: %v", rer)
	}

	handler := tree.match(req)
	assert.NotNil(t, handler)
	err := handler(req)
	assert.ErrorIs(t, err, ErrorRouteNotFound)
}

func TestMiddleware(t *testing.T) {
	tree, _ := CreateTree()

	tree.Use(func(r *http.Request) error {
		if r.Proto != "HTTP/1.1" {
			return errors.New("protocol mismatch")
		}
		return nil
	})

	tree.Use(func(r *http.Request) error {
		if r.Header.Get("X-Content-Type-Options") != "JSONP" {
			return errors.New("codec error")
		}
		return nil
	})

	req, rer := http.NewRequest("PUT", "/api/v1/users", bytes.NewReader([]byte(`{"name":"Hessam","age":42}`)))
	if rer != nil {
		t.Fatalf("with error: %v", rer)
	}
	req.Header.Add("X-Content-Type-Options", "JSONP")

	err := tree.Serve(req)
	t.Logf("error is: %v", err)
}

func TestMountTree(t *testing.T) {
	tree1, _ := CreateTree()

	tree2 := NewMux()
	if err := tree2.AddHandler("GET", "v2/users", func(r *http.Request) error {
		log.Println("Hello from v2/users")
		return nil
	}); err != nil {
		t.Fatalf("add route handler: %v", err)
	}
	if err := tree2.AddHandler("*", "v2/products", nil); err != nil {
		t.Fatalf("add route handler: %v", err)
	}
	if err := tree2.AddHandler("*", "v2/comments", nil); err != nil {
		t.Fatalf("add route handler: %v", err)
	}

	tree1.Mount("/api", tree2)

	req, rer := http.NewRequest("PUT", "/api/v2/users", bytes.NewReader([]byte(`{"name":"Hessam","age":42}`)))
	if rer != nil {
		t.Fatalf("with error: %v", rer)
	}
	req.Header.Add("X-Content-Type-Options", "JSONP")

	err := tree1.Serve(req)
	t.Logf("error is: %v", err)
}

func CreateTree() (*MuxTree, error) {
	tree := NewMux()

	var v1 Handler = func(r *http.Request) error { return nil }
	if err := tree.AddHandler("GET", "/", v1); err != nil {
		return nil, fmt.Errorf("cant create handler: %v", err)
	}

	if err := tree.AddHandler("GET", "/api/v1/users", v1); err != nil {
		return nil, fmt.Errorf("cant create handler %v", err)
	}

	if err := tree.AddHandler("POST", "/api/v1/users", v1); err != nil {
		return nil, fmt.Errorf("cant create handler %v", err)
	}

	if err := tree.AddHandler("PUT", "/api/v1/users", func(r *http.Request) error {
		user := &User{}

		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			return err
		}
		log.Printf("user request: %+v", user)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("cant create handler %v", err)
	}

	if err := tree.AddHandler("GET", "/api/v1/products", v1); err != nil {
		return nil, fmt.Errorf("cant create handler %v", err)
	}

	if err := tree.AddHandler("POST", "/api/v1/products", func(r *http.Request) error {
		log.Printf("incoming req path: %v", r.URL.Path)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("cant create handler %v", err)
	}

	return tree, nil
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func BenchmarkTree(b *testing.B) {
	tree, _ := CreateTree()

	tree.Use(func(r *http.Request) error {
		if r.Proto != "HTTP1.1" {
			return errors.New("protocol mismatch")
		}
		return nil
	})

	tree.Use(func(r *http.Request) error {
		if r.Header.Get("X-Content-Type-Options") != "JSONP" {
			return errors.New("codec error")
		}
		return nil
	})

	tree.DisablePath("/api/v1/users")

	tree2 := NewMux()
	if err := tree2.AddHandler("GET", "v2/users", func(r *http.Request) error {
		log.Println("Hello from v2/users")
		return nil
	}); err != nil {
		b.Fatalf("add route handler: %v", err)
	}
	if err := tree2.AddHandler("*", "v2/products", nil); err != nil {
		b.Fatalf("add route handler: %v", err)
	}
	if err := tree2.AddHandler("*", "v2/comments", nil); err != nil {
		b.Fatalf("add route handler: %v", err)
	}

	tree.Mount("/api", tree2)

	req, rer := http.NewRequest("GET", "/api/v2/users", nil)
	if rer != nil {
		b.Fatalf("with error: %v", rer)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = tree.Serve(req)
	}
}
