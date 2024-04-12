package models

import (
	"errors"
	"strconv"
	"strings"
)

type MessageModel struct {
	Model
	Sender   int    `json:"sender" db:"sender"`
	Receiver int    `json:"receiver" db:"receiver"`
	Content  string `json:"content" db:"content"`
}

func NewMessageModel() *MessageModel {
	return &MessageModel{}
}

func MessageValueMapFactory(queryParam string) (map[string]int, error) {
	queries := strings.Split(queryParam, "&")
	if len(queries) != 2 {
		return nil, errors.New("missing required field")
	}

	queryValues := make(map[string]int)
	for _, query := range queries {
		pair := strings.Split(query, "=")
		value, err := strconv.Atoi(pair[1])
		if err != nil {
			return nil, err
		}

		queryValues[pair[0]] = value
	}

	if _, ok := queryValues["sender"]; ok == false {
		return nil, errors.New("missing required field")
	}

	if _, ok := queryValues["receiver"]; ok == false {
		return nil, errors.New("missing required field")
	}

	return queryValues, nil
}
