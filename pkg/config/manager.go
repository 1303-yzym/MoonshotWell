package config

import (
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Manager 配置文件线程安全封装.
type Manager[E any] struct {
	vp *viper.Viper

	config   *Config[E]
	configMu sync.RWMutex
	once     sync.Once
}

func New[E any]() *Manager[E] {
	return &Manager[E]{
		vp:       viper.GetViper(),
		config:   nil,
		configMu: sync.RWMutex{},
		once:     sync.Once{},
	}
}

func (cm *Manager[E]) Init(configPath string) (err error) {
	cm.once.Do(func() {
		err = cm.init(configPath)
	})

	return err
}

// Load 加载配置文件
func (cm *Manager[E]) Load() *Config[E] {
	if cm.config == nil {
		zap.L().Fatal("configuration not initialized")
	}

	return cm.get()
}

func (cm *Manager[E]) init(cfgFilePath string) error {
	// 环境变量前缀
	cm.vp.SetEnvPrefix("MW")
	cm.vp.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	cm.vp.AutomaticEnv()

	// 设置配置文件路径
	if cfgFilePath != "" {
		cm.vp.SetConfigFile(cfgFilePath)
	} else {
		cm.vp.AddConfigPath(".")
		cm.vp.SetConfigName("config")
		cm.vp.SetConfigType("yaml")
	}

	// 读取配置文件
	if err := cm.vp.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read configuration file '%s': %w",
			cm.vp.ConfigFileUsed(), err)
	}

	zap.S().Infof("Config: [%s]", cm.vp.ConfigFileUsed())

	// 加载配置
	if err := cm.loadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// 设置文件监听
	cm.setupFileWatcher()

	return nil
}

func (cm *Manager[E]) loadConfig() error {
	cm.configMu.Lock()
	defer cm.configMu.Unlock()

	config := &Config[E]{}
	if err := cm.vp.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.Remain != nil {
		zap.L().Info("there are unmapped fields", zap.Any("other_config", config.Remain))
	}

	if err := config.Validation(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	cm.config = config

	return nil
}

func (cm *Manager[E]) setupFileWatcher() {
	cm.vp.WatchConfig()
	cm.vp.OnConfigChange(func(e fsnotify.Event) {
		zap.L().Info("config changed", zap.String("filename", e.Name))
		cm.reloadConfig()
	})
}

func (cm *Manager[E]) reloadConfig() {
	cm.configMu.Lock()
	defer cm.configMu.Unlock()

	if err := cm.vp.Unmarshal(&cm.config); err != nil {
		zap.L().Error("configuration reload failed", zap.Error(err))
	}

	// TODO 重载配置文件应该验证，并重启应该重启的服务

	zap.L().Info("configuration reloaded")
}

func (cm *Manager[E]) get() *Config[E] {
	cm.configMu.RLock()
	defer cm.configMu.RUnlock()

	return cm.config
}

// Validation 配置文件验证
func (c BasicConfig) Validation() error {
	// 验证环境是否正确
	switch c.Env {
	case EnvDev, EnvProd, EnvTest:
	default:
		return fmt.Errorf("invalid env: %s", c.Env)
	}

	// TODO 验证其他配置

	return nil
}

// IsDev 是否是开发环境
func (c BasicConfig) IsDev() bool {
	return c.Env == EnvDev
}

// IsProd 是否是开发环境
func (c BasicConfig) IsProd() bool {
	return c.Env == EnvProd
}

func (c BasicConfig) IsTest() bool {
	return c.Env == EnvTest
}
