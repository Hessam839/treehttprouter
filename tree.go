package treehttprouter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
)

var (
	ErrorRouteNotFound = errors.New("route not found")
)

type ServeHTTP interface {
	Serve(conn net.Conn) error
}

type MuxTree struct {
	root *node

	routeNotFound *node

	methodNotAllowed *node

	middlewares []*Handler
}

//var Tree *MuxTree

func NewMux() *MuxTree {
	var rnf Handler = func(ctx *Context) error {
		return ErrorRouteNotFound
	}

	n := newNode("*")
	_ = n.addRoute("*", &rnf)

	var mnf Handler = func(ctx *Context) error {
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

func (t *MuxTree) Serve(c net.Conn) error {
	buff := make([]byte, 1024)
	readLen, rer := c.Read(buff)
	if rer != nil {
		return fmt.Errorf("read from connection failed: %v", rer)
	}

	req, qer := http.ReadRequest(bufio.NewReader(bytes.NewReader(buff[:readLen])))
	if qer != nil {
		return fmt.Errorf("read http 1.1 request failed: %v", qer)
	}

	handler := t.match(req)
	ctx, err := NewCtx(c)
	if err != nil {
		return fmt.Errorf("read request failed: %v", err)
	}

	for _, middleware := range t.middlewares {
		h := *middleware

		if err := h(ctx); err != nil {
			return err
		}
	}
	return handler(ctx)
}

func (t *MuxTree) Mount(path string, tree *MuxTree) {
	node := t.root.search(path)
	node.addChild(tree.root)
}
