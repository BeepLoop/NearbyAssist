package models

type TransactionModel struct {
	Model
	UpdateableModel
	VendorId  int    `json:"vendorId" db:"vendorId"`
	ClientId  int    `json:"clientId" db:"clientId"`
	ServiceId int    `json:"serviceId" db:"serviceId"`
	Status    string `json:"status" db:"status"`
}

func (t *TransactionModel) Create() (int, error) {
	return 0, nil
}

func (t *TransactionModel) Update(id int) error {
	return nil
}

func (t *TransactionModel) Delete(id int) error {
	return nil
}
