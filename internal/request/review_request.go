package request

type NewReview struct {
	ServiceId int    `json:"serviceId" db:"serviceId" validate:"required"`
	Rating    string `json:"rating" db:"rating" validate:"required"`
}
