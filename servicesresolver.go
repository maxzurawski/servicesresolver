package servicesresolver

import (
	"github.com/labstack/echo"
	"github.com/xdevices/servicesresolver/config"
	"github.com/xdevices/servicesresolver/handlers"
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
}
