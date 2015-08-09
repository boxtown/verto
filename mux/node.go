package mux

import (
	"net/http"
)

// -------------------------------------
// ---------- Path Muxer Nodes ----------

// Node is an interface for endpoint nodes that
// allows for the addition of per-route plugin
// handlers.
type Node interface {
	// Use adds a PluginHandler onto the end of the chain of plugins
	// for a node.
	Use(handler PluginHandler) Node

	// UseHandler wraps the handler as a PluginHandler and adds it onto the end
	// of the plugin chain.
	UseHandler(hander http.Handler) Node
}

// muxNode is a private struct used to keep track of handlers
// and plugins per method+path.
type muxNode struct {
	mux *PathMuxer

	handlers map[string]http.Handler
	chains   map[string]*plugins

	path string
}

// newMuxNode returns a pointer to a newly initialized muxNode.
func newMuxNode(mux *PathMuxer, path string) *muxNode {
	return &muxNode{
		mux:      mux,
		handlers: make(map[string]http.Handler),
		chains:   make(map[string]*plugins),
		path:     path,
	}
}

// ServeHTTP delegates to the appropriate handler based on
// the request method or calls the Not Implemented handler if
// the desired method handler does not exist. If there is a method-appropriate
// chain of plugins, those will be run first.
func (node *muxNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := node.handlers[r.Method]
	if !ok {
		node.mux.NotImplemented.ServeHTTP(w, r)
		return
	}

	chain, _ := node.chains[r.Method]
	if chain == nil || chain.length == 0 {
		handler.ServeHTTP(w, r)
	} else {
		chain.run(w, r)
	}
}

// Private implementation of Node that can map plugin chains
// to specific methods.
type nodeImpl struct {
	method   string
	handlers map[string]http.Handler
	chains   map[string]*plugins
}

func newNodeImpl(method string, node *muxNode) *nodeImpl {
	return &nodeImpl{
		method:   method,
		handlers: node.handlers,
		chains:   node.chains,
	}
}

// Use adds a PluginHandler onto the end of the chain of plugins
// for a node.
func (node *nodeImpl) Use(handler PluginHandler) Node {
	chain, ok := node.chains[node.method]
	if !ok {
		node.chains[node.method] = newPlugins()
		chain = node.chains[node.method]
	}

	// Since we always add node.handler as the last handler,
	// we have to pop it off first before adding the desired handler.
	if chain.length > 0 {
		chain.popTail()
	}

	chain.use(handler)
	chain.use(PluginFunc(
		func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			node.handlers[node.method].ServeHTTP(w, r)
		}))

	return node
}

// UseHandler wraps the handler as a PluginHandler and adds it onto the end
// of the plugin chain.
func (node *nodeImpl) UseHandler(handler http.Handler) Node {
	pluginHandler := PluginFunc(
		func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			handler.ServeHTTP(w, r)
			next(w, r)
		})

	node.Use(pluginHandler)
	return node
}
