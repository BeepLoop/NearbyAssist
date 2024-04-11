package models

// import (
//
//	"context"
//	"errors"
//	"fmt"
//	"nearbyassist/internal/db"
//	"time"
//
// )
type ServiceModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	VendorId    int    `json:"vendorId" db:"vendorId" validate:"required"`
	Title       string `json:"title" db:"title" validate:"required"`
	Description string `json:"description" db:"description" validate:"required"`
	Rate        string `json:"rate" db:"rate" validate:"required"`
	CategoryId  int    `json:"categoryId" db:"categoryId" validate:"required"`
}

func NewServiceModel() *ServiceModel {
	return &ServiceModel{}
}

//
// func NewServiceModelWithLocation(latitude, longitude float64) *ServiceModel {
// 	return &ServiceModel{
// 		GeoSpatialModel: GeoSpatialModel{
// 			Latitude:  latitude,
// 			Longitude: longitude,
// 		},
// 	}
// }
//
// func (s *ServiceModel) Create() (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
// 	        INSERT INTO
// 	            Service
// 	                (vendorId, title, description, rate, location, categoryId)
// 	        VALUES
//                 (
//                     :vendorId,
//                     :title,
//                     :description,
//                     :rate,
//                     ST_GeomFromText(:location, 4326),
//                     :categoryId
//                 )
// 	    `
//
// 	ConstructLocationFromLatLong(&s.GeoSpatialModel)
//
// 	res, err := s.Db.Conn.NamedExecContext(ctx, query, s)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	insertId, err := res.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return 0, context.DeadlineExceeded
// 	}
//
// 	return int(insertId), nil
// }
//
// func (s *ServiceModel) Update(id int) error {
// 	return nil
// }
//
// func (s *ServiceModel) Delete(id int) error {
// 	return nil
// }
//
// func (s *ServiceModel) FindAll() ([]*ServiceModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         SELECT
//             id, vendorId, title, description, rate, ST_AsText(location) as location, categoryId
//         FROM
//             Service
//         LIMIT
//             10
//     `
//
// 	services := make([]*ServiceModel, 0)
// 	err := s.Db.Conn.SelectContext(ctx, &services, query)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, service := range services {
// 		err := ExtractLatLongFromLocation(&service.GeoSpatialModel)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return services, nil
// }
//
// func (s *ServiceModel) FindById(id int) (*ServiceModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         SELECT
//             id, vendorId, title, description, rate, ST_AsText(location) as location, categoryId
//         FROM
//             Service
//         WHERE
//             id = ?
//     `
//
// 	service := new(ServiceModel)
// 	err := s.Db.Conn.GetContext(ctx, service, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	err = ExtractLatLongFromLocation(&service.GeoSpatialModel)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return service, nil
// }
//
// func (s *ServiceModel) FindByVendorId(vendorId int) ([]*ServiceModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         SELECT
//             id, vendorId, title, description, rate, ST_AsText(location) as location, categoryId, createdAt, updatedAt
//         FROM
//             Service
//         WHERE
//             vendorId = ?
//     `
//
// 	services := make([]*ServiceModel, 0)
// 	err := s.Db.Conn.SelectContext(ctx, &services, query, vendorId)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, service := range services {
// 		err := ExtractLatLongFromLocation(&service.GeoSpatialModel)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return services, nil
// }
//
// func (s *ServiceModel) GeoSpatialSearch(searchTerm string, radius float64) ([]*ServiceModel, error) {
// 	if s.Latitude == 0 || s.Longitude == 0 {
// 		return nil, errors.New("Latitude and Longitude must be set")
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := fmt.Sprintf(`
//         SELECT
//             id, vendorId, title, description, format(rate, 2) as rate, ST_AsText(location) as location, categoryId
//         FROM
//             Service
//         WHERE
//             title LIKE '%%%s%%'
//         AND
//             ST_Distance_Sphere(
//                 location,
//                 ST_GeomFromText('POINT(%f %f)', 4326)
//             ) < ?;
//     `, searchTerm, s.Latitude, s.Longitude)
//
// 	services := make([]*ServiceModel, 0)
// 	err := s.Db.Conn.SelectContext(ctx, &services, query, radius)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, service := range services {
// 		err := ExtractLatLongFromLocation(&service.GeoSpatialModel)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return services, nil
// }
