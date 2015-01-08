package cji

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)

type Cji struct {
	middlewares []interface{}
}

func Use(middlewares ...interface{}) *Cji {
	return (&Cji{}).Use(middlewares...)
}

func (z *Cji) Use(middlewares ...interface{}) *Cji {
	c := &Cji{z.middlewares}
	for _, mw := range middlewares {
		switch t := mw.(type) {
		default:
			panic(fmt.Sprintf("unsupported middleware type: %T", t))
		case func(http.Handler) http.Handler:
		case func(*web.C, http.Handler) http.Handler:
		}
		c.middlewares = append(c.middlewares, mw)
	}
	return c
}

// Compose together the middleware chain and wrap the handler with it
func (z *Cji) On(handler interface{}) web.HandlerFunc {
	var hfn web.HandlerFunc
	switch t := handler.(type) {
	case web.Handler:
		hfn = t.ServeHTTPC
	case func(web.C, http.ResponseWriter, *http.Request): // web.HandlerFunc
		hfn = t
	case func(http.ResponseWriter, *http.Request): // http.HandlerFunc
		hfn = func(c web.C, w http.ResponseWriter, r *http.Request) {
			t(w, r)
		}
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
func wrap(middleware interface{}) func(web.HandlerFunc) web.HandlerFunc {
	fn := func(hf web.HandlerFunc) web.HandlerFunc {
		return func(c web.C, w http.ResponseWriter, r *http.Request) {
			newFn := func(ww http.ResponseWriter, rr *http.Request) {
				hf(c, ww, rr)
			}

			var fn http.HandlerFunc
			switch mw := middleware.(type) {
			default:
				panic(fmt.Sprintf("unsupported middleware type: %T", mw))
			case func(http.Handler) http.Handler:
				fn = mw(http.HandlerFunc(newFn)).(http.HandlerFunc)
			case func(*web.C, http.Handler) http.Handler:
				fn = mw(&c, http.HandlerFunc(newFn)).(http.HandlerFunc)
			}
			fn(w, r)
		}
	}
	return fn
}
