package transaction

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/cache"
	service_service "nearbyassist/internal/service/service"
	transaction_service "nearbyassist/internal/service/transaction"
	user_service "nearbyassist/internal/service/user"
	"nearbyassist/internal/utils"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
)

type transactionHandler struct {
	transactionService *transaction_service.Service
	serviceService     *service_service.Service
	userService        *user_service.Service
}

func NewHandler(transactionService *transaction_service.Service, serviceService *service_service.Service, useService *user_service.Service) *transactionHandler {
	return &transactionHandler{
		transactionService: transactionService,
		serviceService:     serviceService,
		userService:        useService,
	}
}

func (h *transactionHandler) CreateTransaction(c echo.Context) error {
	req := new(request.NewTransactionPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	transactionId, err := h.transactionService.CreateTransaction(req)
	if err != nil {
		if strings.Contains(err.Error(), "You already have an confirmed or pending transaction for this service") {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "You already have an confirmed or pending transaction for this service",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating transaction",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transaction": transactionId,
	})
}

func (h *transactionHandler) GetTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Transaction ID is required",
			Error:   "Transaction ID is required",
		})
	}

	var transaction *models.TransactionModel
	inCache, exists := cache.NewGoCache().Get(c.Request().RequestURI)
	if exists {
		transaction = inCache.(*models.TransactionModel)
	} else {
		res, err := h.transactionService.GetTransaction(transactionId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting transaction",
				Error:   err.Error(),
			})
		}

		transaction = res
	}

	service, err := h.serviceService.GetService(transaction.ServiceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Could not get service information of transaction",
			Error:   err.Error(),
		})
	}

	response := response.Transaction{
		Id: transaction.Id,
		Vendor: response.User{
			Id:   transaction.VendorId,
			Name: transaction.Vendor,
		},
		Client: response.User{
			Id:   transaction.ClientId,
			Name: transaction.Client,
		},
		Cost: transaction.Cost,
		Extras: slices.Collect(utils.Map(transaction.Extras, func(x *models.ExtraModel) response.Extra {
			return response.Extra{
				Id:          x.Id,
				Title:       x.Title,
				Description: x.Description,
				Price:       x.Price,
			}
		})),
		Service: response.ServiceBareInfo{
			Id:          transaction.ServiceId,
			VendorId:    transaction.VendorId,
			Title:       service.Service.Title,
			Description: service.Service.Description,
			Rate:        service.Service.Rate,
			Tags: slices.Collect(
				utils.Map(service.Service.Tags, func(t *models.TagModel) response.Tag {
					return response.Tag{
						Id:    t.Id,
						Title: t.Title,
					}
				}),
			),
			Location: response.Location{
				Latitude:  service.Service.Latitude,
				Longitude: service.Service.Longitude,
			},
		},
		Status:       string(transaction.Status),
		CreatedAt:    transaction.CreatedAt,
		UpdatedAt:    transaction.UpdatedAt,
		ScheduledAt:  transaction.ScheduledAt.String,
		CancelReason: transaction.CancelReason.String,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) Cancel(c echo.Context) error {
	req := new(request.CancelRequestPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.transactionService.CancelTransaction(bearerToken, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error cancellation request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *transactionHandler) Accept(c echo.Context) error {
	req := new(request.AcceptTransactionPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.transactionService.AcceptTransactionRequest(bearerToken, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error accepting transaction request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *transactionHandler) Reject(c echo.Context) error {
	req := new(request.RejectRequestPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)
	if err := h.transactionService.RejectTransactionRequest(bearerToken, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error rejecting transaction request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *transactionHandler) GetUserTransactionList(c echo.Context) error {
	filter := c.QueryParam("filter")
	bearerToken := utils.BearerTokenFromHeader(c)

	transactions := make([]*models.TransactionModel, 0)
	switch filter {
	case "sent":
		if result, err := h.transactionService.GetTransactionUserSent(bearerToken); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting user transactions",
				Error:   err.Error(),
			})
		} else {
			transactions = result
		}
	case "received":
		if result, err := h.transactionService.GetTransactionUserReceived(bearerToken); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting user transactions",
				Error:   err.Error(),
			})
		} else {
			transactions = result
		}
	}

	resp := make([]response.Transaction, 0)
	for _, t := range transactions {
		service, err := h.serviceService.GetService(t.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of transaction",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Transaction{
			Id: t.Id,
			Vendor: response.User{
				Id:   t.VendorId,
				Name: t.Vendor,
			},
			Client: response.User{
				Id:   t.ClientId,
				Name: t.Client,
			},
			Cost: t.Cost,
			Extras: slices.Collect(utils.Map(t.Extras, func(x *models.ExtraModel) response.Extra {
				return response.Extra{
					Id:          x.Id,
					Title:       x.Title,
					Description: x.Description,
					Price:       x.Price,
				}
			})),
			Service: response.ServiceBareInfo{
				Id:          t.ServiceId,
				VendorId:    t.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags: slices.Collect(
					utils.Map(service.Service.Tags, func(t *models.TagModel) response.Tag {
						return response.Tag{
							Id:    t.Id,
							Title: t.Title,
						}
					}),
				),
				Location: response.Location{
					Latitude:  service.Service.Latitude,
					Longitude: service.Service.Longitude,
				},
			},
			Status:       string(t.Status),
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
			ScheduledAt:  t.ScheduledAt.String,
			CancelReason: t.CancelReason.String,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": resp,
	})
}

