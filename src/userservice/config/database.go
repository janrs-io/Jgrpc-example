package config

// Database Database config
type Database struct {
	Driver              string `json:"driver" yaml:"driver"`
	Host                string `json:"host" yaml:"host"`
	Port                int    `json:"port" yaml:"port"`
	UserName            string `json:"username" yaml:"username"`
	Password            string `json:"password" yaml:"password"`
	Database            string `json:"database" yaml:"database"`
	Charset             string `json:"charset" yaml:"charset"`
	MaxIdleCons         int    `json:"maxIdleCons" yaml:"maxIdleCons"`
	MaxOpenCons         int    `json:"maxOpenCons" yaml:"maxOpenCons"`
	LogMode             string `json:"logMode" yaml:"logMode"`
	EnableFileLogWriter bool   `json:"enableFileLogWriter" yaml:"enableFileLogWriter"`
	LogFilename         string `json:"logFilename" yaml:"logFilename"`
}
