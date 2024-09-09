package admin

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/id_generator"
	store "nearbyassist/internal/store/admin"
)

type Handler struct {
	store store.IAdminStore
	jwt   authenticator.Authenticator
	idGen id_generator.IdGenerator
}

func NewHandler(store store.IAdminStore, jwt authenticator.Authenticator, idGen id_generator.IdGenerator) *Handler {
	return &Handler{
		store: store,
		jwt:   jwt,
		idGen: idGen,
	}
}
