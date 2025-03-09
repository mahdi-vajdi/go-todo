package auth

import "time"

type User struct {
	Id        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	//UpdatedAt time.Time `json:"updatedAt"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
