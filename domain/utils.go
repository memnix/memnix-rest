package domain

import "github.com/gofiber/fiber/v2"

// Permission is the permission level of a user.
type Permission int64

const (
	PermissionNone  Permission = iota // PermissionNone is the default permission level.
	PermissionUser                    // PermissionUser is the permission level of a user.
	PermissionVip                     // PermissionVip is the permission level of a vip.
	PermissionAdmin                   // PermissionAdmin is the permission level of an admin.
)

func (p Permission) String() string {
	return [...]string{"none", "user", "vip", "admin"}[p]
}

func (p Permission) IsValid() bool {
	return p >= PermissionNone && p <= PermissionAdmin
}

// Route is a route for the API
// It contains the handler, method and permission level.
type Route struct {
	Handler    func(c *fiber.Ctx) error // Handler is the handler function for the route.
	Method     string                   // Method is the method of the route.
	Permission Permission               // Permission is the permission level of the route.
}

type Model interface {
	TableName() string
}
