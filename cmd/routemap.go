package main

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/go-chi/chi"
)

type endpoint struct {
	Method      string   `json:"method"`
	Route       string   `json:"route"`
	Handler     string   `json:"handler"`
	Middlewares []string `json:"middleware"`
}

func nameFn(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func nameFns(fns []func(http.Handler) http.Handler) (names []string) {
	for _, fn := range fns {
		names = append(names, nameFn(fn))
	}
	return
}

func maproutes(r chi.Routes) (endpoints []endpoint) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		e := endpoint{
			Method:      method,
			Route:       route,
			Handler:     nameFn(handler),
			Middlewares: nameFns(middlewares),
		}
		endpoints = append(endpoints, e)
		return nil
	}
	chi.Walk(r, walkFunc)
	return
}
