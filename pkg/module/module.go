package router

import (
	"fmt"
	"net/http"
	"strings"

	"example.com/filmserver/pkg/route"
)

type Router struct {
	BasePath string
	Routes   []route.Route
	Server   *http.ServeMux
}

func New(basePath string, server *http.ServeMux) Router {
	return Router{BasePath: basePath, Routes: []route.Route{}, Server: server}
}

func (m *Router) RegisterRoute(operation Operation, path string, handler http.Handler) {
	m.Routes = append(m.Routes, route.Route{Path: operation.String() + " " + m.BasePath + path, Handler: handler})
}

func (m Router) Register() {
	for _, r := range m.Routes {
		m.Server.Handle(r.Path, r.Handler)
	}
}

func (m Router) PrintRoutesDocumentation() {
	fmt.Printf("\n\n========== %s ==========\n", m.BasePath)
	for _, r := range m.Routes {
		splittedRoute := strings.Split(r.Path, "/")
		fmt.Printf("%s\n", string(splittedRoute[0])+"/"+strings.Join(splittedRoute[2:], "/"))
	}
}
