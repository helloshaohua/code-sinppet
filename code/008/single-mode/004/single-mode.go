package main

type config struct {
}

var cfg *config = new(config)

// GetConfigInstance 提供获取实例的方法
func GetConfigInstance() *config {
	return cfg
}
