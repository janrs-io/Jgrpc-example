package config

// Trace Trace Config
type Trace struct {
	TracerName  string `json:"tracerName" yaml:"tracerName"`
	ServiceName string `json:"serviceName" yaml:"serviceName"`
	EndPoint    string `json:"endPoint" yaml:"endPoint"`
}
