package utils

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetUploadParams(c echo.Context, params ...string) (map[string]int, error) {
	paramIds := make(map[string]int)

	var err error
	for _, param := range params {
		paramId := c.QueryParam(param)

		var id int
		if id, err = strconv.Atoi(paramId); err != nil {
			break
		}

		if _, ok := paramIds[param]; ok {
			err = errors.New("duplicate parameter")
			break
		}

		paramIds[param] = id
	}

	if err != nil {
		return nil, err
	}

	return paramIds, nil
}
