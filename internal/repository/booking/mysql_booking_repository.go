package booking_repo

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlBookingRepository struct {
	db *sqlx.DB
}

func NewMysqlBookingRepository(db *sqlx.DB) *MysqlBookingRepository {
	return &MysqlBookingRepository{
		db: db,
	}
}

func (s *MysqlBookingRepository) Create(data *models.BookingModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	checkDuplicate := `
        SELECT CASE
            WHEN EXISTS
                (
                    SELECT
                        1
                    FROM
                        Booking
                    WHERE
                        (clientId = ? AND serviceId = ?)
                        AND (status = 'confirmed' OR status = 'pending')
                )
            THEN 1
            ELSE 0
        END AS duplicate_booking
    `
	alreadyBooked := false
	if err := tx.GetContext(ctx, &alreadyBooked, checkDuplicate, data.ClientId, data.ServiceId); err != nil {
		return "", err
	}

	if alreadyBooked {
		return "", errors.New("You already have an confirmed or pending booking for this service")
	}

	query := `
        INSERT INTO
            Booking (id, vendorId, clientId, serviceId, cost)
        VALUES
            (:id, :vendorId, :clientId, :serviceId, :cost)
    `

	if _, err := tx.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	insertBookingExtras := `
        INSERT INTO
            BookingExtra (bookingId, extraId)
        VALUES 
            (?, ?)
    `
	for _, extra := range data.Extras {
		if _, err := tx.ExecContext(ctx, insertBookingExtras, data.Id, extra.Id); err != nil {
			fmt.Println("extra: ", extra)
			fmt.Println("error: ", err.Error())
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlBookingRepository) FindById(id string) (*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bookingQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.createdAt,
            t.scheduledAt,
            t.cancelReason,
            t.updatedAt,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.id = ?
    `

	booking := new(models.BookingModel)
	if err := s.db.GetContext(ctx, booking, bookingQuery, id); err != nil {
		fmt.Println("error get booking: ", err.Error())
		return nil, err
	}

	if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
		return nil, err
	} else {
		booking.IsReviewed = isReviewed
	}

	serviceQuery := "SELECT * FROM Service WHERE id = ?"

	service := new(models.ServiceModel)
	if err := s.db.GetContext(ctx, service, serviceQuery, booking.ServiceId); err != nil {
		fmt.Println("error get service: ", err.Error())
		return nil, err
	}
	booking.Service = service

	extrasQuery := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            BookingExtra te
            JOIN Extra e ON e.id = te.extraId
        WHERE
            te.bookingId = ?
    `

	extras := make([]*models.ExtraModel, 0)
	if err := s.db.SelectContext(ctx, &extras, extrasQuery, booking.Id); err != nil {
		fmt.Println("error get extra: ", err.Error())
		return nil, err
	}
	booking.Extras = extras

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return booking, nil
}

func (s *MysqlBookingRepository) GetAll() ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bookings := make([]*models.BookingModel, 0)

	query := "SELECT * FROM Booking"
	if err := s.db.SelectContext(ctx, &bookings, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) GetConfirmedBookingsOfVendor(vendorId string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            *
        FROM
            Booking
        WHERE
            vendorId = ? AND status = 'confirmed'
    `

	bookings := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &bookings, query, vendorId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) GetBookingSent(id string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bookingQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.createdAt,
            t.updatedAt,
            t.scheduledAt,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND (t.status = 'pending' OR t.status = 'confirmed')
        ORDER BY
            t.updatedAt DESC
    `

	bookings := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &bookings, bookingQuery, id); err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Booking t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN BookingExtra te ON e.id = te.extraId
        WHERE
            te.bookingId = ?
    `

	for _, booking := range bookings {
		if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
			return nil, err
		} else {
			booking.IsReviewed = isReviewed
		}

		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, booking.Id); err != nil {
			return nil, err
		}
		booking.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, booking.Id); err != nil {
			return nil, err
		}
		booking.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) GetBookingReceived(id string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bookingQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.createdAt,
            t.updatedAt,
            t.scheduledAt,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? AND t.status = 'pending'
        ORDER BY
            t.updatedAt DESC
    `

	bookings := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &bookings, bookingQuery, id); err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Booking t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN BookingExtra te ON e.id = te.extraId
        WHERE
            te.bookingId = ?
    `

	for _, booking := range bookings {
		if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
			return nil, err
		} else {
			booking.IsReviewed = isReviewed
		}

		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, booking.Id); err != nil {
			return nil, err
		}
		booking.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, booking.Id); err != nil {
			return nil, err
		}
		booking.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) GetRecent(userId string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookingQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? OR t.clientId = ?
        ORDER BY
            t.updatedAt DESC
        LIMIT
            10
    `

	bookings := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &bookings, bookingQuery, userId, userId); err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Booking t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN BookingExtra te ON e.id = te.extraId
        WHERE
            te.bookingId = ?
    `

	for _, booking := range bookings {
		if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
			return nil, err
		} else {
			booking.IsReviewed = isReviewed
		}

		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, booking.Id); err != nil {
			return nil, err
		}
		booking.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, booking.Id); err != nil {
			return nil, err
		}
		booking.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) GetConfirmed(id, filter string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	allowedFilters := []string{"vendor", "client"}
	if !slices.Contains(allowedFilters, filter) {
		return nil, errors.New("invalid filter")
	}

	queryForVendor := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.createdAt,
            t.updatedAt,
            t.scheduledAt,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? AND t.status = 'confirmed'
        ORDER BY
            t.updatedAt DESC
    `

	queryForClient := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.createdAt,
            t.updatedAt,
            t.scheduledAt,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND t.status = 'confirmed'
        ORDER BY
            t.updatedAt DESC
    `

	bookings := make([]*models.BookingModel, 0)
	switch filter {
	case "vendor":
		err := s.db.SelectContext(ctx, &bookings, queryForVendor, id)
		if err != nil {
			return nil, err
		}
	case "client":
		err := s.db.SelectContext(ctx, &bookings, queryForClient, id)
		if err != nil {
			return nil, err
		}
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Booking t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN BookingExtra te ON e.id = te.extraId
        WHERE
            te.bookingId = ?
    `

	for _, booking := range bookings {
		if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
			return nil, err
		} else {
			booking.IsReviewed = isReviewed
		}

		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, booking.Id); err != nil {
			return nil, err
		}
		booking.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, booking.Id); err != nil {
			return nil, err
		}
		booking.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) GetHistory(id, filter string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	allowedFilters := []string{"vendor", "client"}
	if !slices.Contains(allowedFilters, filter) {
		return nil, errors.New("invalid filter")
	}

	queryForVendor := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.createdAt,
            t.updatedAt,
            t.scheduledAt,
            t.cancelReason,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? AND (t.status = 'done' OR t.status = 'cancelled')
        ORDER BY
            t.updatedAt DESC
    `

	queryForClient := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.createdAt,
            t.updatedAt,
            t.scheduledAt,
            t.cancelReason,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND (t.status = 'done' OR t.status = 'cancelled')
        ORDER BY
            t.updatedAt DESC
    `

	bookings := make([]*models.BookingModel, 0)
	switch filter {
	case "vendor":
		err := s.db.SelectContext(ctx, &bookings, queryForVendor, id)
		if err != nil {
			return nil, err
		}
	case "client":
		err := s.db.SelectContext(ctx, &bookings, queryForClient, id)
		if err != nil {
			return nil, err
		}
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Booking t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN BookingExtra te ON e.id = te.extraId
        WHERE
            te.bookingId = ?
    `

	for _, booking := range bookings {
		if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
			return nil, err
		} else {
			booking.IsReviewed = isReviewed
		}

		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, booking.Id); err != nil {
			return nil, err
		}
		booking.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, booking.Id); err != nil {
			return nil, err
		}
		booking.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) GetReviewableBookings(userId string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bookingQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.createdAt,
            t.updatedAt,
            t.scheduledAt,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND t.status = 'done' AND NOT EXISTS (
                SELECT 1 FROM Review r WHERE r.bookingId = t.id
            )
        ORDER BY
            t.updatedAt DESC
    `

	bookings := make([]*models.BookingModel, 0)
	err := s.db.SelectContext(ctx, &bookings, bookingQuery, userId)
	if err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Booking t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN BookingExtra te ON e.id = te.extraId
        WHERE
            te.bookingId = ?
    `

	for _, booking := range bookings {
		if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
			return nil, err
		} else {
			booking.IsReviewed = isReviewed
		}

		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, booking.Id); err != nil {
			return nil, err
		}
		booking.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, booking.Id); err != nil {
			return nil, err
		}
		booking.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) Cancel(bookingId, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE 
            Booking 
        SET 
            cancelReason = ?, status = 'cancelled'
        WHERE 
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, reason, bookingId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlBookingRepository) Accept(bookingId, schedule string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Booking
        SET
            scheduledAt = ?, status = 'confirmed'
        WHERE
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, schedule, bookingId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlBookingRepository) Reject(bookingId, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Booking
        SET
            cancelReason = ?,
            status = 'rejected'
        WHERE 
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, reason, bookingId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlBookingRepository) MarkComplete(bookingId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE Booking SET status = 'done' WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, bookingId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlBookingRepository) IsReviewed(bookingId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT CASE
            WHEN EXISTS (SELECT 1 FROM Review WHERE bookingId = ?)
            THEN 1
            ELSE 0
        END AS is_reviewed
    `

	isReviewed := false
	if err := s.db.GetContext(ctx, &isReviewed, query, bookingId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isReviewed, nil
}
