package config

import (
	"github.com/xdevices/utilities/config"
	"github.com/xdevices/utilities/rabbit"
	"github.com/xdevices/utilities/rabbit/crosscutting"
	"github.com/xdevices/utilities/rabbit/publishing"
)

type resolverManager struct {
	*config.Manager
	rabbit.RabbitMQManager
	Logger *publishing.Publisher
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
	c.Logger = c.RabbitMQManager.InitPublisher()
	c.Logger.Reset()
	c.Logger.DeclareTopicExchange(crosscutting.TopicLogs.String())
}
