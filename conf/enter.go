package conf

type Config struct {
	System System `yaml:"system"`
	DB     []DB   `yaml:"db"`
	Jwt    Jwt    `yaml:"jwt"`
}
