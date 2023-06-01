package config

// Client Grpc 客户端配置
type Client struct {
	// product 产品服务客户端配置
	ProductHost string `json:"productHost" yaml:"productHost"`
	ProductPort string `json:"productPort" yaml:"productPort"`
}
