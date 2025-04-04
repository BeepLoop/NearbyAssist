package suggestion_engine

import (
	"nearbyassist/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedScoring(t *testing.T) {
	t.Run("Test GetTopScore", func(t *testing.T) {
		testData := []*models.GeoSpatialSearchResult{
			{
				Rate:              100,
				Rating:            2.2,
				Distance:          100,
				CompletedBookings: 100,
			},
			{
				Rate:              100,
				Rating:            2.3,
				Distance:          110,
				CompletedBookings: 100,
			},
			{
				Rate:              10,
				Rating:            2.3,
				Distance:          99.9,
				CompletedBookings: 101,
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
