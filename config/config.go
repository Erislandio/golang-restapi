package config

import "github.com/tkanos/gonfig"

// Config struc .
type Config struct {
	DB_NAME string
	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASS string
}

// GetConfig .
func GetConfig() Config {
	conf := Config{}

	gonfig.GetConf("config/config.json", &conf)

	return conf
}
