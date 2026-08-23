package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	UpdateProfile(context.Context, uuid.UUID, string, string) error
	DeleteAccount(context.Context, uuid.UUID) error
}
