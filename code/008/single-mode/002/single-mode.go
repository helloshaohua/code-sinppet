package main

import "sync"

type config struct {
}

var cfg *config
var once sync.Once

func GetConfigInstance() *config {
	once.Do(func() {
		cfg = new(config)
	})
	return cfg
}
