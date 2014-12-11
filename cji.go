package cji

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)

// Goji's middleware type signature
type MiddlewareFunc func(*web.C, http.Handler) http.Handler

type cji struct {
	middlewares []MiddlewareFunc
}

func Use(middlewares ...MiddlewareFunc) *cji {
	return &cji{middlewares}
}

func (z *cji) Use(middlewares ...MiddlewareFunc) *cji {
	mw := z.middlewares
	mw = append(mw, middlewares...)
	return &cji{mw}
}

// Compose together the middleware chain and wrap the handler with it
func (z *cji) On(handler interface{}) web.HandlerFunc {
	var hfn web.HandlerFunc
	switch t := handler.(type) {
	case web.Handler:
		hfn = t.ServeHTTPC
	case func(web.C, http.ResponseWriter, *http.Request): // web.HandlerFunc
		hfn = t
	default:
		panic(fmt.Sprintf("unsupported handler type: %T", t))
	}

	if len(z.middlewares) == 0 {
		return hfn
	}

	m := wrap(z.middlewares[len(z.middlewares)-1])(hfn)
	for i := len(z.middlewares) - 2; i >= 0; i-- {
		f := wrap(z.middlewares[i])
		m = f(m)
	}
	return m
}

// Wrap takes a middleware that works on http.Handler and returns a function that
// takes a web.HandlerFunc and returns a web.HandlerFunc. We use this to wrap HandlerFuncs
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
				panic("unsupported handler passed to the chain!")
			}
		}
	}
	return fn
}
