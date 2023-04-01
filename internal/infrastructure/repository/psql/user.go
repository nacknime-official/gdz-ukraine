package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nacknime-official/gdz-ukraine/internal/entity"
	"github.com/nacknime-official/gdz-ukraine/internal/infrastructure/repository/psql/sqlc"
)

type userRepo struct {
	db *sqlc.Queries
}

func NewUserRepository(db sqlc.DBTX) *userRepo {
	return &userRepo{db: sqlc.New(db)}
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := r.db.GetUserByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return user.ToEntity(), nil
}

func (r *userRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error) {
	user, err := r.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return user.ToEntity(), nil
}
