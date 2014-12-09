package cji

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

var defaultErrResponder ErrResponderFunc = func(c web.C, w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

type HandlerFunc func(web.C, http.ResponseWriter, *http.Request, func(int, error))

type ErrResponderFunc func(web.C, http.ResponseWriter, *http.Request, int, error)

type cji struct {
	handlers     []HandlerFunc
	errResponder ErrResponderFunc
}

func Link(handlers ...HandlerFunc) *cji {
	return &cji{handlers: handlers}
}

func SetErrResponder(errResponder ErrResponderFunc) {
	defaultErrResponder = errResponder
}

func (j *cji) To(handler web.HandlerFunc) web.HandlerFunc {
	return web.HandlerFunc(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var status int
		var err error
		for _, h := range j.handlers {
			h(c, w, r, func(rStatus int, rErr error) {
				status = rStatus
				err = rErr
			})
			if err != nil {
				var respondErr ErrResponderFunc = defaultErrResponder
				if j.errResponder != nil {
					respondErr = j.errResponder
				}
				respondErr(c, w, r, status, err)
				return
			}
		}
		handler(c, w, r)
	})
}

func (j *cji) Link(handlers ...HandlerFunc) *cji {
	jj := &cji{handlers: j.handlers, errResponder: j.errResponder}
	jj.handlers = append(jj.handlers, handlers...)
	return jj
}

func (j *cji) WithErrResponder(errResponder ErrResponderFunc) *cji {
	j.errResponder = errResponder
	return j
}
