package models

import (
	"context"
	"nearbyassist/internal/id_generator"
	"time"

	"github.com/jmoiron/sqlx"
)

type VendorStatusFilter string

const (
	VENDOR_STATUS_RESTRICTED   VendorStatusFilter = "restricted"
	VENDOR_STATUS_UNRESTRICTED VendorStatusFilter = "unrestricted"
	VENDOR_STATUS_ALL          VendorStatusFilter = "all"
)

type VendorModel struct {
	Model
	UpdateableModel
	VendorId   string `json:"vendorId" db:"vendorId"`
	Rating     string `json:"rating" db:"rating"`
	Job        string `json:"job" db:"job"`
	Restricted int    `json:"restricted" db:"restricted"`

	// Additional fields for joins
	Vendor   string `json:"vendor" db:"vendor"`
	ImageUrl string `json:"imageUrl" db:"imageUrl"`
}

func NewVendorModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *VendorModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &VendorModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func (v *VendorModel) DecryptVendor(decryptFunc func(string) (string, error)) (*VendorModel, error) {
	if decrypted, err := decryptFunc(v.Vendor); err != nil {
		return nil, err
	} else {
		v.Vendor = decrypted
	}
	return v, nil
}

func (v *VendorModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO Vendor (id, vendorId, rating, job) VALUES (:id, :vendorId, :rating, :job)"
	if _, err := v.Conn.NamedExecContext(ctx, query, v); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return "", nil
}

func (v *VendorModel) Count(filter VendorStatusFilter) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM Vendor"

	switch filter {
	case VENDOR_STATUS_RESTRICTED:
		query += " WHERE restricted = 1"
	case VENDOR_STATUS_UNRESTRICTED:
		query += " WHERE restricted = 0"
	case VENDOR_STATUS_ALL:
		//
	}

	count := 0
	err := v.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil

}

func (v *VendorModel) FindById(id string) (*VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, vendorId, rating, job, restricted FROM Vendor WHERE vendorId = ?"

	err := v.Conn.GetContext(ctx, v, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return v, nil
}

func (v *VendorModel) FindByServiceId(id string) (*VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            v.vendorId,
            v.rating,
            v.job,
            u.name as vendor,
            u.imageUrl as imageUrl
        FROM
            Vendor v
            JOIN Service s ON s.vendorId = v.vendorId
            JOIN User u ON u.id = v.vendorId
        WHERE 
            s.id = ?
    `

	err := v.Conn.GetContext(ctx, v, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return v, nil
}

func (v *VendorModel) Restrict(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE Vendor SET restricted = 1 WHERE vendorId = ?"
	if _, err := v.Conn.ExecContext(ctx, query, id); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (v *VendorModel) Unrestrict(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE Vendor SET restricted = 0 WHERE vendorId = ?"
	if _, err := v.Conn.ExecContext(ctx, query, id); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
