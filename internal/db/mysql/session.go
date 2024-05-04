package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) FindActiveSessionByToken(token string) (*models.SessionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, token, status FROM Session WHERE token = ? AND status = 'online'"

	session := new(models.SessionModel)
	err := m.Conn.GetContext(ctx, session, query, token)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return session, nil
}

func (m *Mysql) FindSessionByToken(token string) (*models.SessionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, token, status FROM Session WHERE token = ?"

	session := new(models.SessionModel)
	err := m.Conn.GetContext(ctx, session, query, token)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return session, nil
}

func (m *Mysql) NewSession(session *models.SessionModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Session (token) VALUES (:token)"

	res, err := m.Conn.NamedExecContext(ctx, query, session)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) LogoutSession(sessionId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE Session SET status = 'offline' WHERE id = ?"

	_, err := m.Conn.ExecContext(ctx, query, sessionId)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (m *Mysql) FindBlacklistedToken(token string) (*models.BlacklistModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, token FROM Blacklist WHERE token = ?"

	blacklist := new(models.BlacklistModel)
	err := m.Conn.GetContext(ctx, blacklist, query, token)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return blacklist, nil
}

func (m *Mysql) BlacklistToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Blacklist (token) VALUES (?)"

	_, err := m.Conn.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
