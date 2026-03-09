package route

import "net/http"

type Route struct {
	Path    string
	Handler http.Handler
}
