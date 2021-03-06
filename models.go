package gostp

import "time"

// Common is a base model structure
type Common struct {
	ID        uint       `gorm:"primary_key" json:"id" security:"protected"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// User contains minimal information about user
type User struct {
	Common
	Name     string `json:"name"`
	Username string `json:"username" gorm:"type:varchar(100);unique_index" security:"create_only" regex:"username"`
	Password string `json:"password" security:"hidden_out" regex:"password" function:"hashpwd"`
}

// UserTokens contains info about access tokens. Will not be saved in Db
type UserTokens struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessExpiresIn int64  `json:"access_expires_in"`
}

// Access struct contains all access posibilities on site.
type Access struct {
	Common
	UserID uint `json:"user_id" security:"protected"`
	Admin  bool `json:"-"`
	User   bool `json:"-"`
}
