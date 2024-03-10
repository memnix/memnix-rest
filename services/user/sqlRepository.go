package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/memnix/memnix-rest/db/sqlc"
)

// SQLRepository is the repository for the user.
type SQLRepository struct {
	q *db.Queries // q is the sqlc queries.
}

// NewRepository returns a new repository.
func NewRepository(dbConn *pgxpool.Pool) IRepository {
	q := db.New(dbConn)
	return &SQLRepository{q: q}
}

// GetName returns the name of the user.
func (r *SQLRepository) GetName(ctx context.Context, id int32) string {
	userName, err := r.q.GetUserName(ctx, id)
	if err != nil {
		return ""
	}
	return userName.String
}

// GetByID returns the user with the given id.
func (r *SQLRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.q.GetUser(ctx, id)
}

// GetByEmail returns the user with the given email.
func (r *SQLRepository) GetByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	return user, err
}

// Create creates a new user.
func (r *SQLRepository) Create(ctx context.Context, email, password, username string) error {
	_, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:    email,
		Password: password,
		Username: pgtype.Text{String: username, Valid: true},
	})
	return err
}

// Update updates the user with the given id.
func (r *SQLRepository) Update(ctx context.Context, id int32, email, password, username string) error {
	_, err := r.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:       id,
		Email:    email,
		Password: password,
		Username: pgtype.Text{String: username, Valid: true},
	})

	return err
}

// Delete deletes the user with the given id.
func (r *SQLRepository) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

// GetByOauthID returns the user with the given oauth id.
func (r *SQLRepository) GetByOauthID(_ context.Context, _ string) (db.User, error) {
	// err := r.DBConn.WithContext(ctx).Where("oauth_id = ?", id).First(&user).Error
	return db.User{}, nil
}
