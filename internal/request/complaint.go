package request

type BugReportPayload struct {
	Title  string `json:"title" validate:"required"`
	Detail string `json:"detail" validate:"required"`
}
