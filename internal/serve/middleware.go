package serve

import "net/http"

type Middleware interface {
	Next(http.Handler) http.Handler
}
