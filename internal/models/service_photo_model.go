package models

import "path/filepath"

type ServicePhotoModel struct {
	Model
	UpdateableModel
	ServiceId string `json:"serviceId" db:"serviceId"`
	VendorId  string `json:"vendorId" db:"vendorId"`
	Url       string `json:"url" db:"url"`
}

func NewServicePhotoModel(servicePhotoId, vendorId, serviceId string, filename string) *ServicePhotoModel {
	fileLocation := filepath.Join("/resource/service", filename)

	return &ServicePhotoModel{
		Model:     Model{Id: servicePhotoId},
		ServiceId: serviceId,
		VendorId:  vendorId,
		Url:       fileLocation,
	}
}
