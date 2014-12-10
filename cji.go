/*

*/

package cji

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// Goji's middleware type signature
type MiddlewareFunc func(*web.C, http.Handler) http.Handler

type cji struct {
	middlewares []MiddlewareFunc
}

func Use(middlewares ...MiddlewareFunc) *cji {
	return &cji{middlewares: middlewares}
}

//Compose together the middleware chain and wrap the handler with it
func (j *cji) On(handler web.HandlerFunc) web.HandlerFunc {
	//if len(j.middlewares) == 0 {
	//TODO handle error better
	//panic()
	////	} else {
	m := (wrap(j.middlewares[0]))(handler)
	for i := len(j.middlewares) - 2; i >= 0; i-- {
		f := wrap(j.middlewares[i])
		m = f(m)
	}
	return m
}

// wrap takes a middleware that works on http.Handler and returns a function that takes a web.HandlerFunc and returns a web.HandlerFunc. We use this to wrap HandlerFuncs with
func wrap(middleware MiddlewareFunc) func(web.HandlerFunc) web.HandlerFunc {
	fn := func(hf web.HandlerFunc) web.HandlerFunc {
		return func(c web.C, w http.ResponseWriter, r *http.Request) {
			newFn := func(ww http.ResponseWriter, rr *http.Request) {
				hf(c, ww, rr)
			}
			fn, ok := middleware(&c, http.HandlerFunc(newFn)).(http.HandlerFunc)
			if ok {
				fn(w, r)
			} else {
				// something went wrong!
			}
		}
	}
	return fn
}
