package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/xdevices/servicesresolver/config"
	"github.com/xdevices/servicesresolver/dto"
	"github.com/xdevices/utilities/stringutils"
)

func HandleAppRequest(c echo.Context) error {
	app := c.Param("app")
	if stringutils.IsZero(app) {
		return c.JSON(http.StatusBadRequest, &dto.AppMetadata{})
	}

	eurekaAppUrl := fmt.Sprintf("%s/%s/%s", config.Config().EurekaService(), "eureka/apps", app)

	request, err := http.NewRequest("GET", eurekaAppUrl, nil)
	request.Header.Set("Accept", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		config.Config().Logger.PublishError(uuid.New().String(),
			"",
			config.Config().ServiceName(),
			fmt.Sprintf("problem requesting eureka regarding app: [%s]", app),
			err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var application dto.EurekaResponse
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		config.Config().Logger.PublishError(
			uuid.New().String(),
			"",
			config.Config().ServiceName(),
			fmt.Sprintf("could not read body of response."),
			err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if response.StatusCode == http.StatusNoContent || response.StatusCode == http.StatusNotFound {
		config.Config().Logger.PublishInfo(
			uuid.New().String(), "",
			config.Config().ServiceName(),
			fmt.Sprintf("no content found in eureka for app: [%s]", app))
		return c.JSON(http.StatusNoContent, &dto.AppMetadata{})
	}

	err = json.Unmarshal(bytes, &application)
	if err != nil {
		config.Config().Logger.PublishError(uuid.New().String(), "",
			config.Config().ServiceName(), fmt.Sprintf("encountered error during unmarshaling message: [%s]", string(bytes)),
			err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	upInstance := &dto.Instance{}
	for _, item := range application.Application.Instance {
		if item.Status == "UP" {
			upInstance = &item
			break
		}
	}

	if stringutils.IsZero(upInstance.IpAddr) {
		return c.JSON(http.StatusInternalServerError, "error. up instance not found")
	}
	port := strconv.Itoa(upInstance.Port.Port)

	metadata := dto.AppMetadata{
		Name: &application.Application.Name,
		Ip:   &upInstance.IpAddr,
		Port: &port,
	}

	return c.JSON(http.StatusOK, metadata)
}
