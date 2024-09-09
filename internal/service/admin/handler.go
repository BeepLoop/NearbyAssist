package service

import (
	store "nearbyassist/internal/store/admin"
)

type AdminService struct {
	store store.AdminStore
}

func NewAdminService(store store.AdminStore) *AdminService {
	return &AdminService{
		store: store,
	}
}
