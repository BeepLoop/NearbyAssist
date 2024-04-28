package main

import (
	"nearbyassist/internal/config"
	"nearbyassist/internal/db/mysql"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
)

func main() {
	conf := config.LoadConfig()
	db := mysql.NewMysqlDatabase(conf)

	// Seed categories
	_, err := db.Conn.NamedExec("INSERT INTO Category (title) values (:title)", []types.Category{
		{Title: "food"},
		{Title: "service"},
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
	services := []types.ServiceRegister{
		{
			VendorId:    1,
			Title:       "Computer Repair & Maintenance",
			Description: "We offer computer repair and maintenance services.",
			Rate:        100.00,
			Latitude:    7.422302,
			Longitude:   125.824747,
			CategoryId:  2,
		},
		{
			VendorId:    1,
			Title:       "Plumbing Services",
			Description: "We offer plumbing services.",
			Rate:        100.00,
			Latitude:    7.419594,
			Longitude:   125.824616,
			CategoryId:  2,
		},
	}
	for _, service := range services {
		data, err := utils.TransformServiceData(service)
		if err != nil {
			panic("Error transforming service data: " + err.Error())
		}

		_, err = db.Conn.NamedExec("INSERT INTO Service (vendorId, title, description, rate, location, categoryId) values (:vendorId, :title, :description, :rate, ST_GeomFromText(:point, 4326), :categoryId)", data)
		if err != nil {
			panic("Error inserting service: " + err.Error())
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
