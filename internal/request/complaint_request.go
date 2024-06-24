package request

type NewComplaint struct {
	Code    int    `json:"code" db:"code" validate:"required"`
	Title   string `json:"title" db:"title" validate:"required"`
	Content string `json:"content" db:"content" validate:"required"`
}

type SystemComplaint struct {
	Title  string `json:"title" db:"title" validate:"required"`
	Detail string `json:"detail" db:"detail" validate:"required"`
}
