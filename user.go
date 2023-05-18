package twitter

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNameTaken = errors.New("username is already taken")
	ErrEmailTaken    = errors.New("email is already taken")
)

type UserRepo interface {
	Create(ctx context.Context, user User) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
