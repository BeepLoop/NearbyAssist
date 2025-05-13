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

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO
            Booking (id, vendorId, clientId, serviceId, serviceTitle, serviceDescription, price, pricingType, quantity, cost)
        VALUES
            (:id, :vendorId, :clientId, :serviceId, :serviceTitle, :serviceDescription, :price, :pricingType, :quantity, :cost)
    `

	data.Id = utils.GenerateId()
	if _, err := tx.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	insertBookingExtras := `
        INSERT INTO
            BookingExtra (bookingId, extraTitle, extraDescription, price)
        VALUES 
            (:bookingId, :extraTitle, :extraDescription, :price)
    `
	for _, extra := range data.Extras {
		extra.BookingId = data.Id
		if _, err := tx.NamedExecContext(ctx, insertBookingExtras, extra); err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
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
            t.serviceTitle,
            t.serviceDescription,
            t.price,
            t.pricingType,
            t.status,
            t.quantity,
            t.cost,
            t.createdAt,
            t.updatedAt,
            t.scheduleStart,
            t.scheduleEnd,
            t.cancelReason,
            t.cancelledBy
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

	if client, err := s.getClient(booking.Id); err != nil {
		return nil, err
	} else {
		booking.Client = *client
	}

	if vendor, err := s.getVendor(booking.Id); err != nil {
		return nil, err
	} else {
		booking.Vendor = *vendor
	}

	if isReviewed, err := s.IsReviewed(booking.Id); err != nil {
		return nil, err
	} else {
		booking.IsReviewed = isReviewed
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
            t.id
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND t.status = 'pending'
        ORDER BY
            t.updatedAt DESC
    `
	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query, id); err != nil {
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

func (s *MysqlBookingRepository) GetBookingReceived(id string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bookingQuery := `
        SELECT
            t.id
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? AND t.status = 'pending'
        ORDER BY
            t.updatedAt DESC
    `
	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, bookingQuery, id); err != nil {
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

func (s *MysqlBookingRepository) GetRecent(userId string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookingQuery := `
        SELECT
            t.id
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
	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, bookingQuery, userId, userId); err != nil {
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

func (s *MysqlBookingRepository) GetConfirmed(id, filter string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	allowedFilters := []string{"vendor", "client"}
	if !slices.Contains(allowedFilters, filter) {
		return nil, errors.New("invalid filter")
	}

	queryForVendor := `
        SELECT
            t.id
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
            t.id
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND t.status = 'confirmed'
        ORDER BY
            t.updatedAt DESC
    `

	ids := make([]string, 0)
	switch filter {
	case "vendor":
		if err := s.db.SelectContext(ctx, &ids, queryForVendor, id); err != nil {
			return nil, err
		}
	case "client":
		if err := s.db.SelectContext(ctx, &ids, queryForClient, id); err != nil {
			return nil, err
		}
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

func (s *MysqlBookingRepository) GetHistory(id, filter string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	allowedFilters := []string{"vendor", "client"}
	if !slices.Contains(allowedFilters, filter) {
		return nil, errors.New("invalid filter")
	}

	queryForVendor := `
        SELECT
            t.id
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? AND (t.status = 'done' OR t.status = 'cancelled' OR t.status = 'rejected')
        ORDER BY
            t.updatedAt DESC
    `

	queryForClient := `
        SELECT
            t.id
        FROM 
            Booking t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND (t.status = 'done' OR t.status = 'cancelled' OR t.status = 'rejected')
        ORDER BY
            t.updatedAt DESC
    `

	ids := make([]string, 0)
	switch filter {
	case "vendor":
		if err := s.db.SelectContext(ctx, &ids, queryForVendor, id); err != nil {
			return nil, err
		}
	case "client":
		if err := s.db.SelectContext(ctx, &ids, queryForClient, id); err != nil {
			return nil, err
		}
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

func (s *MysqlBookingRepository) GetReviewableBookings(userId string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bookingQuery := `
        SELECT
            t.id
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

	ids := make([]string, 0)
	err := s.db.SelectContext(ctx, &ids, bookingQuery, userId)
	if err != nil {
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

func (s *MysqlBookingRepository) Cancel(bookingId, cancelledBy, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE 
            Booking 
        SET 
            cancelReason = ?, status = 'cancelled', cancelledBy = ?
        WHERE 
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, reason, cancelledBy, bookingId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlBookingRepository) Accept(bookingId string, scheduleStart, scheduleEnd time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Booking
        SET
            scheduleStart = ?, scheduleEnd = ?, status = 'confirmed'
        WHERE
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, scheduleStart, scheduleEnd, bookingId); err != nil {
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

	query := "UPDATE Booking SET status = 'done', updatedAt = CURRENT_TIMESTAMP() WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, bookingId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlBookingRepository) Reschedule(bookingId string, scheduleStart, scheduleEnd time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Booking
        SET
            scheduleStart = ?, scheduleEnd = ?
        WHERE
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, scheduleStart, scheduleEnd, bookingId); err != nil {
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

func (s *MysqlBookingRepository) getBookingExtras(bookingId string) ([]*models.BookingExtraModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	extrasQuery := `
        SELECT
            bookingId, extraTitle, extraDescription, price
        FROM 
            BookingExtra
        WHERE
            bookingId = ?
    `
	extras := make([]*models.BookingExtraModel, 0)
	if err := s.db.SelectContext(ctx, &extras, extrasQuery, bookingId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return extras, nil
}

// Deprecated: service info is now copied to the booking for snapshot
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
            FORMAT(price, 2) AS price,
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

func (s *MysqlBookingRepository) getClient(bookingId string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            u.id, u.name, u.email, u.imageUrl
        FROM
            User u
            JOIN Booking b ON b.clientId = u.id
        WHERE
            b.id = ?
    `

	user := new(models.UserModel)
	if err := s.db.GetContext(ctx, user, query, bookingId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (s *MysqlBookingRepository) getVendor(bookingId string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            u.id, u.name, u.email, u.imageUrl
        FROM
            User u
            JOIN Booking b ON b.vendorId = u.id
        WHERE
            b.id = ?
    `

	user := new(models.UserModel)
	if err := s.db.GetContext(ctx, user, query, bookingId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}
