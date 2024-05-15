package pomment

import (
	"github.com/pomment/pomment/config"
	"github.com/pomment/pomment/dao"
	"path"
)

// InitConfig Initialize Pomment from config file
func InitConfig(basePath string) (err error) {
	dao.InitDataBasePath(path.Join(basePath))
	return config.InitConfig(basePath)
}

func InitExampleConfig(basePath string) (err error) {
	return config.InitExampleConfig(basePath)
}

// ConnectToRedisServer Connect to a Redis server for pomment service
func ConnectToRedisServer(enabled bool, addr string, password string, db int) {
	dao.ConnectToRedisServer(enabled, addr, password, db)
}
