package models

type SessionModel struct {
	Model
	UpdateableModel
	Status       string `json:"status" db:"status"`
	RefreshToken string `json:"refreshToken" db:"refreshToken"`
}

func NewSessionModel(refreshToken string) *SessionModel {
	return &SessionModel{
		RefreshToken: refreshToken,
	}
}

// func (s *SessionModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "INSERT INTO Session (id, refreshToken) VALUES (:id, :refreshToken)"
// 	if _, err := s.Conn.NamedExecContext(ctx, query, s); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return s.Id, nil
// }
//
// func (s *SessionModel) Logout() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "UPDATE Session SET status = 'offline' WHERE id = ?"
// 	if _, err := s.Conn.ExecContext(ctx, query, s.Id); err != nil {
// 		return err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return context.DeadlineExceeded
// 	}
//
// 	return nil
// }
//
// func (s *SessionModel) FindByToken() (*SessionModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, refreshToken, status FROM Session WHERE refreshToken = ?"
//
// 	if err := s.Conn.GetContext(ctx, s, query, s.RefreshToken); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return s, nil
// }
//
// func (s *SessionModel) GetIfActive() (*SessionModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, refreshToken, status FROM Session WHERE refreshToken = ? AND status = 'online'"
//
// 	err := s.Conn.GetContext(ctx, s, query, s.RefreshToken)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return s, nil
// }
