package suggestion_engine

import (
	"nearbyassist/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedScoring(t *testing.T) {
	t.Run("Test GetTopScore", func(t *testing.T) {
		testData := []dto.GeospatialOperation{
			{
				Price:              100,
				Rating:             2.2,
				DistanceFromOrigin: 100,
				CompletedBookings:  100,
			},
			{
				Price:              100,
				Rating:             2.3,
				DistanceFromOrigin: 110,
				CompletedBookings:  100,
			},
			{
				Price:              10,
				Rating:             2.3,
				DistanceFromOrigin: 99.9,
				CompletedBookings:  101,
			},
		}

		w := NewWeightedScoring()
		w.getTopScores(testData)

		expected := score{
			lowestPrice:      10,
			highestRating:    2.3,
			shortestDistance: 99.9,
			mostBookings:     101,
		}

		assert.Equal(t, expected, w.score)
	})
}
