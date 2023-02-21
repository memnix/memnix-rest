package auth

import (
	"strings"

	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/internal/user"
	"github.com/pkg/errors"
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
	/*
		userModel, err := a.GetByEmail(email)
		if err != nil {
			return "", config.ErrUserNotFound
		}

		ok, err := ComparePasswords(password, []byte(userModel.Password))
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || !ok {
			return "", config.ErrInvalidPassword
		}
		if err != nil {
			return "", err
		}

		token, err := jwt.GenerateToken(userModel.ID)
		if err != nil {
			return "", err
		}

		return token, nil

	*/

	return "", nil
}

// Register registers a new user
// Returns an error
func (a *UseCase) Register(user domain.User) (domain.User, error) {
	if err := VerifyPassword(user.Password); err != nil {
		return domain.User{}, errors.Wrap(err, "Verify password failed")
	}

	hash, err := GenerateEncryptedPassword(user.Password)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "Generate encrypted password failed")
	}

	user.Password = string(hash)
	user.Email = strings.ToLower(user.Email)

	if err = a.Create(&user); err != nil {
		return domain.User{}, errors.Wrap(err, "failed to create user in register")
	}

	newUser, err := a.GetByEmail(user.Email)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "failed to get user in register")
	}

	return newUser, nil
}

// Logout returns an empty string
// It might be used to invalidate a token in the future
func (a *UseCase) Logout() (string, error) {
	return "", nil
}

// RefreshToken refreshes a token
func (a *UseCase) RefreshToken(user domain.User) (string, error) {
	/*
		token, err := jwt.GenerateToken(user.ID)
		if err != nil {
			return "", err
		}

		return token, nil

	*/

	return "", nil
}
