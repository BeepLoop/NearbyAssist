package db

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/types"
)

type Database interface {
	// Session Queries
	FindSessionByToken(token string) (*models.SessionModel, error)
	FindActiveSessionByToken(token string) (*models.SessionModel, error)
	NewSession(session *models.SessionModel) (int, error)
	LogoutSession(sesionId int) error
	BlacklistToken(token string) error
	FindBlacklistedToken(token string) (*models.BlacklistModel, error)

	// Admin Queries
	FindAdminByUsername(username string) (*models.AdminModel, error)
	NewAdmin(admin *models.AdminModel) (int, error)

	// User Queries
	CountUser() (int, error)
	FindUserById(id int) (*models.UserModel, error)
	FindUserByEmail(email string) (*models.UserModel, error)
	NewUser(user *models.UserModel) (int, error)

	// Vendor Queries
	CountVendor(filter models.VendorStatus) (int, error)
	FindVendorById(id int) (*models.VendorModel, error)
	FindVendorByService(id int) (*response.ServiceVendorDetails, error)
	RestrictVendor(id int) error
	UnrestrictVendor(id int) error

	// Category Queries
	FindAllCategory() ([]models.CategoryModel, error)

	//  Service Queries
	FindServiceById(id int) (*response.ServiceDetails, error)
	FindServiceByVendor(id int) ([]*models.ServiceModel, error)
	FindAllService() ([]*models.ServiceModel, error)
	RegisterService(service *request.NewService) (int, error)
	UpdateService(service *request.UpdateService) error
	DeleteService(id int) error
	GeoSpatialSearch(params *types.SearchParams) ([]*models.ServiceModel, error)

	// Complaint Queries
	CountComplaint() (int, error)
	FileComplaint(complaint *models.ComplaintModel) (int, error)

	// Transaction Queries
	CountTransaction(status models.TransactionStatus) (int, error)
	CreateTransaction(transaction *request.NewTransaction) (int, error)
	FindAllOngoingTransaction(id int, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error)
	GetTransactionHistory(id int, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error)

	// Application Queries
	CountApplication(status models.ApplicationStatus) (int, error)
	CreateApplication(application *request.NewApplication) (int, error)
	FindApplicationById(id int) (*models.ApplicationModel, error)
	FindAllApplication(status models.ApplicationStatus) ([]models.ApplicationModel, error)
	ApproveApplication(id int) error
	RejectApplication(id int) error

	// Review Queries
	CreateReview(review *models.ReviewModel) (int, error)
	FindReviewById(id int) (*models.ReviewModel, error)
	FindAllReviewByService(id int) ([]models.ReviewModel, error)

	// Message Queries
	GetMessages(senderId, receiverId int) ([]models.MessageModel, error)
	GetAllUserConversations(userId int) ([]models.UserModel, error)
	NewMessage(message models.MessageModel) (int, error)

	// Service Photo Queries
	NewServicePhoto(data *models.ServicePhotoModel) (int, error)
	FindAllPhotosByServiceId(serviceId int) ([]response.ServiceImages, error)

	// Application Proof Queries
	NewApplicationProof(data *models.ApplicationProofModel) (int, error)
}
