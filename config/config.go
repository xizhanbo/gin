package config

type Configuration struct {
	App       App       `mapstructure:"app" json:"app" yaml:"app"`
	Log       Log       `mapstructure:"log" json:"log" yaml:"log"`
	Databases Databases `mapstructure:"databases" json:"databases" yaml:"databases"`
	Jwt       Jwt       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis     Redis     `mapstructure:"redis" json:"redis" yaml:"redis"`
	Storage   Storage   `mapstructure:"storage" json:"storage" yaml:"storage"`
	Kong      Kong      `mapstructure:"kong" json:"kong" yaml:"kong"`
	Consul    Consul    `mapstructure:"consul" json:"consul" yaml:"consul"`
	MicroServ MicroServ `mapstructure:"micro_serv" json:"micro_serv" yaml:"micro_serv"`
	Jaeger    Jaeger    `mapstructure:"jaeger" json:"jaeger" yaml:"jaeger"`
}
