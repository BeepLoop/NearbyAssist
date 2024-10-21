package map_handler

import map_service "nearbyassist/internal/service/map"

type mapHandler struct {
	mapService *map_service.Service
}

func NewHandler(mapService *map_service.Service) *mapHandler {
	return &mapHandler{mapService: mapService}
}
