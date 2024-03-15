package review

import (
	review_query "nearbyassist/internal/db/query/review"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateReview(c echo.Context) error {
	review := new(types.Review)
	err := c.Bind(review)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err = c.Validate(review)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err = review_query.CreateReview(*review)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Review created successfully!",
	})
}
