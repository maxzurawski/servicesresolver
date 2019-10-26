package main

import (
	"github.com/labstack/echo"
	"github.com/xdevices/servicesresolver/config"
	"github.com/xdevices/servicesresolver/handlers"
	"github.com/xdevices/servicesresolver/publishers"
)

func main() {
	e := echo.New()
	e.GET("/:app", handlers.HandleAppRequest)
	e.Logger.Fatal(e.Start(config.Config().Address()))
}

func init() {
	config.Config().Init()
	manager := config.EurekaManagerInit()
	manager.SendRegistrationOrFail()
	manager.ScheduleHeartBeat(config.Config().ServiceName(), 10)
	publishers.InitLogger()
}
