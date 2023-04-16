package config

// WhiteList 权限白名单。包含 api 接口白名单以及 rbac 权限白名单。
type WhiteList struct {
	Api        []any `json:"api" yaml:"api"`
	Permission []any `json:"permission" yaml:"permission"`
}
