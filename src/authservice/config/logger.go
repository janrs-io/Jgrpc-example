package config

// Logger Grpc logger middleware config
type Logger struct {
	Path      string `json:"path" yaml:"path"`
	MaxSize   int    `json:"maxSize" yaml:"maxSize"`
	LocalTime bool   `json:"localTime" yaml:"localTime"`
}
