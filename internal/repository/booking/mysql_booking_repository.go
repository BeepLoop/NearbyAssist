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
            t.id = ?
    `

	booking := new(models.BookingModel)
	if err := s.db.GetContext(ctx, booking, bookingQuery, id); err != nil {
		fmt.Println("error get booking by id: ", err.Error())
		return nil, err
	}

	if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
		return nil, err
	} else {
		booking.IsReviewed = isReviewed
	}

	if service, err := s.getService(booking.ServiceId); err != nil {
		return nil, err
	} else {
		booking.Service = service
	}

	if extras, err := s.getBookingExtras(booking.Id); err != nil {
		return nil, err
	} else {
		booking.Extras = extras
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return booking, nil
}

func (s *MysqlBookingRepository) GetAll(limit, offset int) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT * FROM Booking ORDER BY createdAt, updatedAt DESC LIMIT ? OFFSET ?"

	bookings := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &bookings, query, limit, offset); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) HasOngoingBookingForService(data *models.BookingModel) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	if err := s.db.GetContext(ctx, &alreadyBooked, checkDuplicate, data.ClientId, data.ServiceId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return alreadyBooked, nil
}

func (s *MysqlBookingRepository) GetConfirmedBookingsOfVendor(vendorId string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            id
        FROM
            Booking
        WHERE
            vendorId = ? AND status = 'confirmed'
    `

	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query, vendorId); err != nil {
		return nil, err
	}

	bookings := make([]*models.BookingModel, 0)
	for _, id := range ids {
		booking, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}

func (s *MysqlBookingRepository) GetBookingSent(id string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
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
	sent := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &sent, query, id); err != nil {
		return nil, err
	}

	bookings := make([]*models.BookingModel, 0)
	for _, b := range sent {
		booking, err := s.FindById(b.Id)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
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
	received := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &received, bookingQuery, id); err != nil {
		return nil, err
	}

	bookings := make([]*models.BookingModel, 0)
	for _, b := range received {
		booking, err := s.FindById(b.Id)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
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
	recents := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &recents, bookingQuery, userId, userId); err != nil {
		return nil, err
	}

	bookings := make([]*models.BookingModel, 0)
	for _, b := range recents {
		booking, err := s.FindById(b.Id)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
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

	confirmedBookings := make([]*models.BookingModel, 0)
	switch filter {
	case "vendor":
		err := s.db.SelectContext(ctx, &confirmedBookings, queryForVendor, id)
		if err != nil {
			return nil, err
		}
	case "client":
		err := s.db.SelectContext(ctx, &confirmedBookings, queryForClient, id)
		if err != nil {
			return nil, err
		}
	}

	bookings := make([]*models.BookingModel, 0)
	for _, b := range confirmedBookings {
		booking, err := s.FindById(b.Id)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
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

	history := make([]*models.BookingModel, 0)
	switch filter {
	case "vendor":
		err := s.db.SelectContext(ctx, &history, queryForVendor, id)
		if err != nil {
			return nil, err
		}
	case "client":
		err := s.db.SelectContext(ctx, &history, queryForClient, id)
		if err != nil {
			return nil, err
		}
	}

	bookings := make([]*models.BookingModel, 0)
	for _, b := range history {
		booking, err := s.FindById(b.Id)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
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

	reviewables := make([]*models.BookingModel, 0)
	err := s.db.SelectContext(ctx, &reviewables, bookingQuery, userId)
	if err != nil {
		return nil, err
	}

	bookings := make([]*models.BookingModel, 0)
	for _, b := range reviewables {
		booking, err := s.FindById(b.Id)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
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

func (s *MysqlBookingRepository) getBookingExtras(bookingId string) ([]*models.ExtraModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	if err := s.db.SelectContext(ctx, &extras, extrasQuery, bookingId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return extras, nil
}

func (s *MysqlBookingRepository) getService(serviceId string) (*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	service := new(models.ServiceModel)

	query := `
        SELECT
            id,
            vendorId,
            title,
            description,
            format(rate, 2) as rate,
            createdAt,
            updatedAt,
            disabled
        FROM 
            Service
        WHERE
            id = ?
    `
	if err := s.db.GetContext(ctx, service, query, serviceId); err != nil {
		return nil, err
	}

	if address, err := s.getAddress(serviceId); err != nil {
		return nil, err
	} else {
		service.Address = *address
	}

	if extras, err := s.getExtras(serviceId); err != nil {
		return nil, err
	} else {
		service.Extras = extras
	}

	if tags, err := s.getTags(serviceId); err != nil {
		return nil, err
	} else {
		service.Tags = tags
	}

	if images, err := s.getPhotos(serviceId); err != nil {
		return nil, err
	} else {
		service.Images = images
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return service, nil
}

func (s *MysqlBookingRepository) getTags(serviceId string) ([]*models.TagModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.id,
            t.title
        FROM
            ServiceTag st
            JOIN Tag t ON t.id = st.tagId
        WHERE
            st.serviceId = ?;
    `

	tags := make([]*models.TagModel, 0)
	if err := s.db.SelectContext(ctx, &tags, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return tags, nil
}

func (s *MysqlBookingRepository) getPhotos(serviceId string) ([]*models.ServicePhotoModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, serviceId, url, vendorId FROM ServicePhoto WHERE serviceId = ?"

	images := make([]*models.ServicePhotoModel, 0)
	if err := s.db.SelectContext(ctx, &images, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return images, nil
}

func (s *MysqlBookingRepository) getAddress(serviceId string) (*models.AddressModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            a.id, a.address, a.latitude, a.longitude
        FROM
            Address a
            JOIN ServiceAddress sa ON sa.addressId = a.id
        WHERE
            sa.serviceId = ?
    `
	address := new(models.AddressModel)
	if err := s.db.GetContext(ctx, address, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return address, nil
}

func (s *MysqlBookingRepository) getExtras(serviceId string) ([]*models.ExtraModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	extrasQuery := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN ServiceExtra se ON se.extraId = e.id
        WHERE
            se.serviceId = ? AND e.deleted = 0
    `

	extras := make([]*models.ExtraModel, 0)
	if err := s.db.SelectContext(ctx, &extras, extrasQuery, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return extras, nil
}
