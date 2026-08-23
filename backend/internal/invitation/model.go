package invitation

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID        uuid.UUID  `json:"id"`
	Label     string     `json:"label"`
	ExpiresAt *time.Time `json:"expiresAt"`
	MaxUses   int        `json:"maxUses"`
	Uses      int        `json:"uses"`
	CreatedAt time.Time  `json:"createdAt"`
	RevokedAt *time.Time `json:"revokedAt"`
}

type Repository interface {
	Create(context.Context, Invitation, []byte, uuid.UUID) (Invitation, error)
	List(context.Context) ([]Invitation, error)
	Revoke(context.Context, uuid.UUID) error
}
