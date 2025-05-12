package suggestion_engine

import (
	"errors"
	"nearbyassist/internal/utils"
)

const (
	ERR_INVALID_WEIGHTS = "invalid weights"
)

type Weights struct {
	PriceWeight             float32
	RatingWeight            float32
	DistanceWeight          float32
	BookingsCompletedWeight float32
}

func (c *Weights) Validate() error {
	if c.PriceWeight+c.RatingWeight+c.DistanceWeight+c.BookingsCompletedWeight != 1 {
		return errors.New(ERR_INVALID_WEIGHTS)
	}

	return nil
}

func NewWeightsFromStrings(price, rating, distance, bookings string) *Weights {
	p := utils.StringToFloat32ElseZero(price)
	r := utils.StringToFloat32ElseZero(rating)
	d := utils.StringToFloat32ElseZero(distance)
	b := utils.StringToFloat32ElseZero(bookings)

	return &Weights{
		PriceWeight:             p,
		RatingWeight:            r,
		DistanceWeight:          d,
		BookingsCompletedWeight: b,
	}
}
