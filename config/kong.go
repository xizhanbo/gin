package config

type Kong struct {
	Protocol      string `mapstructure:"protocol" json:"protocol" yaml:"protocol"`
	Host          string `mapstructure:"host" json:"host" yaml:"host"`
	Port          string `mapstructure:"port" json:"port" yaml:"port"`
	UpstreamsPath string `mapstructure:"upstreams_path" json:"upstreams_path" yaml:"upstreams_path"`
	TargetsPath   string `mapstructure:"targets_path" json:"targets_path" yaml:"targets_path"`
	ServicesPath  string `mapstructure:"services_path" json:"services_path" yaml:"services_path"`
	RoutesPath    string `mapstructure:"routes_path" json:"routes_path" yaml:"routes_path"`
}
