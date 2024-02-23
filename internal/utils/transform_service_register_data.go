package utils

import (
	"fmt"
	"nearbyassist/internal/types"
)

func TransformServiceData(service types.ServiceRegister) (*types.ServiceRegister, error) {

	point := fmt.Sprintf(
		"POINT(%f %f)",
		service.Latitude,
		service.Longitude,
	)

	transformedData := types.ServiceRegister{
		VendorId:    service.VendorId,
		Title:       service.Title,
		Description: service.Description,
		Rate:        service.Rate,
		CategoryId:  service.CategoryId,
		Point:     point,
	}

	return &transformedData, nil
}
