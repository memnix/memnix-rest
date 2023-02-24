package domain

import (
	"strconv"
)

type GithubLogin struct {
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

type DiscordLogin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

func (d *DiscordLogin) ToUser() User {
	return User{
		Username:      d.Username,
		Email:         d.Email,
		Permission:    PermissionUser,
		Avatar:        "https://cdn.discordapp.com/avatars/" + d.ID + "/" + d.Avatar + ".png",
		Oauth:         true,
		OauthProvider: "discord",
		OauthID:       d.ID,
	}
}
