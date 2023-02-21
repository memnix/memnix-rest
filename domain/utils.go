package domain

import "github.com/gofiber/fiber/v2"

type Permission int64

const (
	PermissionNone Permission = iota
	PermissionUser
	PermissionVip
	PermissionAdmin
)

type Route struct {
	Handler    func(c *fiber.Ctx) error
	Method     string
	Permission Permission
}
