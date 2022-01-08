package domain

import "time"

// User ...
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json: "email"`
	Password  string    `json:"password"`
	Salt      string    `json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
