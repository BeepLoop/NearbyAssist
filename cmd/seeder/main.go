package main

import (
	"log"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/db/query/service"
	"nearbyassist/internal/db/query/user"
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

	// Seed customers
	customers := []types.User{
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
	}
	for _, customer := range customers {
		err := user_query.RegisterUser(customer)
		if err != nil {
			panic("Error inserting customer: " + err.Error())
		}
	}

	// Seed sevices
	services := []types.ServiceRegister{
		{
			VendorId:    1,
			Title:       "Sugar & Leerd Bakery",
			Description: "We bake the best cakes in town",
			Rate:        100,
			Latitude:    7.419594,
			Longitude:   125.824616,
			CategoryId:  1,
		},
		{
			VendorId:    1,
			Title:       "11:11 Cafe",
			Description: "We serve the best coffee in town",
			Rate:        100,
			Latitude:    7.422325,
			Longitude:   125.824777,
			CategoryId:  1,
		},
		{
			VendorId:    1,
			Title:       "Minute Burger",
			Description: "We serve the best burgers in town",
			Rate:        100,
			Latitude:    7.4234,
			Longitude:   125.828901,
			CategoryId:  1,
		},
		{
			VendorId:    1,
			Title:       "Hugo Bistro",
			Description: "We serve the best pizza in town",
			Rate:        100,
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

		err = service_query.RegisterService(*data)
		if err != nil {
			panic("Error inserting service: " + err.Error())
		}
	}

	// Seed vendors
	_, err = db.Connection.NamedExec("INSERT INTO Vendor (vendorId, rating) values (:vendorId, :rating)", []struct {
		VendorId int     `db:"vendorId"`
		Rating   float32 `db:"rating"`
	}{
		{VendorId: 1, Rating: 4.5},
		{VendorId: 2, Rating: 3.5},
	})
	if err != nil {
		panic("Error inserting vendors: " + err.Error())
	}

	// Seed reviews
	_, err = db.Connection.NamedExec("INSERT INTO Review (vendorId, rating) values (:vendorId, :rating)", []types.Review{
		{VendorId: 1, Rating: 5},
		{VendorId: 1, Rating: 3},
		{VendorId: 1, Rating: 3},
	})
}
