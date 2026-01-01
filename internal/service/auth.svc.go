package service

import (
	"context"
	"strings"
	"time"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/lib/tokens"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log  zerolog.Logger
	usvc *UsersService
	tsvc *TokensService
}

func NewAuthService(usvc *UsersService, tsvc *TokensService) (*AuthService, error) {
	log := logging.L().With().Str("service", "auth.svc").Logger()
	return &AuthService{log: log, usvc: usvc, tsvc: tsvc}, nil
}

type Tokens struct {
	Access     string    `json:"access"`
	AccessExp  time.Time `json:"access_exp"`
	Refresh    string    `json:"refresh"`
	RefreshExp time.Time `json:"refresh_exp"`
}

func (s *AuthService) Signup(ctx context.Context, email, username, password string) (*models.Users, *tokens.TokensPair, error) {
	u, err := s.usvc.CreateUser(ctx, &dto.CreateUserReqBody{Email: email, Username: username, Password: password})
	if err != nil {
		return nil, nil, err
	}
	t, err := s.tsvc.TokensPair(ctx, u.UserID, u.Username, u.Email)
	if err != nil {
		s.log.Err(err).Msg("could not create tokens")
		return nil, nil, err
	}
	return u, t, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.Users, *tokens.TokensPair, error) {
	u, err := s.usvc.GetUserByField(ctx, "username", username)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil, huma.Error404NotFound("user not found")
		}
		return nil, nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, nil, huma.Error401Unauthorized("invalid credentials")
	}
	t, err := s.tsvc.TokensPair(ctx, u.UserID, u.Username, u.Email)
	if err != nil {
		s.log.Err(err).Msg("could not create tokens")
		return nil, nil, huma.Error500InternalServerError(err.Error())
	}
	return u, t, nil
}

func (s *AuthService) Refresh(ctx context.Context, token string) (string, uint, error) {
	return s.tsvc.RefreshTokens(ctx, token)
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.tsvc.RevokeRefreshToken(ctx, token)
}

func (s *AuthService) Revoke(ctx context.Context, token string) error {
	return s.tsvc.RevokeRefreshToken(ctx, token)
}

func (s *AuthService) Verify(ctx context.Context, token string) (*tokens.UserClaims, error) {
	return s.tsvc.ValidateAccessToken(ctx, token)
}
