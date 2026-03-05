package repository

import (
	"amass/internal/models"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type IUserRepository interface {
	GetUser(u *models.UserLogin) (*models.User, error)
	Create(u *models.User) error
}

type userRepository struct {
	ctx    context.Context
	client *pgxpool.Pool
}

func NewUserRepository(ctx context.Context, client *pgxpool.Pool) (IUserRepository, error) {
	return &userRepository{
		ctx:    ctx,
		client: client,
	}, nil
}

func (r *userRepository) GetUser(u *models.UserLogin) (*models.User, error) {

	query := `SELECT username, password_hash, status, created_at, updated_at FROM users WHERE username = $1`

	row := r.client.QueryRow(r.ctx, query, u.Username)

	var user models.User
	err := row.Scan(&user.Username, &user.PasswordHash, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(u *models.User) error {

	query := `INSERT INTO users (username, password_hash, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.client.Exec(r.ctx,
		query,
		u.Username,
		u.PasswordHash,
		u.Status,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
