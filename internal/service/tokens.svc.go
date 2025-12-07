package service

import (
	"context"
	"strings"
	"time"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/lib/tokens"
	"github.com/ICan-TC/users/internal/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type TokensService struct {
	db  *bun.DB
	tp  *tokens.TokenProvider
	log zerolog.Logger
}

func NewTokensService(tp *tokens.TokenProvider, db *bun.DB) *TokensService {
	logger := logging.L().With().Str("service", "tokens.svc").Logger()
	return &TokensService{
		tp:  tp,
		log: logger,
		db:  db,
	}
}

func (s *TokensService) TokensPair(ctx context.Context, sub string, username string, email string) (*tokens.TokensPair, error) {
	tokenID := ulid.Make().String()
	t, err := s.tp.GetTokensPair(ctx, sub, username, email, tokenID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get tokens pair")
		return nil, huma.Error500InternalServerError("could not create tokens")
	}
	m := models.RefreshTokens{
		ID:        tokenID,
		UserID:    sub,
		Token:     t.RefreshTokenID,
		Device:    "",
		ExpiresAt: t.RefreshExp,
		RevokedAt: nil,
	}
	if err := s.db.NewInsert().Model(&m).Returning("*").Scan(ctx, &m); err != nil {
		s.log.Error().Err(err).Msg("failed to insert refresh token")
		return nil, huma.Error500InternalServerError("could not save token")
	}

	return t, nil
}

func (s *TokensService) RefreshTokens(ctx context.Context, refreshToken string) (string, uint, error) {
	claims, err := s.tp.ParseRefresh(ctx, refreshToken)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to parse refresh token")
		return "", 0, huma.Error500InternalServerError("could not parse refresh token")
	}
	t := models.RefreshTokens{ID: claims.TokenID}
	err = s.db.NewSelect().Model(&t).WherePK("id").Scan(ctx, &t)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to select refresh token")
		return "", 0, huma.Error500InternalServerError("could not select refresh token")
	}
	if t.RevokedAt != nil {
		// TODO: log the user out of all sessions
		return "", 0, huma.Error400BadRequest("refresh token has been revoked")
	}

	aToken, aexp, err := s.tp.GetAccess(ctx, claims.Subject, claims.Username, claims.Email, claims.TokenID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get access token")
		return "", 0, huma.Error500InternalServerError("could not create access token")
	}
	return aToken.String(), uint(aexp.Unix()), nil
}

func (s *TokensService) ValidateAccessToken(ctx context.Context, token string) (*tokens.UserClaims, error) {
	claims, err := s.tp.ParseAccess(ctx, token)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, huma.Error400BadRequest(err.Error())
		} else if strings.Contains(err.Error(), "invalid token type") {
			s.log.Err(err).Msg(err.Error())
			return nil, huma.Error400BadRequest("invalid token type")
		}
		s.log.Error().Err(err).Msg("failed to parse access token")
		return nil, huma.Error500InternalServerError("could not parse access token")
	}

	t := models.RefreshTokens{ID: claims.TokenID}
	if err := s.db.NewSelect().Model(&t).WherePK("id").Scan(ctx, &t); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error401Unauthorized("invalid token")
		}
		return nil, huma.Error500InternalServerError("could not select refresh token")
	}
	if t.RevokedAt != nil {
		// TODO: log the user out of all sessions
		return nil, huma.Error401Unauthorized("refresh token has been revoked")
	}
	return claims, nil
}

func (s *TokensService) RevokeRefreshToken(ctx context.Context, token string) error {
	claims, err := s.tp.ParseRefresh(ctx, token)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to parse refresh token")
		return huma.Error500InternalServerError("could not parse refresh token")
	}
	revokedAt := time.Now()
	t := models.RefreshTokens{ID: claims.TokenID, RevokedAt: &revokedAt}
	if _, err := s.db.NewUpdate().Model(&t).OmitZero().WherePK("id").Exec(ctx); err != nil {
		s.log.Error().Err(err).Msg("failed to revoke refresh token")
		return huma.Error500InternalServerError("could not revoke refresh token")
	}

	return nil
}
