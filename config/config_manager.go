package config

import (
	"github.com/xdevices/utilities/config"
	"github.com/xdevices/utilities/rabbit"
)

type resolverManager struct {
	*config.Manager
	rabbit.RabbitMQManager
}

var instance *resolverManager

func Config() *resolverManager {
	if instance == nil {
		instance = new(resolverManager)
		instance.Manager = new(config.Manager)
	}
	return instance
}

func (c *resolverManager) Init() {
	c.Manager.Init()
	c.RabbitMQManager.InitConnection(c.RabbitURL())
}
