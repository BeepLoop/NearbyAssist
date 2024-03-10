package utils

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetUploadParams(c echo.Context) (vendor int, service int, err error) {
	vendorId := c.QueryParam("vendorId")
	serviceId := c.QueryParam("serviceId")
	if vendorId == "" || serviceId == "" {
		return 0, 0, errors.New("missing fields")
	}

	vendor, err = strconv.Atoi(vendorId)
	if err != nil {
		return 0, 0, err
	}

	service, err = strconv.Atoi(serviceId)
	if err != nil {
		return 0, 0, err
	}

	return vendor, service, nil
}
