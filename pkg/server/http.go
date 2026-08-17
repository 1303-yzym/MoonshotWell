package server

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// RunServer 启动http服务
func RunServer(host string, port int, handler http.Handler, tlsConfig *tls.Config) *http.Server {
	address := host + ":" + strconv.Itoa(port)
	srv := initServer(address, handler, tlsConfig)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				zap.L().Panic(err.Error())
			}

			return
		}
	}()

	zap.L().Info("http server listening",
		zap.String("Addr", srv.Addr),
		zap.String("address", "http://"+address),
	)

	return srv
}

// 自定义HTTP配置
func initServer(address string, handler http.Handler, tlsConfig *tls.Config) *http.Server {
	return &http.Server{
		TLSConfig:      tlsConfig,
		Addr:           address,
		Handler:        handler,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   0, // 禁用WriteTimeout，支持SSE长连接
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
