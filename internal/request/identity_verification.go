package request

type IdentityVerificationPayload struct {
	UserId          string `json:"userId" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Address         string `json:"address" validate:"required"`
	IdType          string `json:"idType" validate:"required"`
	IdNumber        string `json:"idNumber" validate:"required"`
	FrontIdImageUrl string `json:"frontIdImageUrl" validate:"required"`
	BackIdImageUrl  string `json:"backIdImageUrl" validate:"required"`
	FaceImageUrl    string `json:"faceImageUrl" validate:"required"`
}
