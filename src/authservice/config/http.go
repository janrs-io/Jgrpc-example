package config

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Http Http server config
type Http struct {
	Host      string `json:"host" yaml:"host"`
	Port      string `json:"port" yaml:"port"`
	Name      string `json:"name" yaml:"name"`
	Server    *http.Server
	ServerMux *runtime.ServeMux
}
