package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/rupi/timetracking/internal/domain"
)

type Service struct {
	db     *sqlx.DB
	jwtSvc *JWTService
}

func NewService(db *sqlx.DB, jwtSvc *JWTService) *Service {
	return &Service{db: db, jwtSvc: jwtSvc}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*domain.User, *TokenPair, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, nil, domain.ErrInvalidInput
	}
	if req.Timezone == "" {
		req.Timezone = "UTC"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{}
	err = s.db.QueryRowxContext(ctx,
		`INSERT INTO users (email, password_hash, name, timezone)
         VALUES ($1, $2, $3, $4)
         RETURNING *`,
		req.Email, string(hash), req.Name, req.Timezone,
	).StructScan(user)
	if err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokenPair(ctx, user.ID, user.Email)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*domain.User, *TokenPair, error) {
	user := &domain.User{}
	err := s.db.QueryRowxContext(ctx,
		`SELECT * FROM users WHERE email = $1`, req.Email,
	).StructScan(user)
	if err != nil {
		return nil, nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, domain.ErrUnauthorized
	}

	tokens, err := s.generateTokenPair(ctx, user.ID, user.Email)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	hash := hashToken(refreshToken)
	var userID uuid.UUID
	err := s.db.QueryRowContext(ctx,
		`DELETE FROM refresh_tokens WHERE token_hash = $1 AND expires_at > NOW() RETURNING user_id`,
		hash,
	).Scan(&userID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	user := &domain.User{}
	if err := s.db.QueryRowxContext(ctx, `SELECT * FROM users WHERE id = $1`, userID).StructScan(user); err != nil {
		return nil, domain.ErrNotFound
	}

	return s.generateTokenPair(ctx, user.ID, user.Email)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	hash := hashToken(refreshToken)
	_, err := s.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, hash)
	return err
}

func (s *Service) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	err := s.db.QueryRowxContext(ctx, `SELECT * FROM users WHERE id = $1`, userID).StructScan(user)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (s *Service) generateTokenPair(ctx context.Context, userID uuid.UUID, email string) (*TokenPair, error) {
	accessToken, err := s.jwtSvc.Generate(userID, email)
	if err != nil {
		return nil, err
	}

	rawRefresh := make([]byte, 32)
	if _, err := rand.Read(rawRefresh); err != nil {
		return nil, err
	}
	refreshToken := hex.EncodeToString(rawRefresh)
	hash := hashToken(refreshToken)

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, hash, time.Now().Add(7*24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
