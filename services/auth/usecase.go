package auth

import (
	"context"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/services/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UseCase is the auth use case.
type UseCase struct {
	user.IRepository
}

// NewUseCase creates a new auth use case.
func NewUseCase(repo user.IRepository) IUseCase {
	return &UseCase{IRepository: repo}
}

// Login logs in a user
// Returns a token and error.
func (a *UseCase) Login(ctx context.Context, password string, email string) (string, error) {
	userModel, err := a.GetByEmail(ctx, email)
	if err != nil {
		log.WithContext(ctx).Error("user not found", slog.Any("error", err), slog.String("email", email))
		return "", errors.New("user not found")
	}

	ok, err := ComparePasswords(ctx, password, []byte(userModel.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || !ok {
		return "", errors.New("invalid password")
	}
	if err != nil {
		return "", err
	}

	token, err := config.GetJwtInstance().GetJwt().GenerateToken(ctx, userModel.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Register registers a new user
// Returns an error.
func (a *UseCase) Register(ctx context.Context, registerStruct domain.Register) (domain.User, error) {
	if err := VerifyPassword(registerStruct.Password); err != nil {
		return domain.User{}, errors.Wrap(err, "Verify password failed")
	}

	hash, err := GenerateEncryptedPassword(ctx, registerStruct.Password)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "Generate encrypted password failed")
	}

	registerStruct.Password = string(hash)
	registerStruct.Email = strings.ToLower(registerStruct.Email)
	userModel := registerStruct.ToUser()

	if err = a.Create(ctx, &userModel); err != nil {
		log.WithContext(ctx).Error("failed to create registerStruct in register", slog.Any("error", err))
		return domain.User{}, errors.Wrap(err, "failed to create registerStruct in register")
	}

	userModel, err = a.GetByEmail(ctx, registerStruct.Email)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "failed to get registerStruct in register")
	}

	return userModel, nil
}

// Logout returns an empty string
// It might be used to invalidate a token in the future.
func (*UseCase) Logout(_ context.Context) (string, error) {
	return "", nil
}

// RefreshToken refreshes a token.
func (*UseCase) RefreshToken(ctx context.Context, user domain.User) (string, error) {
	token, err := config.GetJwtInstance().GetJwt().GenerateToken(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// RegisterOauth registers a new user with oauth.
func (a *UseCase) RegisterOauth(ctx context.Context, user domain.User) error {
	return a.Create(ctx, &user)
}

func (a *UseCase) LoginOauth(ctx context.Context, user domain.User) (string, error) {
	userModel, err := a.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithContext(ctx).Error("failed to get user", slog.Any("error", err))
		return "", err
	}

	if err == nil && userModel.OauthProvider != user.OauthProvider && userModel.OauthProvider != "" {
		log.WithContext(ctx).Warn("user is already registered with another provider")
		return "", errors.New("user is already registered with another provider")
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = a.RegisterOauth(ctx, user)
		if err != nil {
			log.WithContext(ctx).Error("failed to register user", slog.Any("error", err))
			return "", errors.Wrap(err, "failed to register user")
		}

		userModel, err = a.GetByEmail(ctx, user.Email)
		if err != nil {
			log.WithContext(ctx).Error("failed to get user", slog.Any("error", err))
			return "", errors.New("failed to get user")
		}
	}

	// Check if user is up to date
	if userModel.OauthID == "" {
		userModel.OauthID = user.OauthID
		userModel.OauthProvider = user.OauthProvider
		if user.Avatar != "" {
			userModel.Avatar = user.Avatar
		}
		err = a.Update(ctx, &userModel)
		if err != nil {
			log.WithContext(ctx).Error("failed to update user", slog.Any("error", err))
			return "", errors.Wrap(err, "failed to update user")
		}
	}

	return config.GetJwtInstance().GetJwt().GenerateToken(ctx, userModel.ID)
}
