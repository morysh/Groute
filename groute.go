package groute

import (
	"net/http"
	"regexp"
)

type handlerFunc = func(http.ResponseWriter, *http.Request)

var (
	pathVariablePattern = regexp.MustCompile(`\{[a-zA-Z]+}`)
	notFound = http.NotFound
)

type Router interface {
	http.Handler
	Handle(method string, path string, handler http.Handler)
	HandleFunc(method string, path string, handler handlerFunc)
}

type router struct {
	root tree
}