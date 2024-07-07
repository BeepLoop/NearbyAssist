package db

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/types"
)

type DummyDatabase struct{}

func NewDummyDatabase() *DummyDatabase {
	return &DummyDatabase{}
}

func (d *DummyDatabase) FindSessionByToken(token string) (*models.SessionModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindActiveSessionByToken(token string) (*models.SessionModel, error) {
	return nil, nil
}

func (d *DummyDatabase) NewSession(session *models.SessionModel) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) LogoutSession(sessionId int) error {
	return nil
}

func (d *DummyDatabase) BlacklistToken(token string) error {
	return nil
}

func (d *DummyDatabase) FindBlacklistedToken(token string) (*models.BlacklistModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindAdminByUsernameHash(hash string) (*models.AdminModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindAdminById(id int) (*models.AdminModel, error) {
	return nil, nil
}

func (d *DummyDatabase) NewAdmin(admin *models.AdminModel) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) NewStaff(staff *models.AdminModel) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) CountUser() (int, error) {
	return 0, nil
}

func (d *DummyDatabase) FindUserById(id int) (*models.UserModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindUserByEmailHash(hash string) (*models.UserModel, error) {
	return nil, nil
}

func (d *DummyDatabase) NewUser(user *models.UserModel) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) CountVendor(filter models.VendorStatus) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) FindVendorById(id int) (*models.VendorModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindVendorByService(id int) (*response.ServiceVendorDetails, error) {
	return nil, nil
}

func (d *DummyDatabase) RestrictVendor(id int) error {
	return nil
}

func (d *DummyDatabase) UnrestrictVendor(id int) error {
	return nil
}

func (d *DummyDatabase) FindAllTags() ([]models.TagModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindAllTagByServiceId(serviceId int) ([]string, error) {
	return nil, nil
}

func (d *DummyDatabase) CountServices() (int, error) {
	return 0, nil
}

func (d *DummyDatabase) FindServiceById(id int) (*response.ServiceDetails, error) {
	return nil, nil
}

func (d *DummyDatabase) FindServiceByVendor(id int) ([]*models.ServiceModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindAllService() ([]*models.ServiceModel, error) {
	return nil, nil
}

func (d *DummyDatabase) RegisterService(service *request.NewService) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) UpdateService(service *request.UpdateService) error {
	return nil
}

func (d *DummyDatabase) DeleteService(id int) error {
	return nil
}

func (d *DummyDatabase) GeoSpatialSearch(params *types.SearchParams) ([]*models.ServiceSearchResult, error) {
	return nil, nil
}

func (d *DummyDatabase) FindServiceOwner(id int) (*response.ServiceOwner, error) {
	return nil, nil
}

func (d *DummyDatabase) CountSystemComplaint() (int, error) {
	return 0, nil
}

func (m *DummyDatabase) FindAllSystemComplaints() ([]*response.SystemComplaint, error) {
	return nil, nil
}

func (m *DummyDatabase) FindSystemComplaintById(id int) (*models.SystemComplaintModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FileVendorComplaint(complaint *request.NewComplaint) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) FileSystemComplaint(complaint *request.SystemComplaint) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) NewSystemComplaintImage(model *models.SystemComplaintImageModel) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) FindSystemComplaintImagesByComplaintId(id int) ([]models.SystemComplaintImageModel, error) {
	return nil, nil
}

func (d *DummyDatabase) CountTransaction(status models.TransactionStatus) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) CreateTransaction(transaction *request.NewTransaction) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) CompleteTransaction(id int) error {
	return nil
}

func (d *DummyDatabase) FindAllOngoingTransaction(id int, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindTransactionById(id int) (*models.TransactionModel, error) {
	return nil, nil
}

func (d *DummyDatabase) GetTransactionHistory(id int, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error) {
	return nil, nil
}

func (d *DummyDatabase) CountApplication(status models.ApplicationStatus) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) CreateApplication(application *request.NewApplication) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) FindApplicationById(id int) (*models.ApplicationModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindAllApplication(status models.ApplicationStatus) ([]response.Application, error) {
	return nil, nil
}

func (d *DummyDatabase) ApproveApplication(id int) error {
	return nil
}

func (d *DummyDatabase) RejectApplication(id int) error {
	return nil
}

func (d *DummyDatabase) CreateReview(review *request.NewReview) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) FindReviewById(id int) (*models.ReviewModel, error) {
	return nil, nil
}

func (d *DummyDatabase) FindAllReviewByService(id int) ([]models.ReviewModel, error) {
	return nil, nil
}

func (d *DummyDatabase) GetMessages(senderId, receiverId int) ([]models.MessageModel, error) {
	return nil, nil
}

func (d *DummyDatabase) GetAllUserConversations(userId int) ([]models.UserModel, error) {
	return nil, nil
}

func (d *DummyDatabase) NewMessage(message models.MessageModel) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) NewServicePhoto(data *models.ServicePhotoModel) (int, error) {
	return 0, nil
}

func (d *DummyDatabase) FindAllPhotosByServiceId(serviceId int) ([]response.ServiceImages, error) {
	return nil, nil
}

func (d *DummyDatabase) NewApplicationProof(data *models.ApplicationProofModel) (int, error) {
	return 0, nil
}

func (m *DummyDatabase) FindAllIdentityVerification() ([]response.AllVerification, error) {
	return nil, nil
}

func (m *DummyDatabase) NewIdentityVerification(model *models.IdentityVerificationModel) (int, error) {
	return 0, nil
}

func (m *DummyDatabase) FindIdentityVerificationById(id int) (*models.IdentityVerificationModel, error) {
	return nil, nil
}

func (m *DummyDatabase) NewFrontId(model *models.FrontIdModel) (int, error) {
	return 0, nil
}

func (m *DummyDatabase) NewBackId(model *models.BackIdModel) (int, error) {
	return 0, nil
}

func (m *DummyDatabase) NewFace(model *models.FaceModel) (int, error) {
	return 0, nil
}
