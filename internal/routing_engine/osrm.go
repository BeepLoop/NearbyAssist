package routing_engine

import (
	"context"
	"encoding/json"
	"io"
	"nearbyassist/internal/config"
	"nearbyassist/internal/models"
	"net/http"
	"time"
)

type osrmRoute struct {
	Geometry string  `json:"geometry"`
	Duration float32 `json:"duration"`
	Distance float32 `json:"distance"`
}

type osrmResponse struct {
	Code   string      `json:"code"`
	Routes []osrmRoute `json:"routes"`
}

type OSRM struct {
	engineUrl      string
	requestTimeout time.Duration
}

func NewOSRM(conf *config.Config) *OSRM {
	return &OSRM{
		engineUrl:      conf.ROUTE_ENGINE_URL,
		requestTimeout: 5 * time.Second,
	}
}

func (e *OSRM) constructUrl(origin, destination models.GeoSpatialModel) string {
	return e.engineUrl + "/route/v1/driving/" + origin.StringReverseOrder() + ";" + destination.StringReverseOrder()
}

func (e *OSRM) GetPolyline(origin, destination *models.GeoSpatialModel) (PolylineCode, error) {
	data, err := e.findRoute(origin, destination)
	if err != nil {
		return "", err
	}

	return PolylineCode(data.Routes[0].Geometry), nil
}

func (e *OSRM) GetDistance(origin, destination *models.GeoSpatialModel) (float32, error) {
	data, err := e.findRoute(origin, destination)
	if err != nil {
		return 0, err
	}

	return data.Routes[0].Distance, nil
}

func (e *OSRM) findRoute(origin, destination *models.GeoSpatialModel) (*osrmResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.requestTimeout)
	defer cancel()

	url := e.constructUrl(*origin, *destination)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := new(osrmResponse)
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}
