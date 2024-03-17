package main

import (
	"log"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
)

func init() {
	if err := config.Init(); err != nil {
		log.Fatal("Error initializing config: ", err)
	}

	if err := db.Init(); err != nil {
		log.Fatal("Error initializing database: ", err)
	}
}

func main() {
	// Seed categories
	_, err := db.Connection.NamedExec("INSERT INTO Category (title) values (:title)", []types.Category{
		{Title: "food"},
		{Title: "service"},
	})
	if err != nil {
		panic("Error inserting category: " + err.Error())
	}

	// Seed users
	_, err = db.Connection.NamedExec("INSERT INTO User (name, email, imageUrl) values (:name, :email, :imageUrl)", []types.User{
		{
			Name:     "John Loyd Mulit",
			Email:    "jlmulit@email.com",
			ImageUrl: "https://dummyimage.com/400x400/000/fff",
		},
		{
			Name:     "Cherry Lyn Burlat",
			Email:    "clburlat@email.com",
			ImageUrl: "https://dummyimage.com/400x400/000/fff",
		},
		{
			Name:     "Adrian Juntilla",
			Email:    "ajuntilla@email.com",
			ImageUrl: "https://dummyimage.com/400x400/000/fff",
		},
	})
	if err != nil {
		panic("Error inserting user: " + err.Error())
	}

	// Seed sevices
	services := []types.ServiceRegister{
		{
			VendorId:    1,
			Title:       "Sugar & Leerd Bakery",
			Description: "We bake the best cakes in town",
			Rate:        100.00,
			Latitude:    7.419594,
			Longitude:   125.824616,
			CategoryId:  1,
		},
		{
			VendorId:    2,
			Title:       "11:11 Cafe",
			Description: "We serve the best coffee in town",
			Rate:        100.00,
			Latitude:    7.422325,
			Longitude:   125.824777,
			CategoryId:  1,
		},
		{
			VendorId:    2,
			Title:       "Minute Burger",
			Description: "We serve the best burgers in town",
			Rate:        100.00,
			Latitude:    7.4234,
			Longitude:   125.828901,
			CategoryId:  1,
		},
		{
			VendorId:    1,
			Title:       "Hugo Bistro",
			Description: "We serve the best pizza in town",
			Rate:        100.00,
			Latitude:    7.424179,
			Longitude:   125.829182,
			CategoryId:  1,
		},
	}
	for _, service := range services {
		data, err := utils.TransformServiceData(service)
		if err != nil {
			panic("Error transforming service data: " + err.Error())
		}

		_, err = db.Connection.NamedExec("INSERT INTO Service (vendor, title, description, rate, location, category) values (:vendorId, :title, :description, :rate, ST_GeomFromText(:point, 4326), :categoryId)", data)
		if err != nil {
			panic("Error inserting service: " + err.Error())
		}
	}

	// Seed vendors
	_, err = db.Connection.NamedExec("INSERT INTO Vendor (vendorId, rating, role) values (:vendorId, :rating, :role)", []types.VendorData{
		{VendorId: 1, Role: "plumber"},
		{VendorId: 2, Role: "electrician"},
	})
	if err != nil {
		panic("Error inserting vendors: " + err.Error())
	}

	// Seed reviews
	_, err = db.Connection.NamedExec("INSERT INTO Review (serviceId, rating) values (:serviceId, :rating)", []types.Review{
		{ServiceId: 1, Rating: 5},
		{ServiceId: 1, Rating: 3},
		{ServiceId: 1, Rating: 3},
	})
	if err != nil {
		panic("Error inserting reviews: " + err.Error())
	}

	// Seed service photos
	_, err = db.Connection.NamedExec("INSERT INTO Photo (vendorId, serviceId, url) values (:vendorId, :serviceId, :url)", []types.Photo{
		{ServiceId: 1, VendorId: 1, Url: "https://dummyimage.com/400x400/000/fff"},
		{ServiceId: 1, VendorId: 1, Url: "https://dummyimage.com/400x400/000/fff"},
		{ServiceId: 2, VendorId: 2, Url: "https://dummyimage.com/400x400/000/fff"},
		{ServiceId: 2, VendorId: 2, Url: "https://dummyimage.com/400x400/000/fff"},
	})
	if err != nil {
		panic("Error inserting photos: " + err.Error())
	}
}
