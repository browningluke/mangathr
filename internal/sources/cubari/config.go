package cubari

var config Config

type Config struct {
}

func SetConfig(cfg Config) {
	config = cfg
}

func (c *Config) Default() {
}
