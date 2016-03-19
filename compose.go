package mux

import "net/http"

// If creates a new multiplexer that will be executed if all the mux matchers passes.
func If(muxes ...*Layer) *Layer {
	mx := New()
	for _, mm := range muxes {
		mx.AddMatcher(mm.Matchers...)
	}
	return mx
}

// Every is an alias to If().
func Every(muxes ...*Layer) *Layer {
	return If(muxes...)
}

// Or creates a new multiplexer that will be executed if at least one mux matcher passes.
func Or(muxes ...*Layer) *Layer {
	return Match(func(req *http.Request) bool {
		for _, mm := range muxes {
			if mm.Match(req) {
				return true
			}
		}
		return false
	})
}

// Some is an alias to Or().
func Some(muxes ...*Layer) *Layer {
	return Or(muxes...)
}
