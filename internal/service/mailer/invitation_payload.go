package mailer

import "fmt"

type invitationPayload struct {
	joinUrl    string
	inviteCode string
}

func NewInvitationPayload(joinUrl, code string) *invitationPayload {
	return &invitationPayload{
		joinUrl:    joinUrl,
		inviteCode: code,
	}
}

func (p *invitationPayload) GetContent() string {
	return fmt.Sprintf("%s?code=%s", p.joinUrl, p.inviteCode)
}
