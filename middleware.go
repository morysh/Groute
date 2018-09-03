package groute

import (
	"net/http"
)

type middleware = func(next http.Handler) http.Handler