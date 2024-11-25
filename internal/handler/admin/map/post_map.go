package map_handler

import (
	"context"
	"nearbyassist/views/pages/map"

	"github.com/labstack/echo/v4"
)

func (h *mapHandler) PostMap(c echo.Context) error {
	query := c.FormValue("query")

	markers := make([]pages.Marker, 0)
	tags := make([]string, 0)

	pageData := pages.MapPageData{
		Markers:       markers,
		Tags:          tags,
		PreviousQuery: query,
	}

	result, err := h.tagService.GetTags()
	if err != nil {
		page := pages.Map(pageData)
		return page.Render(context.Background(), c.Response().Writer)
	}

	for _, tag := range result {
		pageData.Tags = append(pageData.Tags, tag.Title)
	}

	services, err := h.mapService.GetServices(query)
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

	page := pages.Map(pageData)
	return page.Render(context.Background(), c.Response().Writer)
}
