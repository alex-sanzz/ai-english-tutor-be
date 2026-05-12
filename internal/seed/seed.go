package seed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Seeder struct{
	pgPool *pgxpool.Pool
}

func NewSeeder(pgPool *pgxpool.Pool) *Seeder {
	return &Seeder{
		pgPool: pgPool,
	}
}

func (s *Seeder) SeedRoomType(ctx context.Context) error{
	query := "INSERT INTO room_types (type) VALUES($1) ON CONFLICT(type) DO NOTHING"

	types := []string{"basic", "question", "roleplay"}

	for idx := range(types) {
		_, err := s.pgPool.Exec(ctx, query, types[idx])

		if err != nil {
			return fmt.Errorf("seed room type function error: %w", err) 
		}
	}

	return nil
	
}