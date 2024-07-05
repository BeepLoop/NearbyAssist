package db

import (
	"nearbyassist/internal/config"
	"nearbyassist/internal/db/mysql"
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
	LogoutSession(sessionId int) error
	BlacklistToken(token string) error
	FindBlacklistedToken(token string) (*models.BlacklistModel, error)

	// Admin Queries
	FindAdminByUsernameHash(hash string) (*models.AdminModel, error)
	FindAdminById(id int) (*models.AdminModel, error)
	NewAdmin(admin *models.AdminModel) (int, error)
	NewStaff(staff *models.AdminModel) (int, error)

	// User Queries
	CountUser() (int, error)
	FindUserById(id int) (*models.UserModel, error)
	FindUserByEmailHash(hash string) (*models.UserModel, error)
	NewUser(user *models.UserModel) (int, error)

	// Vendor Queries
	CountVendor(filter models.VendorStatus) (int, error)
	FindVendorById(id int) (*models.VendorModel, error)
	FindVendorByService(id int) (*response.ServiceVendorDetails, error)
	RestrictVendor(id int) error
	UnrestrictVendor(id int) error

	// Tag Queries
	FindAllTags() ([]models.TagModel, error)
	FindAllTagByServiceId(serviceId int) ([]string, error)

	//  Service Queries
	FindServiceById(id int) (*response.ServiceDetails, error)
	FindServiceByVendor(id int) ([]*models.ServiceModel, error)
	FindAllService() ([]*models.ServiceModel, error)
	RegisterService(service *request.NewService) (int, error)
	UpdateService(service *request.UpdateService) error
	DeleteService(id int) error
	GeoSpatialSearch(params *types.SearchParams) ([]*models.ServiceSearchResult, error)
	FindServiceOwner(id int) (*response.ServiceOwner, error)
	CountServices() (int, error)

	// Complaint Queries
	CountSystemComplaint() (int, error)
	FindAllSystemComplaints() ([]*response.SystemComplaint, error)
	FindSystemComplaintById(id int) (*models.SystemComplaintModel, error)
	FileVendorComplaint(complaint *request.NewComplaint) (int, error)
	FileSystemComplaint(complaint *request.SystemComplaint) (int, error)
	NewSystemComplaintImage(model *models.SystemComplaintImageModel) (int, error)
	FindSystemComplaintImagesByComplaintId(id int) ([]models.SystemComplaintImageModel, error)

	// Transaction Queries
	CountTransaction(status models.TransactionStatus) (int, error)
	CreateTransaction(transaction *request.NewTransaction) (int, error)
	CompleteTransaction(id int) error
	FindAllOngoingTransaction(id int, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error)
	FindTransactionById(id int) (*models.TransactionModel, error)
	GetTransactionHistory(id int, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error)

	// Application Queries
	CountApplication(status models.ApplicationStatus) (int, error)
	CreateApplication(application *request.NewApplication) (int, error)
	FindApplicationById(id int) (*models.ApplicationModel, error)
	FindAllApplication(status models.ApplicationStatus) ([]models.ApplicationModel, error)
	ApproveApplication(id int) error
	RejectApplication(id int) error

	// Review Queries
	CreateReview(review *request.NewReview) (int, error)
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

	// Verification Queries
	NewIdentityVerification(model *models.IdentityVerificationModel) (int, error)
	NewFrontId(model *models.FrontIdModel) (int, error)
	NewBackId(model *models.BackIdModel) (int, error)
	NewFace(model *models.FaceModel) (int, error)
}

func NewDatabase(conf *config.Config) Database {
	switch conf.DatabaseType {

	case config.DATABASE_MYSQL:
		return mysql.NewMysqlDatabase(conf)

	case config.DATABASE_DUMMY:
		return NewDummyDatabase()

	default:
		panic("Invalid environment. Cannot initialize storage.")
	}
}
