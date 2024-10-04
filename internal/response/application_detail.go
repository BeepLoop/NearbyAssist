package response

type ApplicationDetail struct {
	Id          string `db:"id" json:"id"`
	ApplicantId string `db:"applicantId" json:"applicantId"`
	Job         string `db:"job" json:"job"`
	Status      string `db:"status" json:"status"`
	CreatedAt   string `db:"createdAt" json:"createdAt"`
	ProofUrl    string `db:"proofUrl" json:"proofUrl"`
}
