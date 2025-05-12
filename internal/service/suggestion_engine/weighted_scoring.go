package suggestion_engine

import (
	"math"
	"nearbyassist/internal/dto"
)

var (
	weightedScoringInstance *weightedScoring
)

type score struct {
	lowestPrice      float32
	highestRating    float32
	shortestDistance float32
	mostBookings     float32
}

type weightedScoring struct {
	Weights
	score
}

func NewWeightedScoring() *weightedScoring {
	if weightedScoringInstance != nil {
		return weightedScoringInstance
	}

	weightedScoringInstance = &weightedScoring{
		Weights: Weights{
			PriceWeight:             0.4,
			RatingWeight:            0.3,
			DistanceWeight:          0.2,
			BookingsCompletedWeight: 0.1,
		},
	}
	return weightedScoringInstance
}

func (w *weightedScoring) SetWeights(weights Weights) {
	w.Weights = weights
}

func (w *weightedScoring) GetValues() Weights {
	return w.Weights
}

func (w *weightedScoring) GenerateSuggestions(services []dto.GeospatialOperation) (map[string]float32, error) {
	w.getTopScores(services)

	scores := make(map[string]float32)
	for _, service := range services {
		score, err := w.calculateScore(service)
		if err != nil {
			return nil, err
		}

		scores[service.Id] = score
	}

	return scores, nil
}

func (w *weightedScoring) calculateScore(service dto.GeospatialOperation) (float32, error) {
	priceScore := w.minimize(service.Price, w.lowestPrice)
	ratingScore := w.maximize(service.Rating, w.highestRating)
	distanceScore := w.minimize(service.DistanceFromOrigin, w.shortestDistance)
	bookingsScore := w.maximize(service.CompletedBookings, w.mostBookings)

	score := (w.PriceWeight * priceScore) + (w.RatingWeight * ratingScore) + (w.DistanceWeight * distanceScore) + (w.BookingsCompletedWeight * bookingsScore)
	return score, nil
}

func (w *weightedScoring) getTopScores(services []dto.GeospatialOperation) {
	var lowestPrice float32 = math.MaxFloat32
	var highestRating float32 = math.SmallestNonzeroFloat32
	var shortestDistance float32 = math.MaxFloat32
	var mostBookings float32 = math.SmallestNonzeroFloat32

	for _, service := range services {
		if service.Price < lowestPrice {
			lowestPrice = service.Price
		}

		if service.Rating > highestRating {
			highestRating = service.Rating
		}

		if service.DistanceFromOrigin < shortestDistance {
			shortestDistance = service.DistanceFromOrigin
		}

		if service.CompletedBookings > mostBookings {
			mostBookings = service.CompletedBookings
		}
	}

	w.lowestPrice = lowestPrice
	w.highestRating = highestRating
	w.shortestDistance = shortestDistance
	w.mostBookings = mostBookings
}

func (w *weightedScoring) maximize(score, maxScore float32) float32 {
	if maxScore == 0 {
		return 0
	}
	return score / maxScore
}

func (w *weightedScoring) minimize(score, minScore float32) float32 {
	if score == 0 {
		return 0
	}

	if score == 0 && minScore == 0 {
		return 1
	}

	return minScore / score
}
