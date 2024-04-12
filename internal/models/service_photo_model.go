package models

import "path/filepath"

type ServicePhotoModel struct {
	Model
	UpdateableModel
	ServiceId int    `json:"serviceId" db:"serviceId"`
	VendorId  int    `json:"vendorId" db:"vendorId"`
	Url       string `json:"url" db:"url"`
}

func NewServicePhotoModel(vendorId, serviceId int, filename string) *ServicePhotoModel {
	fileLocation := filepath.Join("/resource/service", filename)

	return &ServicePhotoModel{
		ServiceId: serviceId,
		VendorId:  vendorId,
		Url:       fileLocation,
	}
}
