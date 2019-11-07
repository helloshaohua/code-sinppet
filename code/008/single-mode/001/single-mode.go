package main

type config struct {
}

var cfg *config

func GetConfigInstance() *config {
	if cfg == nil {
		cfg = new(config)
		return cfg
	}
	return cfg
}
