package postgres

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type conversationQuestionRepository struct {
	db     *pgxpool.Pool
	logger log.Logger
}

func NewConversationQuestionRepository(db *pgxpool.Pool, logger log.Logger) repository.ConversationQuestionRepository {
	return &conversationQuestionRepository{
		db:     db,
		logger: logger,
	}
}

func (cq conversationQuestionRepository) FindAll(ctx context.Context, sessionRoomId string, limit *int, desc bool) ([]*models.ConversationQuestion, error) {
	query := "SELECT id, question, answer, review_result, session_room_id, answered_at, created_at FROM conversation_questions WHERE session_room_id = $1 ORDER BY created_at "

	args := []interface{}{}

	args = append(args, sessionRoomId)

	if desc {
		query += "DESC"
	} else {
		query += "ASC"
	}

	if limit != nil {
		query += " LIMIT $2"

		args = append(args, limit)
	}

	rows, err := cq.db.Query(ctx, query, args...)

	if err != nil {
		return nil, apperr.Internal(fmt.Errorf("conversation question repository find all error: %w", err))
	}
	defer rows.Close()

	questions := []*models.ConversationQuestion{}

	for rows.Next() {
		var q models.ConversationQuestion
		var answer sql.NullString
		var reviewResult sql.NullString
		var answeredAt sql.NullTime

		if err := rows.Scan(&q.ID, &q.Question, &answer, &reviewResult, &q.SessionRoomId, &answeredAt, &q.CreatedAt); err != nil {
			return nil, apperr.Internal(apperr.Internal(fmt.Errorf("conversation question repository find all error: %w", err)))
		}

		if answer.Valid {
			q.Answer = answer.String
		}

		if reviewResult.Valid {
			q.ReviewResult = reviewResult.String
		}

		if answeredAt.Valid {
			q.AnsweredAt = answeredAt.Time
		}

		questions = append(questions, &q)
	}

	if err := rows.Err(); err != nil {
		return nil, apperr.Internal(apperr.Internal(fmt.Errorf("conversation question repository find all error: %w", err)))
	}

	return questions, nil
}

func (cq conversationQuestionRepository) FindById(ctx context.Context, id string) (*models.ConversationQuestion, error) {
	query := "SELECT id, question, answer, review_result, session_room_id, answered_at, created_at FROM conversation_questions WHERE id = $1"
	
	var c models.ConversationQuestion

	var answer sql.NullString

	var reviewResult sql.NullString

	var answeredAt sql.NullTime

	err := cq.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Question, &answer, &reviewResult, &c.SessionRoomId, &answeredAt, &c.CreatedAt)

	if err != nil {
		return nil, apperr.Internal(fmt.Errorf("conversation question repository find by id error: %w", err))
	}

	if answer.Valid {
		c.Answer = answer.String
	}

	if reviewResult.Valid {
		c.ReviewResult = reviewResult.String
	}

	if answeredAt.Valid{
		c.AnsweredAt = answeredAt.Time
	}

	return &c, nil 

}

func (cq conversationQuestionRepository) Create(ctx context.Context, m *models.ConversationQuestion) (string, error) {
	query := "INSERT INTO conversation_questions(question, session_room_id) VALUES ($1) RETURNING id"
	var id string
	err := cq.db.QueryRow(ctx, query, m.Question, m.SessionRoomId).Scan(&id)

	if err != nil {
		return "", apperr.Internal(fmt.Errorf("conversation question repository create error: %w", err))
	}

	return id, err
}

func (cq conversationQuestionRepository) CreateBatch(ctx context.Context, questions []*models.ConversationQuestion) error {

	if len(questions) == 0 {
		return fmt.Errorf("conversation question repository create batch error: questions is empty")
	}

	query := "INSERT INTO conversation_questions(question, session_room_id, created_at) VALUES "

	parameters := []string{}

	values := []interface{}{}

	i := 1

	for _, q := range questions {

		parameters = append(parameters, fmt.Sprintf("($%d, $%d, $%d)", i, i+1, i+2))

		values = append(values, q.Question, q.SessionRoomId, q.CreatedAt)

		i += 3

	}

	query = query + strings.Join(parameters, ",")

	_, err := cq.db.Exec(ctx, query, values...)

	if err != nil {
		return apperr.Internal(fmt.Errorf("conversation question repository create batch error: %w", err))
	}

	return nil
}

func (cq conversationQuestionRepository) CheckIfThereIsAnUnansweredQuestion(ctx context.Context, sessionRoomId string) (bool, error) {
	// EXIST is subquery which has job to check whether there is any data, if it's then return true, otherwise return false
	// select 1 means return 1 if there is any that row that is matched with criteria
	query := `SELECT EXISTS (SELECT 1 FROM conversation_questions WHERE session_room_id = $1)`

	var exist bool

	if err := cq.db.QueryRow(ctx, query, sessionRoomId).Scan(&exist); err != nil {
		return false, apperr.Internal(fmt.Errorf("conversation question repository check if there is an unanswered question error: %w", err))
	}

	return exist, nil
}

func (cq conversationQuestionRepository) FindOneUnansweredQuestion(ctx context.Context, sessionRoomId string) (*models.ConversationQuestion, error) {

	query := `
        SELECT id, question, answer, review_result, session_room_id, answered_at, created_at
        FROM conversation_questions
        WHERE session_room_id = $1 AND answer IS NULL
        ORDER BY created_at ASC
        LIMIT 1
    `
	rows, err := cq.db.Query(ctx, query, sessionRoomId)

	if err != nil {
		return nil, apperr.Internal(fmt.Errorf("Find one unanswered question error: %w", err))
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var data models.ConversationQuestion
	var answer sql.NullString
	var reviewResult sql.NullString
	var answeredAt sql.NullTime

	if err := rows.Scan(
		&data.ID,
		&data.Question,
		&answer,
		&reviewResult,
		&data.SessionRoomId,
		&answeredAt,
		&data.CreatedAt,
	); err != nil {
		return nil, apperr.Internal(fmt.Errorf("Find one unanswered question error: %w", err))
	}

	if answer.Valid {
		data.Answer = answer.String
	}

	if reviewResult.Valid {
		data.ReviewResult = reviewResult.String
	}

	if answeredAt.Valid {
		data.AnsweredAt = answeredAt.Time
	}

	if err := rows.Err(); err != nil {
		return nil, apperr.Internal(fmt.Errorf("Find one unanswered question error: %w", err))
	}

	return &data, nil
}

func (cq conversationQuestionRepository) UpdateAnswer(ctx context.Context, id, answer, reviewResult string) error {
	query := "UPDATE conversation_questions SET answer = $1, review_result = $2, answered_at = $3 WHERE id = $4"
	cd, err := cq.db.Exec(ctx, query, answer, reviewResult, time.Now(), id)

	if err != nil {
		return apperr.Internal(fmt.Errorf("conversation question repository update answer error: %w", err))
	}

	if cd.RowsAffected() == 0 {
		return apperr.BadRequest("400", "the requested question deletion is not found", fmt.Errorf("%s", "conversation question with an id '"+id+"' is not found"))
	}

	return nil
}

func (cq conversationQuestionRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM conversation_questions WHERE id = $1"
	cd, err := cq.db.Exec(ctx, query, id)

	if err != nil {
		return apperr.Internal(fmt.Errorf("conversation question repository delete error: %w", err))
	}

	if cd.RowsAffected() == 0 {
		return apperr.BadRequest("400", "the requested question deletion is not found", fmt.Errorf("%s", "conversation question with an id '"+id+"' is not found"))
	}

	return nil
}
