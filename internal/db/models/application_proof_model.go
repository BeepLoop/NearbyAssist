package models

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"nearbyassist/internal/db"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ApplicationProofModel struct {
	Model
	UpdateableModel
	ApplicationId int    `json:"applicationId" db:"applicationId"`
	ApplicantId   int    `json:"applicantId" db:"applicantId"`
	Url           string `json:"url" db:"url"`
}

func NewApplicationProofModel(applicationId, applicantId int, db *db.DB) *ApplicationProofModel {
	return &ApplicationProofModel{
		Model:         Model{Db: db},
		ApplicationId: applicationId,
		ApplicantId:   applicantId,
	}
}

func (a *ApplicationProofModel) Create() (int, error) {
	return 0, nil
}

func (a *ApplicationProofModel) Update(id int) error {
	return nil
}

func (a *ApplicationProofModel) Delete(id int) error {
	return nil
}

func (a *ApplicationProofModel) SaveToDisk(uuid string, file *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mimeType := strings.Split(file.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%s.%s", uuid, mimeType)

	storageDir := a.Disk.ApplicationProofLocation
	path := filepath.Join(storageDir, filename)

	dist, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dist.Close()

	_, err = io.Copy(dist, src)
	if err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return filename, nil

}

func (a *ApplicationProofModel) SaveToDb(filename string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.Url = "/resource/application/" + filename

	query := `
        INSERT INTO
            ApplicationProof (applicationId, applicantId, url)
        VALUES
            (:applicationId, :applicantId, :url)
    `

	res, err := a.Db.Conn.NamedExec(query, a)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}
