// Package mux implements an HTTP domain-specific traffic multiplexer
// with built-in matchers and features for easy plugin composition and activable logic.
package mux

import (
	"gopkg.in/vinci-proxy/vinci.v0/middleware"
	"net/http"
)

// Layer is a HTTP request/response/error multiplexer who implements both
// middleware and plugin interfaces.
// It has been designed for easy plugin composition based on HTTP matchers/filters.
type Layer struct {
	// Matchers stores a list of matcher functions.
	Matchers []Matcher

	// Middleware stores the multiplexer middleware layer.
	Middleware *middleware.Layer
}

// New creates a new multiplexer with default settings.
func New() *Layer {
	return &Layer{Middleware: middleware.New()}
}

// Match matches the give Context againts a list of matchers and
// returns `true` if all the matchers passed.
func (m *Layer) Match(req *http.Request) bool {
	for _, matcher := range m.Matchers {
		if !matcher(req) {
			return false
		}
	}
	return true
}

// AddMatcher adds a new matcher function in the current mumultiplexer matchers stack.
func (m *Layer) AddMatcher(matchers ...Matcher) *Layer {
	m.Matchers = append(m.Matchers, matchers...)
	return m
}

// If is a semantic alias to AddMatcher.
func (m *Layer) If(matchers ...Matcher) *Layer {
	return m.AddMatcher(matchers...)
}

// Some matches the incoming request if at least one of the matchers passes.
func (m *Layer) Some(matchers ...Matcher) *Layer {
	return m.AddMatcher(func(req *http.Request) bool {
		for _, matcher := range matchers {
			if matcher(req) {
				return true
			}
		}
		return false
	})
}

// Use registers a new plugin in the middleware stack.
func (m *Layer) Use(handler interface{}) *Layer {
	m.Middleware.Use(handler)
	return m
}

// UsePhase registers a new plugin in the middleware stack.
func (m *Layer) UsePhase(phase string, handler interface{}) *Layer {
	m.Middleware.Use(phase, handler)
	return m
}

// UseFinalHandler registers a new plugin in the middleware stack.
func (m *Layer) UseFinalHandler(handler http.Handler) *Layer {
	m.Middleware.UseFinalHandler(handler)
	return m
}

// HandleHTTP returns the function handler to match an incoming HTTP transacion
// and trigger the equivalent middleware phase.
func (m *Layer) HandleHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	if m.Match(r) {
		m.Middleware.Run("request", w, r, h)
		return
	}
	h.ServeHTTP(w, r)
}
