package db

// import (
// 	"nearbyassist/internal/config"
// 	"nearbyassist/internal/db/mysql"
// 	"nearbyassist/internal/models"
// 	"nearbyassist/internal/request"
// 	"nearbyassist/internal/response"
// 	"nearbyassist/internal/types"
// )
//
// type Database interface {
// 	// Session Queries
// 	FindSessionByToken(token string) (*models.SessionModel, error)
// 	FindActiveSessionByToken(token string) (*models.SessionModel, error)
// 	NewSession(session *models.SessionModel) (string, error)
// 	LogoutSession(sessionId string) error
// 	BlacklistToken(token string) error
// 	FindBlacklistedToken(token string) (*models.BlacklistModel, error)
//
// 	// Admin Queries
// 	FindAdminByUsernameHash(hash string) (*models.AdminModel, error)
// 	FindAdminById(id string) (*models.AdminModel, error)
// 	NewAdmin(admin *models.AdminModel) (string, error)
// 	NewStaff(staff *models.AdminModel) (string, error)
//
// 	// User Queries
// 	CountUser() (int, error)
// 	CheckUserVerification(userId string) (bool, error)
// 	FindUserById(userId string) (*models.UserModel, error)
// 	FindUserByEmailHash(hash string) (*models.UserModel, error)
// 	NewUser(user *models.UserModel) (string, error)
//
// 	// Vendor Queries
// 	CountVendor(filter models.VendorStatus) (int, error)
// 	FindVendorById(vendorId string) (*models.VendorModel, error)
// 	FindVendorByService(serviceId string) (*response.ServiceVendorDetails, error)
// 	RestrictVendor(vendorId string) error
// 	UnrestrictVendor(vendorId string) error
//
// 	// Tag Queries
// 	FindAllTags() ([]models.TagModel, error)
// 	FindAllTagByServiceId(serviceTagId string) ([]string, error)
//
// 	//  Service Queries
// 	FindServiceById(serviceId string) (*response.ServiceDetails, error)
// 	FindServiceByVendor(vendorId string) ([]*models.ServiceModel, error)
// 	FindAllService() ([]*models.ServiceModel, error)
// 	RegisterService(service *request.NewService) (string, error)
// 	UpdateService(service *request.UpdateService) error
// 	DeleteService(serviceId string) error
// 	GeoSpatialSearch(params *types.SearchParams) ([]*models.ServiceSearchResult, error)
// 	FindServiceOwner(serviceId string) (*response.ServiceOwner, error)
// 	CountServices() (int, error)
//
// 	// Complaint Queries
// 	CountSystemComplaint() (int, error)
// 	FindAllSystemComplaints() ([]*response.SystemComplaint, error)
// 	FindSystemComplaintById(systemComplaintId string) (*models.SystemComplaintModel, error)
// 	FileVendorComplaint(complaint *request.NewComplaint) (string, error)
// 	FileSystemComplaint(complaint *request.SystemComplaint) (string, error)
// 	NewSystemComplaintImage(model *models.SystemComplaintImageModel) (string, error)
// 	FindSystemComplaintImagesByComplaintId(systemComplaintId string) ([]models.SystemComplaintImageModel, error)
//
// 	// Transaction Queries
// 	CountTransaction(status models.TransactionStatus) (int, error)
// 	CreateTransaction(transaction *request.NewTransaction) (string, error)
// 	CompleteTransaction(transactionId string) error
// 	FindAllOngoingTransaction(id string, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error)
// 	FindUserTransactions(userId string) ([]*models.DetailedTransactionModel, error)
// 	FindTransactionById(transactionId string) (*models.TransactionModel, error)
// 	GetTransactionHistory(id string, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error)
//
// 	// Application Queries
// 	CountApplication(status models.ApplicationStatus) (int, error)
// 	CreateApplication(application *request.NewApplication) (string, error)
// 	FindApplicationById(applicationId string) (*models.ApplicationModel, error)
// 	FindAllApplication(status models.ApplicationStatus) ([]response.Application, error)
// 	ApproveApplication(applicationId string) error
// 	RejectApplication(applicationId string) error
//
// 	// Review Queries
// 	CreateReview(review *request.NewReview) (string, error)
// 	FindReviewById(reviewId string) (*models.ReviewModel, error)
// 	FindAllReviewByService(serviceId string) ([]models.ReviewModel, error)
//
// 	// Message Queries
// 	GetMessages(senderId, receiverId string) ([]models.MessageModel, error)
// 	GetAllUserConversations(userId string) ([]*models.UserModel, error)
// 	NewMessage(message models.MessageModel) (string, error)
//
// 	// Service Photo Queries
// 	NewServicePhoto(data *models.ServicePhotoModel) (string, error)
// 	FindAllPhotosByServiceId(serviceId string) ([]response.ServiceImages, error)
//
// 	// Application Proof Queries
// 	NewApplicationProof(data *models.ApplicationProofModel) (string, error)
//
// 	// Verification Queries
// 	FindAllIdentityVerification() ([]response.AllVerification, error)
// 	NewIdentityVerification(model *models.IdentityVerificationModel) (string, error)
// 	FindIdentityVerificationById(identityVerificationId string) (*models.IdentityVerificationModel, error)
// 	NewFrontId(model *models.FrontIdModel) (string, error)
// 	NewBackId(model *models.BackIdModel) (string, error)
// 	NewFace(model *models.FaceModel) (string, error)
// }
//
// func NewDatabase(conf *config.Config) Database {
// 	switch conf.DatabaseType {
//
// 	case config.DATABASE_MYSQL:
// 		return mysql.NewMysqlDatabase(conf)
//
// 	case config.DATABASE_DUMMY:
// 		return NewDummyDatabase()
//
// 	default:
// 		panic("Invalid environment. Cannot initialize storage.")
// 	}
// }