func (h *transactionHandler) GetRecentTransactions(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

	transactions, err := h.transactionService.GetRecentTransactions(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving recents",
			Error:   err.Error(),
		})
	}

	resp := make([]response.Transaction, 0)
	for _, t := range transactions {
		service, err := h.serviceService.GetService(t.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of transaction",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Transaction{
			Id: t.Id,
			Vendor: response.User{
				Id:   t.VendorId,
				Name: t.Vendor,
			},
			Client: response.User{
				Id:   t.ClientId,
				Name: t.Client,
			},
			Cost: t.Cost,
			Extras: slices.Collect(utils.Map(t.Extras, func(x *models.ExtraModel) response.Extra {
				return response.Extra{
					Id:          x.Id,
					Title:       x.Title,
					Description: x.Description,
					Price:       x.Price,
				}
			})),
			Service: response.ServiceBareInfo{
				Id:          t.ServiceId,
				VendorId:    t.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags: slices.Collect(
					utils.Map(service.Service.Tags, func(t *models.TagModel) response.Tag {
						return response.Tag{
							Id:    t.Id,
							Title: t.Title,
						}
					}),
				),
				Location: response.Location{
					Latitude:  service.Service.Latitude,
					Longitude: service.Service.Longitude,
				},
			},
			Status:       string(t.Status),
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
			ScheduledAt:  t.ScheduledAt.String,
			CancelReason: t.CancelReason.String,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": resp,
	})
}

func (h *transactionHandler) GetConfirmedTransactions(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

	transactions, err := h.transactionService.GetConfirmedTransactions(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting confirmed transactions",
			Error:   err.Error(),
		})
	}

	resp := make([]response.Transaction, 0)
	for _, t := range transactions {
		service, err := h.serviceService.GetService(t.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of transaction",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Transaction{
			Id: t.Id,
			Vendor: response.User{
				Id:   t.VendorId,
				Name: t.Vendor,
			},
			Client: response.User{
				Id:   t.ClientId,
				Name: t.Client,
			},
			Cost: t.Cost,
			Extras: slices.Collect(utils.Map(t.Extras, func(x *models.ExtraModel) response.Extra {
				return response.Extra{
					Id:          x.Id,
					Title:       x.Title,
					Description: x.Description,
					Price:       x.Price,
				}
			})),
			Service: response.ServiceBareInfo{
				Id:          t.ServiceId,
				VendorId:    t.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags: slices.Collect(
					utils.Map(service.Service.Tags, func(t *models.TagModel) response.Tag {
						return response.Tag{
							Id:    t.Id,
							Title: t.Title,
						}
					}),
				),
				Location: response.Location{
					Latitude:  service.Service.Latitude,
					Longitude: service.Service.Longitude,
				},
			},
			Status:       string(t.Status),
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
			ScheduledAt:  t.ScheduledAt.String,
			CancelReason: t.CancelReason.String,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": resp,
	})
}

func (h *transactionHandler) GetReviewableTransactions(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

	reviewables, err := h.transactionService.GetReviewableTransactions(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting confirmed transactions",
			Error:   err.Error(),
		})
	}

	resp := make([]response.Transaction, 0)
	for _, t := range reviewables {
		service, err := h.serviceService.GetService(t.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of transaction",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Transaction{
			Id: t.Id,
			Vendor: response.User{
				Id:   t.VendorId,
				Name: t.Vendor,
			},
			Client: response.User{
				Id:   t.ClientId,
				Name: t.Client,
			},
			Cost: t.Cost,
			Extras: slices.Collect(utils.Map(t.Extras, func(x *models.ExtraModel) response.Extra {
				return response.Extra{
					Id:          x.Id,
					Title:       x.Title,
					Description: x.Description,
					Price:       x.Price,
				}
			})),
			Service: response.ServiceBareInfo{
				Id:          t.ServiceId,
				VendorId:    t.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags: slices.Collect(
					utils.Map(service.Service.Tags, func(t *models.TagModel) response.Tag {
						return response.Tag{
							Id:    t.Id,
							Title: t.Title,
						}
					}),
				),
				Location: response.Location{
					Latitude:  service.Service.Latitude,
					Longitude: service.Service.Longitude,
				},
			},
			Status:       string(t.Status),
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
			ScheduledAt:  t.ScheduledAt.String,
			CancelReason: t.CancelReason.String,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"reviewables": resp,
	})
}

func (h *transactionHandler) GetTransactionHistory(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

	transactions, err := h.transactionService.GetTransactionHistory(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting transaction history",
			Error:   err.Error(),
		})
	}

	resp := make([]response.Transaction, 0)
	for _, t := range transactions {
		service, err := h.serviceService.GetService(t.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of transaction",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Transaction{
			Id: t.Id,
			Vendor: response.User{
				Id:   t.VendorId,
				Name: t.Vendor,
			},
			Client: response.User{
				Id:   t.ClientId,
				Name: t.Client,
			},
			Cost: t.Cost,
			Extras: slices.Collect(utils.Map(t.Extras, func(x *models.ExtraModel) response.Extra {
				return response.Extra{
					Id:          x.Id,
					Title:       x.Title,
					Description: x.Description,
					Price:       x.Price,
				}
			})),
			Service: response.ServiceBareInfo{
				Id:          t.ServiceId,
				VendorId:    t.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags: slices.Collect(
					utils.Map(service.Service.Tags, func(t *models.TagModel) response.Tag {
						return response.Tag{
							Id:    t.Id,
							Title: t.Title,
						}
					}),
				),
				Location: response.Location{
					Latitude:  service.Service.Latitude,
					Longitude: service.Service.Longitude,
				},
			},
			Status:       string(t.Status),
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
			ScheduledAt:  t.ScheduledAt.String,
			CancelReason: t.CancelReason.String,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"history": resp,
	})
}

func (h *transactionHandler) CompleteTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Transaction ID is required",
			Error:   "Transaction ID is required",
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.transactionService.CompleteTransaction(bearerToken, transactionId); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "You are not authorized to complete this transaction",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error marking transaction as complete",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
