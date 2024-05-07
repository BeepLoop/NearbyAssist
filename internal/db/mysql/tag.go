package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) FindAllTags() ([]models.TagModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, title FROM Tag"

	tags := make([]models.TagModel, 0)
	err := m.Conn.SelectContext(ctx, &tags, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return tags, nil
}

func (m *Mysql) FindAllTagByServiceId(serviceId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.title AS tag
        FROM
            Service_Tag st
            JOIN Tag t ON t.id = st.tagId
        WHERE
            st.serviceId = ?;
    `

	tags := make([]string, 0)
	if err := m.Conn.SelectContext(ctx, &tags, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return tags, nil
}
