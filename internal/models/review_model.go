package models

type ReviewModel struct {
	Model
	BookingId string `json:"bookingId" db:"bookingId" validate:"required"`
	Rating    int    `json:"rating" db:"rating" validate:"required"`
	Text      string `json:"text" db:"text" validate:"required"`
}
