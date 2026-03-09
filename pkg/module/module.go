package module

import (
	"net/http"

	"example.com/filmserver/pkg/route"
)

type Module struct {
	basePath string
	routes   []route.Route
	server   *http.ServeMux
}

func New(m Module) Module {
	module := Module{
		basePath: m.basePath,
		routes:   m.routes,
		server:   m.server,
	}
	return module
}

func (m Module) Register() {
	for _, route := range m.routes {
		m.server.Handle(m.basePath, route.Handler())
	}
}

func (m Module) GetRoutes() []string {
	routes := []string{}
	for _, route := range m.routes {
		routes = append(routes, m.basePath+route.Path())
	}
	return routes
}
