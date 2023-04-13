package config

// Client Client service
type Client struct {
	AuthHost string `json:"authHost" yaml:"authHost"`
	AuthPort string `json:"authPort" yaml:"authPort"`
}
