package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/maxzurawski/servicesresolver/config"
	"github.com/maxzurawski/servicesresolver/dto"
	"github.com/maxzurawski/servicesresolver/publishers"
	"github.com/maxzurawski/utilities/stringutils"
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
		publishers.Logger().Error(uuid.New().String(),
			"",
			fmt.Sprintf("problem requesting eureka regarding app: [%s]", app),
			err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		publishers.Logger().Error(
			uuid.New().String(),
			"",
			fmt.Sprintf("could not read body of response."),
			err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if response.StatusCode == http.StatusNoContent || response.StatusCode == http.StatusNotFound {
		publishers.Logger().Info(
			uuid.New().String(), "",
			fmt.Sprintf("no content found in eureka for app: [%s]", app))
		return c.JSON(http.StatusNoContent, &dto.AppMetadata{})
	}

	var application dto.EurekaResponse
	err = json.Unmarshal(bytes, &application)
	if err != nil {
		publishers.Logger().Error(uuid.New().String(),
			"",
			fmt.Sprintf("encountered error during unmarshaling message: [%s]", string(bytes)),
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
