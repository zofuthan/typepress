package global

import (
	"net/http"
)

var routeBefore = []func(http.ResponseWriter, *http.Request) bool{}

// register function on before HandlerFunc of route executed
func OnRouteBefore(f func(http.ResponseWriter, *http.Request) bool) {
	routeBefore = append(routeBefore, f)
}

// fire functions on before HandlerFunc of route executed
func FireRouteBefore(wr http.ResponseWriter, r *http.Request) bool {
	for _, f := range routeBefore {
		if !f(wr, r) {
			return false
		}
	}
	return true
}

var routeAfter = []func(http.ResponseWriter, *http.Request) bool{}

// register function on after HandlerFunc of route executed
func OnRouteAfter(f func(http.ResponseWriter, *http.Request) bool) {
	routeAfter = append(routeAfter, f)
}

// fire functions on after HandlerFunc of route executed
func FireRouteAfter(wr http.ResponseWriter, r *http.Request) bool {
	for _, f := range routeAfter {
		if !f(wr, r) {
			return false
		}
	}
	return true
}

var renderAfter = []func(*http.Request, error) bool{}

// register function on after rander for template
func OnRenderAfter(f func(*http.Request, error) bool) {
	renderAfter = append(renderAfter, f)
}

// fire functions on after rander for template
func FireRenderAfter(r *http.Request, err error) {
	for _, f := range renderAfter {
		if !f(r, err) {
			return
		}
	}
	return
}

var muxBefore = []func(http.ResponseWriter, *http.Request) bool{}

// register function on before Mux.ServeHTTP executed
func OnMuxBefore(f func(http.ResponseWriter, *http.Request) bool) {
	muxBefore = append(muxBefore, f)
}

// fire functions on before Mux.ServeHTTP executed
func FireMuxBefore(wr http.ResponseWriter, r *http.Request) bool {
	for _, f := range muxBefore {
		if !f(wr, r) {
			return false
		}
	}
	return true
}

var muxAfter = []func(http.ResponseWriter, *http.Request) bool{}

// register function on after Mux.ServeHTTP executed
func OnMuxAfter(f func(http.ResponseWriter, *http.Request) bool) {
	muxAfter = append(muxAfter, f)
}

// fire functions on after Mux.ServeHTTP executed
func FireMuxAfter(wr http.ResponseWriter, r *http.Request) bool {
	for _, f := range muxAfter {
		if !f(wr, r) {
			return false
		}
	}
	return true
}

var events = []func(code int, r *http.Request, i ...interface{}) bool{}

// register function on events
func OnEvent(f func(code int, r *http.Request, i ...interface{}) bool) {
	events = append(events, f)
}

// fire functions on events
func FireEvent(code int, r *http.Request, i ...interface{}) bool {
	for _, f := range events {
		if !f(code, r, i...) {
			return false
		}
	}
	if code != 200 && len(events) == 0 {
		LogDebug(code, r, i...)
	}
	return true
}

var onShutDown []func()

// register function on shutDown
func OnShutDown(f func()) {
	onShutDown = append(onShutDown, f)
}

// fire functions on shutDown
func FireShutDown() {
	for _, f := range onShutDown {
		f()
	}
	onShutDown = []func(){}
}
