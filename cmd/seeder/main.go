package main

import (
	"nearbyassist/internal/config"
	"nearbyassist/internal/db/mysql"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/types"
)

func main() {
	conf := config.LoadConfig()
	db := mysql.NewMysqlDatabase(conf)

	// Seed categories
	_, err := db.Conn.NamedExec("INSERT INTO Tag (title) values (:title)", []types.Category{
		{Title: "computer repair"},
		{Title: "plumbing"},
		{Title: "electric"},
	})
	if err != nil {
		panic("Error inserting category: " + err.Error())
	}

	// Seed Admin
	_, err = db.Conn.NamedExec("INSERT INTO Admin (username, password, role) VALUES (:username, :password, :role)", []struct {
		Username string `db:"username"`
		Password string `db:"password"`
		Role     string `db:"role"`
	}{
		{Username: "admin", Password: "admin", Role: "admin"},
		{Username: "dui", Password: "pass", Role: "staff"},
	})
	if err != nil {
		panic("Error inserting admin: " + err.Error())
	}

	// Seed users
	_, err = db.Conn.NamedExec("INSERT INTO User (name, email, imageUrl) values (:name, :email, :imageUrl)", []types.User{
		{
			Name:     "John Loyd Mulit",
			Email:    "jlmulit68@gmail.com",
			ImageUrl: "https://i.pravatar.cc/100",
		},
	})
	if err != nil {
		panic("Error inserting user: " + err.Error())
	}

	// Seed vendors
	_, err = db.Conn.NamedExec("INSERT INTO Vendor (vendorId, job) values ((SELECT id FROM User WHERE name = :name), :job)", []struct {
		Name string `db:"name"`
		Job  string `db:"job"`
	}{
		{Name: "John Loyd Mulit", Job: "Plumber"},
	})
	if err != nil {
		panic("Error inserting vendors: " + err.Error())
	}

	// Seed sevices
	services := []request.NewService{
		{
			VendorId:    1,
			Description: "We offer computer repair and maintenance services.",
			Rate:        "100",
			Tags: []string{
				"computer repair",
				"electric",
			},
			GeoSpatialModel: models.GeoSpatialModel{
				Latitude:  7.422302,
				Longitude: 125.824747,
			},
		},
		{
			VendorId:    1,
			Description: "We offer plumbing services.",
			Rate:        "100",
			Tags: []string{
				"plumbing",
			},
			GeoSpatialModel: models.GeoSpatialModel{
				Latitude:  7.419594,
				Longitude: 125.824616,
			},
		},
	}

	registerService := `
	        INSERT INTO
	            Service
	                (vendorId, description, rate, latitude, longitude)
	        VALUES 
                (
                    :vendorId,
                    :description,
                    :rate,
                    :latitude,
                    :longitude
                )
	    `

	for _, service := range services {
		tx, err := db.Conn.Beginx()
		if err != nil {
			panic(err)
		}

		models.ConstructLocationFromLatLong(&service.GeoSpatialModel)

		res, err := tx.NamedExec(registerService, service)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				panic("failed to rollback on insert service: " + err.Error())
			}

			panic(err)
		}

		serviceId, err := res.LastInsertId()
		if err != nil {
			panic("unable to get service ID")
		}

		registerTag := `
        INSERT INTO 
            Service_Tag (serviceId, tagId)
        VALUES
            (
                ?,
                (SELECT id FROM Tag WHERE title = ?)
            )
    `

		var tagErr error
		for _, tag := range service.Tags {
			if _, err := tx.Exec(registerTag, serviceId, tag); err != nil {
				tagErr = err
				break
			}
		}

		if tagErr != nil {
			if err := tx.Rollback(); err != nil {
				panic("failed to rollback on insert tag " + err.Error())
			}

			panic(err.Error())
		}

		if err := tx.Commit(); err != nil {
			if err := tx.Rollback(); err != nil {
				panic("failed to rollback on commit " + err.Error())
			}
		}
	}

	// Seed reviews
	_, err = db.Conn.NamedExec("INSERT INTO Review (serviceId, rating) values (:serviceId, :rating)", []types.Review{
		{ServiceId: 1, Rating: 5},
		{ServiceId: 1, Rating: 3},
		{ServiceId: 1, Rating: 3},
		{ServiceId: 1, Rating: 4},
	})
	if err != nil {
		panic("Error inserting reviews: " + err.Error())
	}

	// Seed service photos
	_, err = db.Conn.NamedExec("INSERT INTO ServicePhoto (vendorId, serviceId, url) values (:vendorId, :serviceId, :url)", []types.ServicePhoto{
		{ServiceId: 1, VendorId: 1, Url: "https://i.pravatar.cc/100"},
		{ServiceId: 1, VendorId: 1, Url: "https://i.pravatar.cc/100"},
	})
	if err != nil {
		panic("Error inserting photos: " + err.Error())
	}
}
