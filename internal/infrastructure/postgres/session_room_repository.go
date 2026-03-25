package postgres

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRoomRepository struct {
	db *pgxpool.Pool
	logger log.Logger
}

func NewSessionRoomRepository(db *pgxpool.Pool, logger log.Logger) *SessionRoomRepository {
	return &SessionRoomRepository{
		db: db,
		logger: logger,
	}
}

func (s *SessionRoomRepository) CreateSessionRoom(ctx context.Context, userId, roomType, icon, topic string) (string, error) {
	var id string
	// On conflict will only be triggered when primary key, unique constraint, unique index is violated
	err := s.db.QueryRow(ctx, "INSERT INTO session_rooms (user_id, room_type_id, icon, topic) VALUES ($1, $2, $3, $4) RETURNING id", userId, roomType, icon, topic).Scan(&id)

	if err != nil {
		
		return "", fmt.Errorf("postgres CreateSessionRoom error: %w", err)
	}

	return id, nil
}

func (s *SessionRoomRepository) FindById(ctx context.Context, id string) (*models.SessionRoom, error){
	query := "SELECT id, room_type_id, icon, topic, user_id, created_at FROM session_rooms WHERE id = $1"
	row := s.db.QueryRow(ctx, query, id)

	var room models.SessionRoom

	err := row.Scan(&room.ID, &room.RoomType, &room.Icon, &room.Topic, &room.UserId, &room.CreatedAt)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows){
			return nil, apperr.NotFound("400", "session room not found", err)
		}

		return nil, apperr.Internal(fmt.Errorf("session room repository find by id error: %w", err))
	}

	return &room, nil 
}

func (s *SessionRoomRepository) FindAllByUserIdAndRoomType(ctx context.Context, userId, roomType string) ([]*models.SessionRoom, error){
	query := "SELECT id, room_type_id, topic, icon, user_id, created_at FROM session_rooms WHERE user_id = $1 AND room_type_id = $2"

	rows, err := s.db.Query(ctx, query, userId, roomType)

	if err != nil {
		return nil, fmt.Errorf("session room repository find all by user id and room type error : %w", err)
	}

	defer rows.Close()

	var rooms []*models.SessionRoom

	for rows.Next() {
		var r models.SessionRoom

		if err := rows.Scan(&r.ID, &r.RoomType, &r.Topic, &r.Icon, &r.UserId, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("session room repository find all by user id and room type error : %w", err)
		}

		rooms = append(rooms, &r)

	}

	return rooms, nil 

}

func (s *SessionRoomRepository) DeleteById(ctx context.Context, id string) error {
	query := "DELETE FROM session_rooms WHERE id = $1"
	ct, err := s.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("session room repository delete by id error: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return apperr.NotFound("404", "the session room is not exist", nil)
	}

	return nil
}