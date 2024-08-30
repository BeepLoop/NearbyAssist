package db

// import (
// 	"nearbyassist/internal/models"
// 	"nearbyassist/internal/request"
// 	"nearbyassist/internal/response"
// 	"nearbyassist/internal/types"
// )
//
// type DummyDatabase struct{}
//
// func NewDummyDatabase() *DummyDatabase {
// 	return &DummyDatabase{}
// }
//
// func (d *DummyDatabase) FindSessionByToken(token string) (*models.SessionModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindActiveSessionByToken(token string) (*models.SessionModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) NewSession(session *models.SessionModel) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) LogoutSession(sessionId string) error {
// 	return nil
// }
//
// func (d *DummyDatabase) BlacklistToken(token string) error {
// 	return nil
// }
//
// func (d *DummyDatabase) FindBlacklistedToken(token string) (*models.BlacklistModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindAdminByUsernameHash(hash string) (*models.AdminModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindAdminById(id string) (*models.AdminModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) NewAdmin(admin *models.AdminModel) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) NewStaff(staff *models.AdminModel) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) CheckUserVerification(id string) (bool, error) {
// 	return false, nil
// }
//
// func (d *DummyDatabase) CountUser() (int, error) {
// 	return 0, nil
// }
//
// func (d *DummyDatabase) FindUserById(id string) (*models.UserModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindUserByEmailHash(hash string) (*models.UserModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) NewUser(user *models.UserModel) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) CountVendor(filter models.VendorStatus) (int, error) {
// 	return 0, nil
// }
//
// func (d *DummyDatabase) FindVendorById(id string) (*models.VendorModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindVendorByService(id string) (*response.ServiceVendorDetails, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) RestrictVendor(id string) error {
// 	return nil
// }
//
// func (d *DummyDatabase) UnrestrictVendor(id string) error {
// 	return nil
// }
//
// func (d *DummyDatabase) FindAllTags() ([]models.TagModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindAllTagByServiceId(serviceId string) ([]string, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) CountServices() (int, error) {
// 	return 0, nil
// }
//
// func (d *DummyDatabase) FindServiceById(id string) (*response.ServiceDetails, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindServiceByVendor(id string) ([]*models.ServiceModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindAllService() ([]*models.ServiceModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) RegisterService(service *request.NewService) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) UpdateService(service *request.UpdateService) error {
// 	return nil
// }
//
// func (d *DummyDatabase) DeleteService(id string) error {
// 	return nil
// }
//
// func (d *DummyDatabase) GeoSpatialSearch(params *types.SearchParams) ([]*models.ServiceSearchResult, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindServiceOwner(id string) (*response.ServiceOwner, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) CountSystemComplaint() (int, error) {
// 	return 0, nil
// }
//
// func (m *DummyDatabase) FindAllSystemComplaints() ([]*response.SystemComplaint, error) {
// 	return nil, nil
// }
//
// func (m *DummyDatabase) FindSystemComplaintById(id string) (*models.SystemComplaintModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FileVendorComplaint(complaint *request.NewComplaint) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) FileSystemComplaint(complaint *request.SystemComplaint) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) NewSystemComplaintImage(model *models.SystemComplaintImageModel) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) FindSystemComplaintImagesByComplaintId(id string) ([]models.SystemComplaintImageModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) CountTransaction(status models.TransactionStatusFilter) (int, error) {
// 	return 0, nil
// }
//
// func (d *DummyDatabase) CreateTransaction(transaction *request.NewTransaction) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) CompleteTransaction(id string) error {
// 	return nil
// }
//
// func (d *DummyDatabase) FindAllOngoingTransaction(id string, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindUserTransactions(id string) ([]*models.DetailedTransactionModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindTransactionById(id string) (*models.TransactionModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) GetTransactionHistory(id string, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) CountApplication(status models.ApplicationStatusFilter) (int, error) {
// 	return 0, nil
// }
//
// func (d *DummyDatabase) CreateApplication(application *request.NewApplication) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) FindApplicationById(id string) (*models.ApplicationModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindAllApplication(status models.ApplicationStatusFilter) ([]response.Application, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) ApproveApplication(id string) error {
// 	return nil
// }
//
// func (d *DummyDatabase) RejectApplication(id string) error {
// 	return nil
// }
//
// func (d *DummyDatabase) CreateReview(review *request.NewReview) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) FindReviewById(id string) (*models.ReviewModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) FindAllReviewByService(id string) ([]models.ReviewModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) GetMessages(senderId, receiverId string) ([]models.MessageModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) GetAllUserConversations(userId string) ([]*models.UserModel, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) NewMessage(message models.MessageModel) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) NewServicePhoto(data *models.ServicePhotoModel) (string, error) {
// 	return "", nil
// }
//
// func (d *DummyDatabase) FindAllPhotosByServiceId(serviceId string) ([]response.ServiceImages, error) {
// 	return nil, nil
// }
//
// func (d *DummyDatabase) NewApplicationProof(data *models.ApplicationProofModel) (string, error) {
// 	return "", nil
// }
//
// func (m *DummyDatabase) FindAllIdentityVerification() ([]response.AllVerification, error) {
// 	return nil, nil
// }
//
// func (m *DummyDatabase) NewIdentityVerification(model *models.IdentityVerificationModel) (string, error) {
// 	return "", nil
// }
//
// func (m *DummyDatabase) FindIdentityVerificationById(id string) (*models.IdentityVerificationModel, error) {
// 	return nil, nil
// }
//
// func (m *DummyDatabase) NewFrontId(model *models.FrontIdModel) (string, error) {
// 	return "", nil
// }
//
// func (m *DummyDatabase) NewBackId(model *models.BackIdModel) (string, error) {
// 	return "", nil
// }
//
// func (m *DummyDatabase) NewFace(model *models.FaceModel) (string, error) {
//     return "", nil
// }
