package config

type MicroServ struct {
	Protocol string   `mapstructure:"protocol" json:"protocol" yaml:"protocol"`
	Host     []string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string   `mapstructure:"port" json:"port" yaml:"port"`
}
