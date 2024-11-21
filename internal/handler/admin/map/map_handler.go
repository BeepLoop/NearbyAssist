package map_handler

import (
	map_service "nearbyassist/internal/service/map"
	tag_service "nearbyassist/internal/service/tag"
)

type mapHandler struct {
	tagService *tag_service.Service
	mapService *map_service.Service
}

func NewHandler(mapService *map_service.Service, tagService *tag_service.Service) *mapHandler {
	return &mapHandler{mapService: mapService, tagService: tagService}
}
