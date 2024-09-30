package suggestion_engine

import (
	"math"
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
)

type criteria struct {
	priceWeight        float32
	ratingWeight       float32
	distanceWeight     float32
	transactionsWeight float32
}

type score struct {
	lowestPrice      float32
	highestRating    float32
	shortestDistance float32
	mostTransactions float32
}

type weightedScoring struct {
	criteria
	score
}

func NewWeightedScoring() *weightedScoring {
	return &weightedScoring{
		criteria: criteria{
			priceWeight:        0.4,
			ratingWeight:       0.3,
			distanceWeight:     0.2,
			transactionsWeight: 0.1,
		},
	}
}

func (w *weightedScoring) GenerateSuggestions(services []*models.GeoSpatialSearchResult) ([]*response.SearchResult, error) {
	scores := make([]*response.SearchResult, 0)

	w.getTopScores(services)

	for _, service := range services {
		score, err := w.calculateScore(service)
		if err != nil {
			return nil, err
		}

		scores = append(scores, &response.SearchResult{
			Id:        service.Id,
			Score:     score,
			Rank:      0,
			Vendor:    service.VendorName,
			Latitude:  service.Latitude,
			Longitude: service.Longitude,
		})
	}

	return scores, nil
}

func (w *weightedScoring) calculateScore(service *models.GeoSpatialSearchResult) (float32, error) {
	priceScore := w.minimize(service.Rate, w.lowestPrice)
	ratingScore := w.maximize(service.Rating, w.highestRating)
	distanceScore := w.minimize(service.Distance, w.shortestDistance)
	transactionsScore := w.maximize(service.CompletedTransactions, w.mostTransactions)

	score := (w.priceWeight * priceScore) + (w.ratingWeight * ratingScore) + (w.distanceWeight * distanceScore) + (w.transactionsWeight * transactionsScore)
	return score, nil
}

func (w *weightedScoring) getTopScores(services []*models.GeoSpatialSearchResult) {
	var lowestPrice float32 = math.MaxFloat32
	var highestRating float32 = 1.0
	var shortestDistance float32 = math.MaxFloat32
	var mostTransactions float32 = 1.0

	for _, service := range services {
		if service.Rate < lowestPrice {
			lowestPrice = service.Rate
		}

		if service.Rating > highestRating {
			highestRating = service.Rating
		}

		if service.Distance < shortestDistance {
			shortestDistance = service.Distance
		}

		if service.CompletedTransactions > mostTransactions {
			mostTransactions = service.CompletedTransactions
		}
	}

	w.lowestPrice = lowestPrice
	w.highestRating = highestRating
	w.shortestDistance = shortestDistance
	w.mostTransactions = mostTransactions
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
	return minScore / score
}
