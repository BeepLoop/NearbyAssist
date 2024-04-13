package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) FindSessionByToken(token string) (*models.SessionModel, error) {
	return nil, nil
}

func (m *Mysql) FindSessionById(id int) (*models.SessionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, userId, token, status FROM Session WHERE id = ?"

	session := new(models.SessionModel)
	err := m.Conn.GetContext(ctx, session, query, id)
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

	query := "INSERT INTO Session (userId, token) VALUES (:userId, :token)"

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
