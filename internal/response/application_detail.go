package response

type ApplicationDetail struct {
	Id              string `db:"id" json:"id"`
	ApplicantId     string `db:"applicantId" json:"applicantId"`
	ExpertiseId     string `db:"expertiseId" json:"expertiseId"`
	Status          string `db:"status" json:"status"`
	CreatedAt       string `db:"createdAt" json:"createdAt"`
	ProofUrl        string `db:"proofUrl" json:"proofUrl"`
	PoliceClearance string `db:"policeClearance" json:"policeClearance"`
}
