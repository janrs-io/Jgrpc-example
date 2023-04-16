package config

// Client Grpc 客户端配置
type Client struct {
	// auth 服务客户端配置
	AuthHost string `json:"authHost" yaml:"authHost"`
	AuthPort string `json:"authPort" yaml:"authPort"`
}
