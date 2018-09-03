package groute

import (
	"net/http"
)

type methodHandlers map[string]http.Handler

func (mh methodHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := mh[r.Method]
	if ok {
		handler.ServeHTTP(w, r)
	} else if len(mh) == 0 {
		notFound(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		for method := range mh {
			w.Header().Add("Allow", method)
		}
	}
}

type tree interface {
	Handle(path []string, w http.ResponseWriter, r *http.Request)
	Add(path []string, method string, handler handlerFunc)
}

type genericTree struct {
	methodHandlers
	children map[string]tree
}

type pathParamTree struct {
	methodHandlers
	child tree
	name string
}

func (t genericTree) Handle(path []string, w http.ResponseWriter, r *http.Request) {
	if len(path) == 0 {
		t.ServeHTTP(w, r)
	} else {
		next, ok := t.children[path[0]]
		if ok {
			next.Handle(path[1:], w, r)
		} else {
			notFound(w, r)
		}
	}
}

func (t pathParamTree) Handle(path []string, w http.ResponseWriter, r *http.Request) {
	if len(path) == 0 {
		t.ServeHTTP(w, r)
	} else {
		if t.child == nil {
			notFound(w, r)
		} else {
			value := path[0]
			r.URL.Query().Set(t.name, value)
			t.child.Handle(path[1:], w, r)
		}
	}
}