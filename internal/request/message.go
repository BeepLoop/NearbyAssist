package request

type MarkMessageSeenPayload struct {
	MessageIDs []string `json:"messageIds" validate:"required"`
}
