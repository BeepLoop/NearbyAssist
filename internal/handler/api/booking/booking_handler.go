package booking

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	booking_service "nearbyassist/internal/service/booking"
	"nearbyassist/internal/service/cache"
	qr_service "nearbyassist/internal/service/qr"
	service_service "nearbyassist/internal/service/service"
	user_service "nearbyassist/internal/service/user"
	"nearbyassist/internal/utils"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
)

type bookingHandler struct {
	bookingService *booking_service.Service
	serviceService *service_service.Service
	userService    *user_service.Service
	qrService      *qr_service.Service
}

func NewHandler(bookingService *booking_service.Service, serviceService *service_service.Service, useService *user_service.Service, qrService *qr_service.Service) *bookingHandler {
	return &bookingHandler{
		bookingService: bookingService,
		serviceService: serviceService,
		userService:    useService,
		qrService:      qrService,
	}
}

func (h *bookingHandler) CreateBooking(c echo.Context) error {
	req := new(request.NewBookingPayload)
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

	bookingId, err := h.bookingService.CreateBooking(req)
	if err != nil {
		if strings.Contains(err.Error(), booking_service.ERR_DISABLED_SERVICE) {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "Service is disabled, booking not allowed",
				Error:   err.Error(),
			})
		}

		if strings.Contains(err.Error(), booking_service.ERR_HAS_PENDING_OR_CONFIRMED) {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "You already have an confirmed or pending booking for this service",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating booking",
			Error:   err.Error(),
		})
	}

	booking, err := h.bookingService.GetBooking(bookingId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve booking information",
			Error:   err.Error(),
		})
	}

	service, err := h.serviceService.GetService(booking.ServiceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Could not get service information of booking",
			Error:   err.Error(),
		})
	}

	response := response.Booking{
		Id: booking.Id,
		Vendor: response.User{
			Id:   booking.VendorId,
			Name: booking.Vendor,
		},
		Client: response.User{
			Id:   booking.ClientId,
			Name: booking.Client,
		},
		Cost: booking.Cost,
		Extras: slices.AppendSeq(
			make([]response.Extra, 0),
			utils.Map(booking.Extras, func(x *models.ExtraModel) response.Extra {
				return response.Extra{
					Id:          x.Id,
					Title:       x.Title,
					Description: x.Description,
					Price:       x.Price,
				}
			}),
		),
		Service: response.ServiceBareInfo{
			Id:          booking.ServiceId,
			VendorId:    booking.VendorId,
			Title:       service.Service.Title,
			Description: service.Service.Description,
			Rate:        service.Service.Rate,
			Tags:        service.Service.Tags,
			Location:    service.Service.Location,
		},
		Status:       string(booking.Status),
		CreatedAt:    booking.CreatedAt,
		UpdatedAt:    booking.UpdatedAt,
		ScheduledAt:  booking.ScheduledAt.String,
		CancelReason: booking.CancelReason.String,
		QRSignature: utils.Must(h.qrService.SignData(&request.QRSignatureInput{
			ClientID:  booking.ClientId,
			VendorID:  booking.VendorId,
			BookingID: booking.Id,
		})),
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"booking": response,
	})
}

func (h *bookingHandler) GetBooking(c echo.Context) error {
	bookingId := c.Param("bookingId")
	if bookingId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Booking ID is required",
			Error:   "Booking ID is required",
		})
	}

	var booking *models.BookingModel
	inCache, exists := cache.NewGoCache().Get(c.Request().RequestURI)
	if exists {
		booking = inCache.(*models.BookingModel)
	} else {
		res, err := h.bookingService.GetBooking(bookingId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting booking",
				Error:   err.Error(),
			})
		}

		booking = res
	}

	service, err := h.serviceService.GetService(booking.ServiceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Could not get service information of booking",
			Error:   err.Error(),
		})
	}

	response := response.Booking{
		Id: booking.Id,
		Vendor: response.User{
			Id:   booking.VendorId,
			Name: booking.Vendor,
		},
		Client: response.User{
			Id:   booking.ClientId,
			Name: booking.Client,
		},
		Cost: booking.Cost,
		Extras: slices.AppendSeq(
			make([]response.Extra, 0),
			utils.Map(booking.Extras, func(x *models.ExtraModel) response.Extra {
				return response.Extra{
					Id:          x.Id,
					Title:       x.Title,
					Description: x.Description,
					Price:       x.Price,
				}
			}),
		),
		Service: response.ServiceBareInfo{
			Id:          booking.ServiceId,
			VendorId:    booking.VendorId,
			Title:       service.Service.Title,
			Description: service.Service.Description,
			Rate:        service.Service.Rate,
			Tags:        service.Service.Tags,
			Location:    service.Service.Location,
		},
		Status:       string(booking.Status),
		CreatedAt:    booking.CreatedAt,
		UpdatedAt:    booking.UpdatedAt,
		ScheduledAt:  booking.ScheduledAt.String,
		CancelReason: booking.CancelReason.String,
		QRSignature: utils.Must(h.qrService.SignData(&request.QRSignatureInput{
			ClientID:  booking.ClientId,
			VendorID:  booking.VendorId,
			BookingID: booking.Id,
		})),
	}

	return c.JSON(http.StatusOK, response)
}

