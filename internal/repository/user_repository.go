package repository

import (
	"context"
	"database/sql"
	"projetpostgre/internal/domain"
	"time"

	"github.com/Masterminds/squirrel"
)

type PostgresUserRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewPostgresUserRepository(db *sql.DB) domain.UserRepository {
	return &PostgresUserRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user domain.User) error {
	query, args, err := r.qb.
		Insert("users").
		Columns("username", "email", "password", "created_at", "updated_at").
		Values(user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	query, args, err := r.qb.
		Select("id", "username", "email", "password", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return user, err
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user domain.User) error {
	query, args, err := r.qb.
		Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// Жесткое удаление
// Перепиши
/*func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.qb.
		Delete("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}*/

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.qb.
		Update("users").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	var users []domain.User

	query, args, err := r.qb.
		Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
