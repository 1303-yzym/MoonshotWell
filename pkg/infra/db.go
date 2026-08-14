package infra

type DBConfig struct {
	MySQL MySQLConfig `json:"mysql" comment:"mysql配置"`
}
type MySQLConfig struct {
	Port int `json:"port" comment:"端口"`
	// TODO
}
