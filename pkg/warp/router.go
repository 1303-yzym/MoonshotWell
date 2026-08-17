package warp

import (
	"context"
	"path"
)

type MiddlewareChain[O context.Context, C Context[O]] []HandlerFunc[O, C]

type RouterGroup[O context.Context, C Context[O]] interface {
	registerHandler(route *Route[O, C], middlewares ...Middleware[O, C])
	BasePath() string
	Group(relativePath string, middlewares ...Middleware[O, C]) RouterGroup[O, C]
	Use(handlers ...Middleware[O, C]) RouterGroup[O, C]
	Handle(httpMethod, relativePath string, handler HandlerFunc[O, C], middlewares ...Middleware[O, C])
}

type routerGroup[O context.Context, C Context[O]] struct {
	rg              RouterGroup[O, C]
	basePath        string
	middlewares     []string
	middlewareChain []HandlerFunc[O, C]
}

var _ RouterGroup[context.Context, Context[context.Context]] = (*routerGroup[context.Context, Context[context.Context]])(nil)

func (r *routerGroup[O, C]) BasePath() string {
	return r.basePath
}

type Middleware[O context.Context, C Context[O]] struct {
	Name    string
	Handler HandlerFunc[O, C]
}

func WrapMiddleware[O context.Context, C Context[O]](name string, handler HandlerFunc[O, C]) Middleware[O, C] {
	return Middleware[O, C]{Name: name, Handler: handler}
}

func (r *routerGroup[O, C]) Use(middlewares ...Middleware[O, C]) RouterGroup[O, C] {
	for _, middleware := range middlewares {
		r.middlewareChain = append(r.middlewareChain, middleware.Handler)
		r.middlewares = append(r.middlewares, middleware.Name)
	}

	return r
}

func (r *routerGroup[O, C]) Group(relativePath string, middlewares ...Middleware[O, C]) RouterGroup[O, C] {
	handlers := make([]HandlerFunc[O, C], 0, len(middlewares))

	middlewareNames := make([]string, 0, len(middlewares))
	for _, middleware := range middlewares {
		handlers = append(handlers, middleware.Handler)
		middlewareNames = append(middlewareNames, middleware.Name)
	}

	return &routerGroup[O, C]{
		rg:              r,
		basePath:        relativePath,
		middlewares:     middlewareNames,
		middlewareChain: handlers,
	}
}

func (r *routerGroup[O, C]) Handle(httpMethod, relativePath string, handler HandlerFunc[O, C], middlewares ...Middleware[O, C]) {
	r.registerHandler(&Route[O, C]{
		MetaInfo: &MetaInfo{
			Path:   relativePath,
			Method: httpMethod,
		},
		Handlers: []HandlerFunc[O, C]{handler},
	}, middlewares...)
}

func (r *routerGroup[O, C]) registerHandler(route *Route[O, C], middlewares ...Middleware[O, C]) {
	if r.rg != nil {
		route.Path = r.calculateAbsolutePath(route.Path)
		route.Middlewares = append(r.middlewares, route.Middlewares...)

		var middlewareHandlers []HandlerFunc[O, C]

		for _, middleware := range middlewares {
			route.Middlewares = append(route.Middlewares, middleware.Name)
			middlewareHandlers = append(middlewareHandlers, middleware.Handler)
		}

		route.Handlers = r.combineHandlers(append(middlewareHandlers, route.Handlers...)...)

		r.rg.registerHandler(route)
	}
}

func (r *routerGroup[O, C]) calculateAbsolutePath(relativePath string) string {
	if relativePath == "" {
		return r.basePath
	}

	finalPath := path.Join(r.basePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}

	return finalPath
}

func (r *routerGroup[O, C]) combineHandlers(middlewareChain ...HandlerFunc[O, C]) []HandlerFunc[O, C] {
	finalSize := len(r.middlewareChain) + len(middlewareChain)
	mergedHandlers := make(MiddlewareChain[O, C], finalSize)
	copy(mergedHandlers, r.middlewareChain)
	copy(mergedHandlers[len(r.middlewareChain):], middlewareChain)

	return mergedHandlers
}
