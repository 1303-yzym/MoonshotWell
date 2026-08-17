package warp

type Config struct {
	NoRouterInfo bool
}

type Options func(*Config)

func WithNoRouterInfo(noRouterInfo bool) Options {
	return func(cfg *Config) {
		cfg.NoRouterInfo = noRouterInfo
	}
}
