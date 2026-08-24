package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrPasswordResetInvalid = errors.New("password reset invalid")

const passwordResetTTL = 30 * time.Minute

type Clock func() time.Time

type Service struct {
	repo Repository
	ttl  time.Duration
	now  Clock
}

func NewService(repo Repository, ttl time.Duration) *Service {
	return &Service{repo: repo, ttl: ttl, now: time.Now}
}

func normalizeEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }
func tokenHash(token string) []byte      { sum := sha256.Sum256([]byte(token)); return sum[:] }

func invitationHash(code string) []byte { return tokenHash(strings.ToUpper(strings.TrimSpace(code))) }

func randomToken(bytes int) (string, error) {
	buffer := make([]byte, bytes)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func (s *Service) Register(ctx context.Context, email, password, nickname, invitationCode string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	seed, err := randomToken(16)
	if err != nil {
		return User{}, err
	}
	return s.repo.RegisterWithInvitation(ctx, NewUser{ID: uuid.New(), Email: normalizeEmail(email), PasswordHash: string(hash), Nickname: strings.TrimSpace(nickname), AvatarSeed: seed}, invitationHash(invitationCode), s.now())
}

func (s *Service) Login(ctx context.Context, email, password, userAgent, ip string) (User, string, error) {
	user, err := s.repo.FindUserByEmail(ctx, normalizeEmail(email))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return User{}, "", ErrInvalidCredentials
	}
	token, err := randomToken(32)
	if err != nil {
		return User{}, "", err
	}
	if err := s.repo.CreateSession(ctx, user.ID, tokenHash(token), s.now().Add(s.ttl), userAgent, ip); err != nil {
		return User{}, "", err
	}
	return user, token, nil
}

func (s *Service) Authenticate(ctx context.Context, token string) (User, error) {
	if token == "" {
		return User{}, ErrInvalidCredentials
	}
	return s.repo.FindUserBySession(ctx, tokenHash(token))
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.repo.DeleteSession(ctx, tokenHash(token))
}

func (s *Service) CreatePasswordReset(ctx context.Context, email string, createdBy uuid.UUID) (PasswordReset, error) {
	user, err := s.repo.FindUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return PasswordReset{}, err
	}
	token, err := randomToken(32)
	if err != nil {
		return PasswordReset{}, err
	}
	expiresAt := s.now().Add(passwordResetTTL)
	if err := s.repo.CreatePasswordReset(ctx, user.ID, createdBy, tokenHash(token), expiresAt); err != nil {
		return PasswordReset{}, err
	}
	return PasswordReset{Token: token, ExpiresAt: expiresAt, User: user}, nil
}

func (s *Service) ResetPassword(ctx context.Context, token, password string) error {
	if strings.TrimSpace(token) == "" {
		return ErrPasswordResetInvalid
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.ConsumePasswordReset(ctx, tokenHash(token), string(hash), s.now())
}
