package domain

// User is the domain model for a user.
type User struct {
	Username      string     `json:"username" validate:"required"`
	Email         string     `json:"email" validate:"email" gorm:"unique"`
	Password      string     `json:"-" validate:"required"`
	Avatar        string     `json:"avatar"`
	OauthProvider string     `json:"oauth_provider" `
	OauthID       string     `json:"oauth_id" gorm:"unique"`
	Learning      []*Deck    `json:"learning" gorm:"many2many:user_decks;"`
	OwnDecks      []Deck     `json:"own_decks" gorm:"foreignKey:OwnerID"`
	ID            uint       `json:"id"`
	Permission    Permission `json:"permission"`
	Oauth         bool       `json:"oauth" gorm:"default:false"`
}

// TableName returns the table name for the user model.
func (*User) TableName() string {
	return "users"
}

// ToPublicUser converts the user to a public user.
func (u *User) ToPublicUser() PublicUser {
	return PublicUser{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Avatar:     u.Avatar,
		Permission: u.Permission,
	}
}

// Validate validates the user.
func (u *User) Validate() error {
	return GetValidatorInstance().Validate().Struct(u)
}

// HasPermission checks if the user has the given permission.
func (u *User) HasPermission(permission Permission) bool {
	return u.Permission >= permission
}

// PublicUser is the public user model.
type PublicUser struct {
	Username   string     `json:"username"`   // Username of the user
	Email      string     `json:"email"`      // Email of the user
	Avatar     string     `json:"avatar"`     // Avatar of the user
	ID         uint       `json:"id"`         // ID of the user
	Permission Permission `json:"permission"` // Permission of the user
}

// Login is the login model.
type Login struct {
	Email    string `json:"email" validate:"required,email,max=254"` // Email of the user
	Password string `json:"password" validate:"required,max=72"`     // Password of the user
}

// Validate validates the login model.
func (l *Login) Validate() error {
	return GetValidatorInstance().Validate().Struct(l)
}

// Register is the register model.
type Register struct {
	Username string `json:"username" validate:"required,max=50,min=3"` // Username of the user
	Email    string `json:"email" validate:"email,required,max=254"`   // Email of the user
	// Password of the user
	// min length is 12 as recommended by the OWASP password guidelines
	// max length is 72 because of bcrypt limitation
	// the password must have a minimum entropy of 72 bits but the recommended is 80 bits (CNIL, 2022)
	Password string `json:"password" validate:"required,min=10,max=72"`
}

// Validate validates the register model.
func (r *Register) Validate() error {
	return GetValidatorInstance().Validate().Struct(r)
}

// ToUser converts the register model to a user.
func (r *Register) ToUser() User {
	return User{
		Username:   r.Username,
		Email:      r.Email,
		Password:   r.Password,
		Permission: PermissionUser,
	}
}
