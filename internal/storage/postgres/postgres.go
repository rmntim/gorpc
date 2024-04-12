package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rmntim/sso/internal/domain/models"
	"github.com/rmntim/sso/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Migrate() error {
	const op = "storage.postgres.Migrate"

	driver, err := postgres.WithInstance(s.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var id int64
	query := "INSERT INTO users(email, pass_hash) VALUES ($1, $2) RETURNING id"
	if err := s.db.QueryRowContext(ctx, query, email, passwordHash).Scan(&id); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.postgres.User"

	var user models.User

	query := "SELECT * FROM users u WHERE u.email = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := stmt.QueryRowContext(ctx, email).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.IsAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var isAdmin bool

	query := "SELECT is_admin FROM users u WHERE u.id = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if err := stmt.QueryRowContext(ctx, userId).Scan(&isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appId int) (*models.App, error) {
	const op = "storage.postgres.App"

	var app models.App

	query := "SELECT * FROM apps a WHERE a.id = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := stmt.QueryRowContext(ctx, appId).Scan(&app.Id, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &app, nil
}
