package installation

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrAlreadyInstalled = errors.New("installation is already initialized")
	ErrSetupDisabled    = errors.New("setup token is not configured")
	ErrInvalidToken     = errors.New("invalid setup token")
)

type Admin struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Nickname string    `json:"nickname"`
	Role     string    `json:"role"`
}
type NewAdmin struct {
	ID                                        uuid.UUID
	Email, PasswordHash, Nickname, AvatarSeed string
}
type Repository interface {
	Required(context.Context) (bool, error)
	CreateFirstAdmin(context.Context, NewAdmin) (Admin, error)
}
