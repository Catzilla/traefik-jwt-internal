package traefik_jwt_internal

type HeaderConfig struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}

type Config struct {
	Alg    string        `yaml:"alg"`
	Secret string        `yaml:"secret"`
	TTL    int64         `yaml:"ttl"`
	Header *HeaderConfig `yaml:"header"`
	Claims string        `yaml:"claims"`
}

func CreateConfig() *Config {
	return &Config{
		Alg:    "HS256",
		Secret: "",
		TTL:    120,
		Header: &HeaderConfig{
			Name:   "Authorization",
			Prefix: "Bearer",
		},
		Claims: "{}",
	}
}
