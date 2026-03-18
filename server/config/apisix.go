package config

type Apisix struct {
	Url    string `mapstructure:"url" json:"url" yaml:"url"`
	ApiKey string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`
}
