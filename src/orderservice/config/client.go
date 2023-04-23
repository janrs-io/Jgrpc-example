package config

// Client Grpc 客户端配置
type Client struct {
	// auth 鉴权服务客户端配置
	AuthHost string `json:"authHost" yaml:"authHost"`
	AuthPort string `json:"authPort" yaml:"authPort"`
	// product 产品服务客户端配置
	ProductHost string `json:"productHost" yaml:"productHost"`
	ProductPort string `json:"productPort" yaml:"productPort"`
}
