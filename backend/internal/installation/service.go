package installation

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo  Repository
	token string
}

func NewService(repo Repository, token string) *Service { return &Service{repo: repo, token: token} }
func (s *Service) Status(ctx context.Context) (bool, bool, error) {
	required, err := s.repo.Required(ctx)
	return required, s.token != "", err
}
func (s *Service) Install(ctx context.Context, token, email, password, nickname string) (Admin, error) {
	required, err := s.repo.Required(ctx)
	if err != nil {
		return Admin{}, err
	}
	if !required {
		return Admin{}, ErrAlreadyInstalled
	}
	if s.token == "" {
		return Admin{}, ErrSetupDisabled
	}
	expected := sha256.Sum256([]byte(s.token))
	received := sha256.Sum256([]byte(token))
	if subtle.ConstantTimeCompare(expected[:], received[:]) != 1 {
		return Admin{}, ErrInvalidToken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Admin{}, err
	}
	seed := make([]byte, 16)
	if _, err = rand.Read(seed); err != nil {
		return Admin{}, err
	}
	return s.repo.CreateFirstAdmin(ctx, NewAdmin{ID: uuid.New(), Email: strings.ToLower(strings.TrimSpace(email)), PasswordHash: string(hash), Nickname: strings.TrimSpace(nickname), AvatarSeed: base64.RawURLEncoding.EncodeToString(seed)})
}
