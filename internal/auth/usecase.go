package auth

import (
	"context"
	"strings"

	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/internal/user"
	"github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UseCase is the auth use case
type UseCase struct {
	user.IRepository
}

// NewUseCase creates a new auth use case
func NewUseCase(repo user.IRepository) IUseCase {
	return &UseCase{IRepository: repo}
}

// Login logs in a user
// Returns a token and error
func (a *UseCase) Login(ctx context.Context, password string, email string) (string, error) {
	userModel, err := a.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	ok, err := ComparePasswords(ctx, password, []byte(userModel.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || !ok {
		return "", errors.New("invalid password")
	}
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(ctx, userModel.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Register registers a new user
// Returns an error
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
		otelzap.Ctx(ctx).Error("failed to create registerStruct in register", zap.Error(err))
		return domain.User{}, errors.Wrap(err, "failed to create registerStruct in register")
	}

	userModel, err = a.GetByEmail(ctx, registerStruct.Email)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "failed to get registerStruct in register")
	}

	return userModel, nil
}

// Logout returns an empty string
// It might be used to invalidate a token in the future
func (*UseCase) Logout(ctx context.Context) (string, error) {
	return "", nil
}

// RefreshToken refreshes a token
func (*UseCase) RefreshToken(ctx context.Context, user domain.User) (string, error) {
	token, err := jwt.GenerateToken(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// RegisterOauth registers a new user with oauth
func (a *UseCase) RegisterOauth(ctx context.Context, user domain.User) error {
	return a.Create(ctx, &user)
}

// LoginOauth logs in a user with oauth
func (a *UseCase) LoginOauth(ctx context.Context, user domain.User) (string, error) {
	userModel, err := a.GetByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = a.RegisterOauth(ctx, user)
			if err != nil {
				otelzap.Ctx(ctx).Error("failed to register user", zap.Error(err))
				return "", errors.Wrap(err, "failed to register user")
			}

			userModel, err = a.GetByEmail(ctx, user.Email)
			if err != nil {
				otelzap.Ctx(ctx).Error("failed to get user", zap.Error(err))
				return "", errors.New("failed to get user")
			}
		} else {
			otelzap.Ctx(ctx).Error("failed to get user", zap.Error(err))
			return "", err
		}
	}

	if userModel.OauthProvider != user.OauthProvider && userModel.OauthProvider != "" {
		otelzap.Ctx(ctx).Warn("user is already registered with another provider")
		return "", errors.New("user is already registered with another provider")
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
			otelzap.Ctx(ctx).Error("failed to update user", zap.Error(err))
			return "", errors.Wrap(err, "failed to update user")
		}
	}

	return jwt.GenerateToken(ctx, userModel.ID)
}
