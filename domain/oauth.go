package domain

import (
	"strconv"

	"gorm.io/gorm"
)

type GithubLogin struct {
	gorm.DB
	Login      string `json:"login"`
	NodeID     string `json:"node_id"`
	AvatarURL  string `json:"avatar_url"`
	GravatarID string `json:"gravatar_id"`
	URL        string `json:"url"`
	Bio        string `json:"bio"`
	Email      string `json:"email"`
	ID         int    `json:"id"`
}

func (g *GithubLogin) ToUser() User {
	return User{
		Username:      g.Login,
		Email:         g.Email,
		Permission:    PermissionUser,
		Avatar:        g.AvatarURL,
		Oauth:         true,
		OauthProvider: "github",
		OauthID:       strconv.Itoa(g.ID),
	}
}
