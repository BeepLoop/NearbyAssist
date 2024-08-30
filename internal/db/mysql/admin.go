package mysql

// import (
// 	"context"
// 	"nearbyassist/internal/models"
// 	"time"
// )
//
// func (m *Mysql) FindAdminByUsernameHash(hash string) (*models.AdminModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, username, password, role FROM Admin WHERE usernameHash = ?"
//
// 	admin := &models.AdminModel{}
// 	err := m.Conn.GetContext(ctx, admin, query, hash)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return admin, nil
// }
//
// func (m *Mysql) FindAdminById(id string) (*models.AdminModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, username, password, role FROM Admin WHERE id = ?"
//
// 	admin := &models.AdminModel{}
// 	err := m.Conn.GetContext(ctx, admin, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return admin, nil
// }
//
// func (m *Mysql) NewAdmin(admin *models.AdminModel) (string, error) {
// 	// TODO: implement this method
// 	return "", nil
// }
//
// func (m *Mysql) NewStaff(staff *models.AdminModel) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "INSERT INTO Admin (id, username, password, usernameHash) VALUES (:id, :username, :password, :usernameHash)"
// 	if _, err := m.Conn.NamedExecContext(ctx, query, staff); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return staff.Id, nil
// }
