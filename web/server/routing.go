package server

import (
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/files"
	"github.com/Galagoshin/GoUtils/requests"
	"net/http"
	"regexp"
)

type option string

type route struct {
	pattern *regexp.Regexp
	methods []requests.Method
	handle  func(r *requests.Request) (requests.Response, error)
	option  option
}

func (option option) len() int {
	return len(regexp.MustCompile("{[^/]+}").FindAllString(string(option), -1))
}

var routes = []route{}

func Route(pattern string, methods []requests.Method, handle func(*requests.Request) (requests.Response, error)) {
	compiled, opt := CompileRoute(pattern)
	logger.Debug(2, false, fmt.Sprintf("Registered route: %s", pattern))
	logger.Debug(3, false, fmt.Sprintf("Registered route regexp: %s", compiled))
	routes = append(routes, route{
		pattern: compiled,
		methods: methods,
		handle:  handle,
		option:  opt,
	})
}

func Static() {
	staticDir := files.Directory{Path: "static"}
	err := staticDir.CreateAll()
	if err != nil {
		logger.Panic(err)
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
