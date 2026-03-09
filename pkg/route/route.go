package route

import "net/http"

type Route struct {
	handler http.Handler
	path    string
}

func New(r Route) Route {
	return Route{
		handler: r.handler,
		path:    r.path,
	}
}
func (r Route) Path() string {
	return r.path
}

func (r Route) Handler() http.Handler {
	return r.handler
}
