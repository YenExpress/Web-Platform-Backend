package utils

import (
	"net/http"
)

// Router is an interface that abstracts routing functionality
type Router interface {
	GET(pattern string, handler http.HandlerFunc)
	POST(pattern string, handler http.HandlerFunc)
	Group(prefix string, middleware ...Middleware) Router
	Use(middleware ...Middleware)
}

// Middleware represents a middleware function
type Middleware func(next http.HandlerFunc) http.HandlerFunc

type RouterFunc func(r Router)

// aggregate routes to be handled by router
func aggregateRouters(routeLists ...[]RouterFunc) []RouterFunc {
	totalLength := 0

	for _, routeList := range routeLists {
		totalLength += len(routeList)
	}

	aggregatedRoutes := make([]RouterFunc, 0, totalLength)

	for _, routeList := range routeLists {
		aggregatedRoutes = append(aggregatedRoutes, routeList...)
	}

	return aggregatedRoutes
}

// function to load all routes aggregrated
func LoadRoutes(r Router, Routers ...[]RouterFunc) {
	allRouters := aggregateRouters(Routers...)

	for _, route := range allRouters {
		route(r)
	}

}
