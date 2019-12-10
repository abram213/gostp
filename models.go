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
	Email    string `json:"email" gorm:"type:varchar(100);unique_index" security:"create_only" regex:"email"`
	Password string `json:"password" security:"hidden_out" regex:"password" function:"hashpwd"`
	Token    Token  `json:"token" security:"protected"`
}

// UserTokens contains info about access tokens. Will not be saved in Db
type UserTokens struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessExpiresIn int64  `json:"access_expires_in"`
}

// Token - structure which contains info about token
type Token struct {
	ID           uint       `gorm:"primary_key" json:"-"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `sql:"index" json:"-"`
	UserID       uint       `json:"-"`
	RefreshToken string     `gorm:"type:varchar(255)" json:"refresh_token"`
}