func (h *bookingHandler) Cancel(c echo.Context) error {
	req := new(request.CancelRequestPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.bookingService.CancelBooking(bearerToken, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error cancellation request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *bookingHandler) Accept(c echo.Context) error {
	req := new(request.AcceptBookingPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.bookingService.AcceptBookingRequest(bearerToken, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error accepting booking request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *bookingHandler) Reject(c echo.Context) error {
	req := new(request.RejectRequestPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)
	if err := h.bookingService.RejectBookingRequest(bearerToken, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error rejecting booking request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *bookingHandler) GetUserBookingList(c echo.Context) error {
	filter := c.QueryParam("filter")
	bearerToken := utils.BearerTokenFromHeader(c)

	bookings := make([]*models.BookingModel, 0)
	switch filter {
	case "sent":
		if result, err := h.bookingService.GetBookingUserSent(bearerToken); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting user bookings",
				Error:   err.Error(),
			})
		} else {
			bookings = result
		}
	case "received":
		if result, err := h.bookingService.GetBookingUserReceived(bearerToken); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting user bookings",
				Error:   err.Error(),
			})
		} else {
			bookings = result
		}
	}

	resp := make([]response.Booking, 0)
	for _, booking := range bookings {
		service, err := h.serviceService.GetService(booking.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of booking",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Booking{
			Id: booking.Id,
			Vendor: response.User{
				Id:   booking.VendorId,
				Name: booking.Vendor,
			},
			Client: response.User{
				Id:   booking.ClientId,
				Name: booking.Client,
			},
			Cost: booking.Cost,
			Extras: slices.AppendSeq(
				make([]response.Extra, 0),
				utils.Map(booking.Extras, func(x *models.ExtraModel) response.Extra {
					return response.Extra{
						Id:          x.Id,
						Title:       x.Title,
						Description: x.Description,
						Price:       x.Price,
					}
				}),
			),
			Service: response.ServiceBareInfo{
				Id:          booking.ServiceId,
				VendorId:    booking.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags:        service.Service.Tags,
				Location:    service.Service.Location,
			},
			Status:       string(booking.Status),
			CreatedAt:    booking.CreatedAt,
			UpdatedAt:    booking.UpdatedAt,
			ScheduledAt:  booking.ScheduledAt.String,
			CancelReason: booking.CancelReason.String,
			QRSignature: utils.Must(h.qrService.SignData(&request.QRSignatureInput{
				ClientID:  booking.ClientId,
				VendorID:  booking.VendorId,
				BookingID: booking.Id,
			})),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"bookings": resp,
	})
}

func (h *bookingHandler) GetRecentBookings(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

	bookings, err := h.bookingService.GetRecentBookings(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving recents",
			Error:   err.Error(),
		})
	}

	resp := make([]response.Booking, 0)
	for _, booking := range bookings {
		service, err := h.serviceService.GetService(booking.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of booking",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Booking{
			Id: booking.Id,
			Vendor: response.User{
				Id:   booking.VendorId,
				Name: booking.Vendor,
			},
			Client: response.User{
				Id:   booking.ClientId,
				Name: booking.Client,
			},
			Cost: booking.Cost,
			Extras: slices.AppendSeq(
				make([]response.Extra, 0),
				utils.Map(booking.Extras, func(x *models.ExtraModel) response.Extra {
					return response.Extra{
						Id:          x.Id,
						Title:       x.Title,
						Description: x.Description,
						Price:       x.Price,
					}
				}),
			),
			Service: response.ServiceBareInfo{
				Id:          booking.ServiceId,
				VendorId:    booking.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags:        service.Service.Tags,
				Location:    service.Service.Location,
			},
			Status:       string(booking.Status),
			CreatedAt:    booking.CreatedAt,
			UpdatedAt:    booking.UpdatedAt,
			ScheduledAt:  booking.ScheduledAt.String,
			CancelReason: booking.CancelReason.String,
			QRSignature: utils.Must(h.qrService.SignData(&request.QRSignatureInput{
				ClientID:  booking.ClientId,
				VendorID:  booking.VendorId,
				BookingID: booking.Id,
			})),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"bookings": resp,
	})
}

