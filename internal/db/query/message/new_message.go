package message_query

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
)

func NewMessage(message types.Message) (*types.Message, error) {
	insertQuery := `
        INSERT INTO
            Message (sender, receiver, content)
        VALUES
            (:sender, :receiver, :content)
    `

	res, err := db.Connection.NamedExec(insertQuery, message)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	retrieveQuery := `
        SELECT
            id, sender, receiver, content 
        FROM
            Message 
        WHERE
            id = ?
    `
	inserted := new(types.Message)
	err = db.Connection.Get(inserted, retrieveQuery, id)
	if err != nil {
		return nil, err
	}

	return inserted, nil
}
