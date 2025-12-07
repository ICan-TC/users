package service

import (
	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/lib/tokens"
	"github.com/ICan-TC/users/internal/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
)

type AuthService struct {
	tp  *tokens.TokenProvider
	log zerolog.Logger
}

func NewAuthService() (*AuthService, error) {
	log := logging.L().With().Str("service", "users.svc").Logger()
	cfg := config.Get()
	tp, err := tokens.NewTokenProvider(tokens.TokenProviderArgs{
		Secret:          cfg.Auth.Secret,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	})
	if err != nil {
		return nil, err
	}
	return &AuthService{tp: tp, log: log}, nil
}

func (s *AuthService) Signup(username, password string) (string, error) {
	return "", huma.Error501NotImplemented("not implemented")
}

func (s *AuthService) Login(username, password string) (string, error) {
	return "", huma.Error501NotImplemented("not implemented")
}

func (s *AuthService) Refresh(token string) (string, error) {
	return "", huma.Error501NotImplemented("not implemented")
}

func (s *AuthService) Logout(token string) error {
	return huma.Error501NotImplemented("not implemented")
}

func (s *AuthService) Verify(token string) (string, error) {
	return "", huma.Error501NotImplemented("not implemented")
}

func (s *AuthService) Revoke(token string) error {
	return huma.Error501NotImplemented("not implemented")
}
