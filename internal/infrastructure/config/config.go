package config

import "github.com/1303-yzym/MoonshotWell/pkg/config"

var cfg = config.New[ConfigInstance]()

type Config config.Config[ConfigInstance]

type ConfigInstance struct {
	// TODO 定义扩展
	// Constant Constant `mapstructure:"constant" json:"constant" comment:"常量"`
}

// 扩展示例
//type Constant struct {
//	LoginSmsTemplateCode    string `mapstructure:"login_sms_template_code" json:"login_sms_template_code" comment:"登录短信模板code"`
//	LoginSmsSign            string `mapstructure:"login_sms_sign" json:"login_sms_sign" comment:"登录短信模板签名"`
//	CDN                     string `mapstructure:"cdn" json:"cdn" comment:"cdn加速地址"`
//}

func InitConfig(path string) error {
	return cfg.Init(path)
}

func Load() *Config {
	return (*Config)(cfg.Load())
}
