package user

import (
	"github.com/google/uuid"
	"time"
)

type Role string
type User struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Version   uint      `json:"version"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      Role      `json:"role"`
}

type DBUser struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   uint      `db:"version"`
	Phone     string    `db:"phone"`
	Email     string    `json:"email"`
	Password  []byte    `db:"password"`
	Role      Role      `db:"role"`
}
