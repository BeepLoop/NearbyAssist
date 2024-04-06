package models

import (
	"context"
	"nearbyassist/internal/db"
	"time"
)

type VendorModel struct {
	Model
	UpdateableModel
	VendorId   int    `json:"vendorId" db:"vendorId"`
	Rating     string `json:"rating" db:"rating"`
	Job        string `json:"job" db:"job"`
	Restricted int    `json:"restricted" db:"restricted"`
}

func NewVendorModel() *VendorModel {
	return &VendorModel{}
}

func (v *VendorModel) Create() (int, error) {
	return 0, nil
}

func (v *VendorModel) Update(id int) error {
	return nil
}

func (v *VendorModel) Delete(id int) error {
	return nil
}

func (v *VendorModel) FindById(id int) (*VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            id, vendorId, rating, job, restricted
        FROM
            Vendor
        WHERE
            id = ?
    `

	vendor := new(VendorModel)
	err := db.Connection.GetContext(ctx, vendor, query, id)
	if err != nil {
		return nil, err
	}

	return vendor, nil
}

func (v *VendorModel) RestrictAccount(vendorId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        UPDATE
            Vendor 
        SET 
            restricted = 1
        WHERE 
            vendorId = ?
    `

	_, err := db.Connection.ExecContext(ctx, query, vendorId)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (v *VendorModel) UnrestrictAccount(vendorId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        UPDATE
            Vendor 
        SET 
            restricted = 0
        WHERE 
            vendorId = ?
    `

	_, err := db.Connection.ExecContext(ctx, query, vendorId)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (v *VendorModel) Count(filter string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            COUNT(*)
        FROM 
            Vendor
    `

	switch filter {
	case "restricted":
		query += " WHERE restricted = 1"
	case "unrestricted":
		query += " WHERE restricted = 0"
	}

	count := 0
	err := db.Connection.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}
