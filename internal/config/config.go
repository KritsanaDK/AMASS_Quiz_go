
package config

type Config struct {
	Port  string
	Debug bool
}

func Load() *Config {
	return &Config{
		Port:  "8080",
		Debug: true,
	}
}
