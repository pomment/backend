package pomment

import (
	"github.com/pomment/pomment/config"
	"github.com/pomment/pomment/dao"
)

// InitConfig Initialize Pomment from config file
func InitConfig(basePath string) (err error) {
	return config.InitConfig(basePath)
}

// ConnectToRedisServer Connect to a Redis server for pomment service
func ConnectToRedisServer(enabled bool, addr string, password string, db int) {
	dao.ConnectToRedisServer(enabled, addr, password, db)
}
