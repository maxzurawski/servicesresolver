package main

import (
	"github.com/labstack/echo"
	"github.com/maxzurawski/servicesresolver/config"
	"github.com/maxzurawski/servicesresolver/handlers"
	"github.com/maxzurawski/servicesresolver/publishers"
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
