package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Nickname     string    `json:"nickname"`
	AvatarSeed   string    `json:"avatarSeed"`
	Theme        string    `json:"theme"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}

type NewUser struct {
	ID, InvitationID                          uuid.UUID
	Email, PasswordHash, Nickname, AvatarSeed string
}

type Repository interface {
	FindUserByEmail(context.Context, string) (User, error)
	FindUserByID(context.Context, uuid.UUID) (User, error)
	RegisterWithInvitation(context.Context, NewUser, []byte, time.Time) (User, error)
	CreateSession(context.Context, uuid.UUID, []byte, time.Time, string, string) error
	FindUserBySession(context.Context, []byte) (User, error)
	DeleteSession(context.Context, []byte) error
	CreatePasswordReset(context.Context, uuid.UUID, uuid.UUID, []byte, time.Time) error
	ConsumePasswordReset(context.Context, []byte, string, time.Time) error
}

type PasswordReset struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user"`
}
