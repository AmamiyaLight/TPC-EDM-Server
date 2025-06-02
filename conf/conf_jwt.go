package conf

type Jwt struct {
	Expire int    `yaml:"expire"` //单位:Hour
	Secret string `yaml:"secret"`
	Issuer string `yaml:"issuer"`
}