func (h *bookingHandler) GetConfirmedBookings(c echo.Context) error {
	filter := c.QueryParam("filter")
	bearerToken := utils.BearerTokenFromHeader(c)

	bookings, err := h.bookingService.GetConfirmedBookings(bearerToken, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting confirmed bookings",
			Error:   err.Error(),
		})
	}

	resp := make([]response.Booking, 0)
	for _, booking := range bookings {
		service, err := h.serviceService.GetService(booking.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of booking",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Booking{
			Id: booking.Id,
			Vendor: response.User{
				Id:   booking.VendorId,
				Name: booking.Vendor,
			},
			Client: response.User{
				Id:   booking.ClientId,
				Name: booking.Client,
			},
			Cost: booking.Cost,
			Extras: slices.AppendSeq(
				make([]response.Extra, 0),
				utils.Map(booking.Extras, func(x *models.ExtraModel) response.Extra {
					return response.Extra{
						Id:          x.Id,
						Title:       x.Title,
						Description: x.Description,
						Price:       x.Price,
					}
				}),
			),
			Service: response.ServiceBareInfo{
				Id:          booking.ServiceId,
				VendorId:    booking.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags:        service.Service.Tags,
				Location:    service.Service.Location,
			},
			Status:       string(booking.Status),
			CreatedAt:    booking.CreatedAt,
			UpdatedAt:    booking.UpdatedAt,
			ScheduledAt:  booking.ScheduledAt.String,
			CancelReason: booking.CancelReason.String,
			QRSignature: utils.Must(h.qrService.SignData(&request.QRSignatureInput{
				ClientID:  booking.ClientId,
				VendorID:  booking.VendorId,
				BookingID: booking.Id,
			})),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"bookings": resp,
	})
}

func (h *bookingHandler) GetReviewableBookings(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

	reviewables, err := h.bookingService.GetReviewableBookings(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting confirmed bookings",
			Error:   err.Error(),
		})
	}

	resp := make([]response.Booking, 0)
	for _, booking := range reviewables {
		service, err := h.serviceService.GetService(booking.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of booking",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Booking{
			Id: booking.Id,
			Vendor: response.User{
				Id:   booking.VendorId,
				Name: booking.Vendor,
			},
			Client: response.User{
				Id:   booking.ClientId,
				Name: booking.Client,
			},
			Cost: booking.Cost,
			Extras: slices.AppendSeq(
				make([]response.Extra, 0),
				utils.Map(booking.Extras, func(x *models.ExtraModel) response.Extra {
					return response.Extra{
						Id:          x.Id,
						Title:       x.Title,
						Description: x.Description,
						Price:       x.Price,
					}
				}),
			),
			Service: response.ServiceBareInfo{
				Id:          booking.ServiceId,
				VendorId:    booking.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags:        service.Service.Tags,
				Location:    service.Service.Location,
			},
			Status:       string(booking.Status),
			CreatedAt:    booking.CreatedAt,
			UpdatedAt:    booking.UpdatedAt,
			ScheduledAt:  booking.ScheduledAt.String,
			CancelReason: booking.CancelReason.String,
			QRSignature: utils.Must(h.qrService.SignData(&request.QRSignatureInput{
				ClientID:  booking.ClientId,
				VendorID:  booking.VendorId,
				BookingID: booking.Id,
			})),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"reviewables": resp,
	})
}

func (h *bookingHandler) GetBookingHistory(c echo.Context) error {
	filter := c.QueryParam("filter")
	bearerToken := utils.BearerTokenFromHeader(c)

	bookings, err := h.bookingService.GetBookingHistory(bearerToken, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting booking history",
			Error:   err.Error(),
		})
	}

	resp := make([]response.Booking, 0)
	for _, booking := range bookings {
		service, err := h.serviceService.GetService(booking.ServiceId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Could not get service information of booking",
				Error:   err.Error(),
			})
		}

		resp = append(resp, response.Booking{
			Id: booking.Id,
			Vendor: response.User{
				Id:   booking.VendorId,
				Name: booking.Vendor,
			},
			Client: response.User{
				Id:   booking.ClientId,
				Name: booking.Client,
			},
			Cost: booking.Cost,
			Extras: slices.AppendSeq(
				make([]response.Extra, 0),
				utils.Map(booking.Extras, func(x *models.ExtraModel) response.Extra {
					return response.Extra{
						Id:          x.Id,
						Title:       x.Title,
						Description: x.Description,
						Price:       x.Price,
					}
				}),
			),
			Service: response.ServiceBareInfo{
				Id:          booking.ServiceId,
				VendorId:    booking.VendorId,
				Title:       service.Service.Title,
				Description: service.Service.Description,
				Rate:        service.Service.Rate,
				Tags:        service.Service.Tags,
				Location:    service.Service.Location,
			},
			Status:       string(booking.Status),
			CreatedAt:    booking.CreatedAt,
			UpdatedAt:    booking.UpdatedAt,
			ScheduledAt:  booking.ScheduledAt.String,
			CancelReason: booking.CancelReason.String,
			QRSignature: utils.Must(h.qrService.SignData(&request.QRSignatureInput{
				ClientID:  booking.ClientId,
				VendorID:  booking.VendorId,
				BookingID: booking.Id,
			})),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"history": resp,
	})
}

func (h *bookingHandler) CompleteBooking(c echo.Context) error {
	bookingId := c.Param("bookingId")
	if bookingId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Booking ID is required",
			Error:   "Booking ID is required",
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.bookingService.CompleteBooking(bearerToken, bookingId); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "You are not authorized to complete this booking",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error marking booking as complete",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
