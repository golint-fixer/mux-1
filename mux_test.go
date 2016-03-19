package mux

import (
	"github.com/nbio/st"
	"gopkg.in/vinci-proxy/vinci.v0/utils"
	"net/http"
	"testing"
)

func TestMuxSimple(t *testing.T) {
	mx := New()
	mx.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("foo", "bar")
		h.ServeHTTP(w, r)
	})
	wrt := utils.NewWriterStub()
	req := newRequest()

	mx.Middleware.Run("request", wrt, req, nil)
	st.Expect(t, wrt.Header().Get("foo"), "bar")
}
