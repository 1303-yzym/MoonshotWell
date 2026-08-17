package warp

import (
	"context"
	"crypto/tls"
	"net/http"
	"reflect"
	"time"
)

type HandlerFunc[O context.Context, C Context[O]] func(ctx C)

type Route[O context.Context, C Context[O]] struct {
	*MetaInfo
	HandlerValue *reflect.Value
	Handlers     []HandlerFunc[O, C]
}

type RouteInfo struct {
	*MetaInfo
	HandlerValue *reflect.Value
}

type Warp[O context.Context, C Context[O]] struct {
	RouterGroup[O, C]

	adapter Adapter[O, C]

	config     *Config
	routeInfos []*RouteInfo
}

var _ RouterGroup[context.Context, Context[context.Context]] = (*Warp[context.Context, Context[context.Context]])(nil)

func New[O context.Context, C Context[O]](adapter Adapter[O, C], opts ...Options) *Warp[O, C] {
	cfg := &Config{
		NoRouterInfo: true,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	warp := &Warp[O, C]{
		adapter:    adapter,
		config:     cfg,
		routeInfos: make([]*RouteInfo, 0),
	}

	warp.RouterGroup = &routerGroup[O, C]{
		rg:       warp,
		basePath: "/",
	}

	return warp
}

func (wp *Warp[O, C]) Handler() http.Handler {
	return wp.adapter.Handler()
}

func (wp *Warp[O, C]) Server(host, port string, tlsConfig *tls.Config) *http.Server {
	return &http.Server{
		TLSConfig:      tlsConfig,
		Addr:           host + ":" + port,
		Handler:        wp.Handler(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (wp *Warp[O, C]) RouteInfos() []*RouteInfo {
	return wp.routeInfos
}

func (wp *Warp[O, C]) registerHandler(route *Route[O, C], _ ...Middleware[O, C]) {
	wp.adapter.RegisterHandler(route)
	wp.addRouteInfo(route)
}

func (wp *Warp[O, C]) addRouteInfo(route *Route[O, C]) {
	if !wp.config.NoRouterInfo {
		wp.routeInfos = append(wp.routeInfos, &RouteInfo{
			MetaInfo:     route.MetaInfo,
			HandlerValue: route.HandlerValue,
		})
	}
}
