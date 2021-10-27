package models

// ResponseHTTP structure
type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Count   int         `json:"count"`
}

type ResponseCard struct {
	Card    Card
	Answers []string
}

type ResponseAuth struct {
	Success bool
	User    User
	Message string
}

type Permission int64

const (
	PermUser Permission = iota
	PermMod
	PermAdmin
)
