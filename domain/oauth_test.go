package domain_test

import (
	"testing"

	"github.com/memnix/memnix-rest/domain"
)

func TestGithubLogin_ToUser(t *testing.T) {
	githubLogin := &domain.GithubLogin{
		Login:     "testuser",
		Email:     "testuser@example.com",
		AvatarURL: "https://example.com/avatar.png",
		ID:        123,
	}

	expected := domain.User{
		Username:      "testuser",
		Email:         "testuser@example.com",
		Permission:    domain.PermissionUser,
		Avatar:        "https://example.com/avatar.png",
		Oauth:         true,
		OauthProvider: "github",
		OauthID:       "123",
	}

	user := githubLogin.ToUser()

	if user.Username != expected.Username {
		t.Errorf("Expected username to be %s, but got %s", expected.Username, user.Username)
	}

	if user.Email != expected.Email {
		t.Errorf("Expected email to be %s, but got %s", expected.Email, user.Email)
	}

	if user.Permission != expected.Permission {
		t.Errorf("Expected permission to be %s, but got %s", expected.Permission, user.Permission)
	}

	if user.Avatar != expected.Avatar {
		t.Errorf("Expected avatar to be %s, but got %s", expected.Avatar, user.Avatar)
	}

	if user.Oauth != expected.Oauth {
		t.Errorf("Expected oauth to be %v, but got %v", expected.Oauth, user.Oauth)
	}

	if user.OauthProvider != expected.OauthProvider {
		t.Errorf("Expected oauth provider to be %s, but got %s", expected.OauthProvider, user.OauthProvider)
	}

	if user.OauthID != expected.OauthID {
		t.Errorf("Expected oauth ID to be %s, but got %s", expected.OauthID, user.OauthID)
	}
}

func TestDiscordLogin_ToUser(t *testing.T) {
	discordLogin := &domain.DiscordLogin{
		Username: "testuser",
		Email:    "testuser@example.com",
		ID:       "123",
		Avatar:   "avatar123",
	}

	expected := domain.User{
		Username:      "testuser",
		Email:         "testuser@example.com",
		Permission:    domain.PermissionUser,
		Avatar:        "https://cdn.discordapp.com/avatars/123/avatar123.png",
		Oauth:         true,
		OauthProvider: "discord",
		OauthID:       "123",
	}

	user := discordLogin.ToUser()

	if user.Username != expected.Username {
		t.Errorf("Expected username to be %s, but got %s", expected.Username, user.Username)
	}

	if user.Email != expected.Email {
		t.Errorf("Expected email to be %s, but got %s", expected.Email, user.Email)
	}

	if user.Permission != expected.Permission {
		t.Errorf("Expected permission to be %s, but got %s", expected.Permission, user.Permission)
	}

	if user.Avatar != expected.Avatar {
		t.Errorf("Expected avatar to be %s, but got %s", expected.Avatar, user.Avatar)
	}

	if user.Oauth != expected.Oauth {
		t.Errorf("Expected oauth to be %v, but got %v", expected.Oauth, user.Oauth)
	}

	if user.OauthProvider != expected.OauthProvider {
		t.Errorf("Expected oauth provider to be %s, but got %s", expected.OauthProvider, user.OauthProvider)
	}

	if user.OauthID != expected.OauthID {
		t.Errorf("Expected oauth ID to be %s, but got %s", expected.OauthID, user.OauthID)
	}
}
