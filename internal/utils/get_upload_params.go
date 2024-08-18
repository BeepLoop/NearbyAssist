package utils

import (
	"errors"

	"github.com/labstack/echo/v4"
)

func GetUploadParams(c echo.Context, params ...string) (map[string]string, error) {
	paramIds := make(map[string]string)

	var err error
	for _, param := range params {
		paramId := c.QueryParam(param)

		if paramId == "" {
			break
		}

		if _, ok := paramIds[param]; ok {
			err = errors.New("duplicate parameter")
			break
		}

		paramIds[param] = paramId
	}

	if err != nil {
		return nil, err
	}

	return paramIds, nil
}
