package auth

import (
	"strings"

	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/internal/user"
	"github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UseCase struct {
	user.IRepository
}

func NewUseCase(repo user.IRepository) IUseCase {
	return &UseCase{IRepository: repo}
}

// Login logs in a user
// Returns a token and error
func (a *UseCase) Login(password string, email string) (string, error) {
	userModel, err := a.GetByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	ok, err := ComparePasswords(password, []byte(userModel.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || !ok {
		return "", errors.New("invalid password")
	}
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(userModel.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Register registers a new user
// Returns an error
func (a *UseCase) Register(registerStruct domain.Register) (domain.User, error) {
	if err := VerifyPassword(registerStruct.Password); err != nil {
		return domain.User{}, errors.Wrap(err, "Verify password failed")
	}

	hash, err := GenerateEncryptedPassword(registerStruct.Password)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "Generate encrypted password failed")
	}

	registerStruct.Password = string(hash)
	registerStruct.Email = strings.ToLower(registerStruct.Email)
	userModel := registerStruct.ToUser()

	if err = a.Create(&userModel); err != nil {
		log.Warn().Err(err).Msg("failed to create user in register")
		return domain.User{}, errors.Wrap(err, "failed to create registerStruct in register")
	}

	userModel, err = a.GetByEmail(registerStruct.Email)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "failed to get registerStruct in register")
	}

	return userModel, nil
}

// Logout returns an empty string
// It might be used to invalidate a token in the future
func (a *UseCase) Logout() (string, error) {
	return "", nil
}

// RefreshToken refreshes a token
func (a *UseCase) RefreshToken(user domain.User) (string, error) {
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *UseCase) RegisterOauth(user domain.User) error {
	return a.Create(&user)
}

func (a *UseCase) LoginOauth(user domain.User) (string, error) {
	userModel, err := a.GetByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = a.RegisterOauth(user)
			if err != nil {
				log.Error().Err(err).Msg("failed to register user")
				return "", errors.New("failed to register user")
			}
		} else {
			log.Error().Err(err).Msg("failed to get user")
			return "", errors.New("failed to get user")
		}
	}

	if userModel.OauthProvider != user.OauthProvider && userModel.OauthProvider != "" {
		log.Error().Err(err).Msg("user is already registered with another provider")
		return "", errors.New("user is already registered with another provider")
	}

	// Check if user is up to date
	if userModel.OauthID == "" {
		userModel.OauthID = user.OauthID
		userModel.OauthProvider = user.OauthProvider
		if user.Avatar != "" {
			userModel.Avatar = user.Avatar
		}
		err = a.Update(&userModel)
		if err != nil {
			log.Error().Err(err).Msg("failed to update user")
			return "", errors.New("failed to update user")
		}
	}

	return jwt.GenerateToken(userModel.ID)
}
