package main

import "pomment-go/config"

// InitConfig Initialize Pomment from config file
func InitConfig(basePath string) (err error) {
	return config.InitConfig(basePath)
}
