package domain

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

// Permission is the permission level of a user.
type Permission int64

type Validator struct {
	validate *validator.Validate
}

var (
	validatorInstance *Validator //nolint:gochecknoglobals //Singleton
	validatorOnce     sync.Once  //nolint:gochecknoglobals //Singleton
)

func GetValidatorInstance() *Validator {
	validatorOnce.Do(func() {
		validatorInstance = &Validator{
			validate: validator.New(),
		}
	})
	return validatorInstance
}

func (v *Validator) Validate() *validator.Validate {
	return v.validate
}

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

type Nonce struct {
	HtmxNonce        string
	TwNonce          string
	HyperscriptNonce string
}
