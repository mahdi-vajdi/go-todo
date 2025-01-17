package auth

import "time"

type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
