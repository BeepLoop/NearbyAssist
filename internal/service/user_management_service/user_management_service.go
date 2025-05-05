package user_management_service

import (
	admin_repo "nearbyassist/internal/repository/admin"
	booking_repo "nearbyassist/internal/repository/booking"
	notification_repo "nearbyassist/internal/repository/notification"
	report_user_repo "nearbyassist/internal/repository/report_user"
	service_repo "nearbyassist/internal/repository/service"
	user_repo "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	"nearbyassist/internal/service/core"
	resource_service "nearbyassist/internal/service/resource"
	"nearbyassist/internal/service/websocket"
)

type Service struct {
	userStore       user_repo.UserRepository
	vendorStore     vendor_repo.VendorRepository
	notifStore      notification_repo.NotificationRepository
	bookingStore    booking_repo.BookingRepository
	serviceStore    service_repo.ServiceRepository
	reportUserStore report_user_repo.ReportUserRepository
	adminStore      admin_repo.AdminRepository
	resourceService *resource_service.Service
	ws              websocket.Socket
	encrypt         core.Encryption
	hash            core.Hash
}

func NewService(
	userStore user_repo.UserRepository,
	vendorStore vendor_repo.VendorRepository,
	notifStore notification_repo.NotificationRepository,
	bookingStore booking_repo.BookingRepository,
	serviceStore service_repo.ServiceRepository,
	reportUserStore report_user_repo.ReportUserRepository,
	adminStore admin_repo.AdminRepository,
	resourceService *resource_service.Service,
	ws websocket.Socket,
	encrypt core.Encryption,
	hash core.Hash,
) *Service {
	return &Service{
		userStore:       userStore,
		vendorStore:     vendorStore,
		notifStore:      notifStore,
		bookingStore:    bookingStore,
		serviceStore:    serviceStore,
		reportUserStore: reportUserStore,
		adminStore:      adminStore,
		resourceService: resourceService,
		ws:              ws,
		encrypt:         encrypt,
		hash:            hash,
	}
}
