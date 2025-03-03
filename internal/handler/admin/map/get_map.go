package map_handler

import (
	"context"
	"nearbyassist/views/pages/map"

	"github.com/labstack/echo/v4"
)

func (h *mapHandler) GetMap(c echo.Context) error {
	markers := make([]pages.Marker, 0)
	tags := make([]string, 0)

	pageData := pages.MapPageData{
		Markers: markers,
		Tags:    tags,
	}

	params := c.QueryParams()
	if params.Has("query") {
		services, err := h.mapService.GetServices(params.Get("query"))
		if err != nil {
			page := pages.Map(pageData)
			return page.Render(context.Background(), c.Response().Writer)
		}

		for _, service := range services {
			pageData.Markers = append(pageData.Markers, pages.Marker{
				Lat: service.Latitude,
				Lon: service.Longitude,
			})
		}
	}

	result, err := h.tagService.GetTags()
	if err != nil {
		page := pages.Map(pageData)
		return page.Render(context.Background(), c.Response().Writer)
	}

	for _, tag := range result {
		pageData.Tags = append(pageData.Tags, tag.Title)
	}

	page := pages.Map(pageData)
	return page.Render(context.Background(), c.Response().Writer)
}
