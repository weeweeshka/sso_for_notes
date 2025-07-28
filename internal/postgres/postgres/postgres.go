package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/weeweeshka/sso_for_notes/internal/domain/models"
	"log/slog"
	"time"
)

type Storage struct {
	db *pgxpool.Pool
}

func configurationPool(config *pgxpool.Config) {
	config.MaxConns = int32(20)
	config.MinConns = int32(5)
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = 5 * time.Second
}

func NewStorage(connString string, log *slog.Logger) (*Storage, error) {
	const op = "postgres.NewStorage"
	log = slog.With(slog.String("op", op))
	ctx := context.Context(context.Background())

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error("failed to parse config", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	configurationPool(config)

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error("failed to connect to database", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("connected to database")

	_, err = dbPool.Exec(ctx, `CREATE TABLE IF NOT EXISTS users(
    	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    	email TEXT NOT NULL UNIQUE,
    	pass_hash BYTEA NOT NULL)`)
	if err != nil {
		log.Error("failed to create table", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("created table users")

	_, err = dbPool.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS apps(
        id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
     	name TEXT NOT NULL UNIQUE,
     	secret TEXT NOT NULL UNIQUE)`)
	if err != nil {
		log.Error("failed to create table", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("created table apps")

	return &Storage{db: dbPool}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) Register(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "postgres.Register"
	log := slog.With(slog.String("op", op))

	var id int64
	err := s.db.QueryRow(ctx, `INSERT INTO users(email, pass_hash) VALUES ($1, $2) RETURNING id`, email, passHash).Scan(&id)
	if err != nil {
		log.Error("failed to insert user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "postgres.Login"
	log := slog.With(slog.String("op", op))

	row := s.db.QueryRow(ctx, `SELECT id, email, pass_hash FROM users WHERE email = $1`, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("found user")
	return user, nil

}

func (s *Storage) App(ctx context.Context, appID int32) (models.App, error) {
	const op = "storage.sqlite.App"

	row := s.db.QueryRow(ctx, `SELECT id, name, secret FROM apps WHERE id = $1`, appID)

	var app models.App
	err := row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) CreateApp(ctx context.Context, appName string, secret string) (int32, error) {
	const op = "postgres.CreateApp"
	log := slog.With(slog.String("op", op))

	var appID int32
	err := s.db.QueryRow(
		ctx,
		`INSERT INTO apps(name, secret) VALUES ($1, $2) RETURNING id`,
		appName, secret,
	).Scan(&appID)

	if err != nil {
		log.Error("failed to create app", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("created app")
	return appID, nil
}
