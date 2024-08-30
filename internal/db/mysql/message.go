package mysql

// import (
// 	"context"
// 	"nearbyassist/internal/models"
// 	"time"
// )

// func (m *Mysql) NewMessage(message models.MessageModel) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         INSERT INTO
//             Message (id, sender, receiver, content)
//         VALUES
//             (:id, :sender, :receiver, :content)
//     `
//
// 	if _, err := m.Conn.NamedExecContext(ctx, query, message); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return message.Id, nil
// }
//
// func (m *Mysql) GetMessages(senderId, receiverId string) ([]models.MessageModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         SELECT
//             id, sender, receiver, content, createdAt
//         FROM
//             Message
//         WHERE
//             sender = ? AND receiver = ?
//         OR
//             sender = ? AND receiver = ?
//         ORDER BY
//             createdAt
//     `
//
// 	messages := make([]models.MessageModel, 0)
// 	if err := m.Conn.SelectContext(
// 		ctx,
// 		&messages,
// 		query,
// 		senderId,
// 		receiverId,
// 		receiverId,
// 		senderId,
// 	); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return messages, nil
// }
//
// func (m *Mysql) GetAllUserConversations(userId string) ([]*models.UserModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         SELECT DISTINCT
//             u.id,
//             u.name,
//             u.imageUrl
//         FROM
//             User u
//         JOIN
//             Message m ON u.id = m.sender OR u.id = m.receiver
//         WHERE
//             u.id <> ?
//     `
//
// 	conversations := make([]*models.UserModel, 0)
// 	if err := m.Conn.SelectContext(ctx, &conversations, query, userId); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return conversations, nil
// }
