package mysql

// import (
// 	"context"
// 	"nearbyassist/internal/models"
// 	"nearbyassist/internal/request"
// 	"nearbyassist/internal/response"
// 	"time"
// )

// func (m *Mysql) CountSystemComplaint() (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT COUNT(*) FROM SystemComplaint"
//
// 	count := 0
// 	err := m.Conn.GetContext(ctx, &count, query)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return 0, context.DeadlineExceeded
// 	}
//
// 	return count, nil
// }

// func (m *Mysql) FindAllSystemComplaints() ([]*response.SystemComplaint, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, title FROM SystemComplaint"
//
// 	complaints := make([]*response.SystemComplaint, 0)
// 	if err := m.Conn.SelectContext(ctx, &complaints, query); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return complaints, nil
// }

// func (m *Mysql) FindSystemComplaintById(systemComplaintId string) (*models.SystemComplaintModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT * FROM SystemComplaint WHERE id = ?"
//
// 	complaint := &models.SystemComplaintModel{}
// 	if err := m.Conn.GetContext(ctx, complaint, query, systemComplaintId); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return complaint, nil
// }

// func (m *Mysql) FileVendorComplaint(complaint *request.NewComplaint) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         INSERT INTO 
//             Complaint (id, vendorId, code, title, content)
//         VALUES
//             (:id, :vendorId, :code, :title, :content)
//     `
//
// 	if _, err := m.Conn.NamedExecContext(ctx, query, complaint); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return complaint.Id, nil
// }

// func (m *Mysql) FileSystemComplaint(complaint *request.SystemComplaint) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         INSERT INTO 
//             SystemComplaint (id, title, detail)
//         VALUES
//             (:id, :title, :detail)
//     `
//
// 	if _, err := m.Conn.NamedExecContext(ctx, query, complaint); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return complaint.Id, nil
// }

// func (m *Mysql) NewSystemComplaintImage(model *models.SystemComplaintImageModel) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         INSERT INTO 
//             SystemComplaintImage (id, complaintId, url)
//         VALUES
//             (:id, :complaintId, :url)
//     `
//
// 	if _, err := m.Conn.NamedExecContext(ctx, query, model); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return model.Id, nil
// }

// func (m *Mysql) FindSystemComplaintImagesByComplaintId(systemComplaintId string) ([]models.SystemComplaintImageModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT * FROM SystemComplaintImage WHERE complaintId = ?"
//
// 	images := make([]models.SystemComplaintImageModel, 0)
// 	if err := m.Conn.SelectContext(ctx, &images, query, systemComplaintId); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return images, nil
// }
