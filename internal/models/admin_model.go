package models

// import (
// 	"context"
// 	"nearbyassist/internal/id_generator"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// 	"golang.org/x/crypto/bcrypt"
// )

type AdminRole string

const (
	ADMIN_ROLE_ADMIN AdminRole = "admin"
	ADMIN_ROLE_STAFF AdminRole = "staff"
)

type AdminModel struct {
	Model
	UpdateableModel
	Username     string    `json:"username" db:"username"`
	Password     string    `json:"password" db:"password"`
	Role         AdminRole `json:"role" db:"role"`
	UsernameHash string    `json:"usernameHash" db:"usernameHash"`
}

// func NewAdminModelWithId(id string, conn *sqlx.DB) *AdminModel {
// 	return &AdminModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }

// func (a *AdminModel) EncryptUsername(encryptFunc func(string) (string, error)) (*AdminModel, error) {
// 	if encrypted, err := encryptFunc(a.Username); err != nil {
// 		return nil, err
// 	} else {
// 		a.Username = encrypted
// 	}
//
// 	return a, nil
// }
//
// func (a *AdminModel) DecryptUsername(decryptFunc func(string) (string, error)) (*AdminModel, error) {
// 	if decrypted, err := decryptFunc(a.Username); err != nil {
// 		return nil, err
// 	} else {
// 		a.Username = decrypted
// 	}
//
// 	return a, nil
// }
//
// func (a *AdminModel) EncryptPassword() (*AdminModel, error) {
// 	if hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost); err != nil {
// 		return nil, err
// 	} else {
// 		a.Password = string(hashed)
// 	}
// 	return a, nil
// }
//
// func (a *AdminModel) HashUsername(plainUsername string, hashFunc func([]byte) (string, error)) (*AdminModel, error) {
// 	if hash, err := hashFunc([]byte(plainUsername)); err != nil {
// 		return nil, err
// 	} else {
// 		a.UsernameHash = hash
// 	}
// 	return a, nil
// }
//
// func (a *AdminModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "INSERT INTO Admin (id, username, password, usernameHash, role) VALUES (:id, :username, :password, :usernameHash, :role)"
// 	if _, err := a.Conn.NamedExecContext(ctx, query, a); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return a.Id, nil
// }
//
// func (a *AdminModel) FindById(id string) (*AdminModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, username, password, role FROM Admin WHERE id = ?"
//
// 	err := a.Conn.GetContext(ctx, a, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return a, nil
// }
//
// func (a *AdminModel) FindByUsernameHash() (*AdminModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, username, password, role FROM Admin WHERE usernameHash = ?"
//
// 	err := a.Conn.GetContext(ctx, a, query, a.UsernameHash)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return a, nil
// }
//
// func (a *AdminModel) IsPasswordMatch(plainPwd string) bool {
// 	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plainPwd)); err != nil {
// 		return false
// 	}
// 	return true
// }
