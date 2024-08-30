package mysql

// import (
// 	"context"
// 	"nearbyassist/internal/models"
// 	"nearbyassist/internal/response"
// 	"time"
// )
//
// func (m *Mysql) NewServicePhoto(data *models.ServicePhotoModel) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	query := `
//         INSERT INTO
//             ServicePhoto (id, vendorId, serviceId, url)
//         VALUES
//             (:id, :vendorId, :serviceId, :url)
//     `
//
// 	if _, err := m.Conn.NamedExecContext(ctx, query, data); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return data.Id, nil
// }
//
// func (m *Mysql) FindAllPhotosByServiceId(serviceId string) ([]response.ServiceImages, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	query := `
//         SELECT 
//             id AS imageId,
//             url AS imageUrl
//         FROM
//             ServicePhoto
//         WHERE
//             serviceId = ?
//     `
//
// 	images := make([]response.ServiceImages, 0)
// 	if err := m.Conn.SelectContext(ctx, &images, query, serviceId); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return images, nil
// }
