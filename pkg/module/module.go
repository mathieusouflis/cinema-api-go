package module

import (
	"fmt"
	"net/http"
	"strings"

	"example.com/filmserver/pkg/route"
)

type Module struct {
	BasePath string
	Routes   []route.Route
	Server   *http.ServeMux
}

func New(basePath string, server *http.ServeMux) Module {
	return Module{BasePath: basePath, Routes: []route.Route{}, Server: server}
}

func (m *Module) RegisterRoute(operation Operation, path string, handler http.Handler) {
	m.Routes = append(m.Routes, route.Route{Path: operation.String() + " " + m.BasePath + path, Handler: handler})
}

func (m Module) Register() {
	for _, r := range m.Routes {
		m.Server.Handle(r.Path, r.Handler)
	}
}

func (m Module) PrintRoutesDocumentation() {
	fmt.Printf("\n\n========== %s ==========\n", m.BasePath)
	for _, r := range m.Routes {
		splittedRoute := strings.Split(r.Path, "/")
		fmt.Printf("%s\n", string(splittedRoute[0])+"/"+strings.Join(splittedRoute[2:], "/"))
	}
}
