package postgres

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ChatRepository struct {
	db *pgxpool.Pool
	logger log.Logger
}

func NewChatRepository(db *pgxpool.Pool, logger log.Logger) *ChatRepository {
	return &ChatRepository{
		db: db,
		logger: logger,
	}
}

func (c *ChatRepository) FindRecentMessages(ctx context.Context, sessionId string , numLastMessages int32) ([]*models.Chat, error){
	// sub is subquery
	// basically gets the first $2 messages, then sort it on descending order
	// then sort it once again on ascending order
	query := `
		SELECT * FROM (
			SELECT * FROM chats 
			WHERE session_room_id = $1 
			ORDER BY created_at DESC 
			LIMIT $2
		) sub
		ORDER BY created_at ASC
	`

	rows, err := c.db.Query(ctx, query, sessionId, numLastMessages)

	if err != nil {
		return nil, fmt.Errorf("chat repository find recent messages error: %w", err)
	}

	defer rows.Close()

	result := []*models.Chat{}

	for rows.Next() {
		var chat models.Chat

		if err := rows.Scan(&chat.ID, &chat.Message, &chat.SessionRoomId, &chat.ReceivedTime, &chat.CreatedAt, &chat.UpdatedAt, &chat.Role); err != nil {
			return nil, fmt.Errorf("chat repository find recent messages error: %w", err)
		}

		result = append(result, &chat)
	}



	return result, nil


}

func (c *ChatRepository) FindByRoomType(ctx context.Context, roomType string) ([]*models.Chat, error) {
	query := `SELECT * FROM chats c INNER JOIN session_rooms sr ON c.session_room_id = sr.id WHERE sr.room_type_id = $1`
	
	rows, err := c.db.Query(ctx, query, roomType)

	if err != nil {
		return nil, fmt.Errorf("chat repository FindByRoomType function: %w", err) 
	}

	defer rows.Close()

	var out []*models.Chat

	for rows.Next(){
		var m models.Chat
		if err := rows.Scan(&m.ID, &m.Message, &m.ReceivedTime, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("chat repository FindByRoomType function scan: %w", err) 
		}
		out = append(out, &m)
	}



	return out, nil  

}

func (c *ChatRepository) Create(ctx context.Context, chat *models.Chat) error {
	c.logger.Debug("Creating chat message", zap.String("message", chat.Message))
	
	query := `INSERT INTO chats (message, session_room_id, role, received_time, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW(), NOW())`
	_, err := c.db.Exec(ctx, query, chat.Message,  chat.SessionRoomId, chat.Role)
	c.logger.Debug("Executed chat insert query", zap.Error(err))
	if err != nil {
		return apperr.Internal(fmt.Errorf("chat repository create function: %w", err)) 
	}

	return nil 

}

func (c *ChatRepository) DeleteAllMessages(ctx context.Context, sessionId string, ignoreFirstMessage bool) error {
	var query string 

	if ignoreFirstMessage {
		query = "DELETE FROM chats WHERE session_room_id = $1 AND id != (SELECT id FROM chats WHERE session_room_id = $1 ORDER BY created_at ASC LIMIT 1)"
	}else {
		query = "DELETE FROM chats WHERE session_room_id = $1"
	}
	
	_, err := c.db.Exec(ctx, query, sessionId)

	c.logger.Debug("Executed chat insert query", zap.Error(err))
	if err != nil {
		return apperr.Internal(fmt.Errorf("chat repository delete all messages function: %w", err) )
	}

	return nil 
}