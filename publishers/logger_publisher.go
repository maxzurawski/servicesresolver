package publishers

import (
	"github.com/labstack/gommon/log"
	"github.com/xdevices/servicesresolver/config"
	"github.com/xdevices/utilities/rabbit/crosscutting"
	"github.com/xdevices/utilities/rabbit/publishing"
)

type logger struct {
	*publishing.Publisher
}

var publisher *publishing.Publisher
var loggerInstance *logger

func InitLogger() {
	if publisher == nil && config.Config().ConnectToRabbit() {
		publisher = config.Config().InitPublisher()
		publisher.DeclareTopicExchange(crosscutting.TopicLogs.String())
	}
}

func Logger() *logger {
	if loggerInstance == nil {
		loggerInstance = new(logger)
		loggerInstance.Publisher = publisher
	}
	return loggerInstance
}

func (l *logger) Info(processId, sensorUuid, msg string) {

	if !config.Config().ConnectToRabbit() {
		log.Info("connection to rabbit disabled")
		return
	}

	l.Reset()
	l.PublishInfo(processId,
		sensorUuid,
		config.Config().ServiceName(),
		msg,
	)
}

func (l *logger) Error(processId, sensorUuid, msg, details string) {
	if !config.Config().ConnectToRabbit() {
		log.Info("connection to rabbit disabled")
		return
	}

	l.Reset()
	l.PublishError(processId,
		sensorUuid,
		config.Config().ServiceName(),
		msg,
		details,
	)
}
