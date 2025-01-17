package todo

import (
	"time"
)

type Todo struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
