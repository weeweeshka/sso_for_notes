package auth

import (
	"context"
	"fmt"
	"github.com/weeweeshka/sso_for_notes/internal/domain/models"
	"github.com/weeweeshka/sso_for_notes/internal/lib/jwt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	authRegister AuthRegistration
	userProvider UserProvider
	appProvider  AppProvider
	tokenTLL     time.Duration
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type AuthRegistration interface {
	Register(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

func New(log *slog.Logger,
	authRegister AuthRegistration,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTLL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		authRegister: authRegister,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTLL:     tokenTLL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "businessLogic.auth.Register"

	log := a.log.With(slog.String("op", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := a.authRegister.Register(ctx, email, passHash)
	if err != nil {
		log.Info("failed to register user", err)
		return 0, err
	}

	log.Info("user was registered!")

	return userID, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "businessLogic.auth.Login"

	log := a.log.With(slog.String("op", op))

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		log.Error("failed to find user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("failed to find app", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User logged successfully!")

	token, err := jwt.NewToken(user, app, a.tokenTLL)
	if err != nil {
		log.Error("failed to generate token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil

}
