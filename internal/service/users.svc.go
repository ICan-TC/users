package service

import (
	"context"
	"strings"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewUsersService(db *bun.DB) (*UsersService, error) {
	log := logging.L().With().Str("service", "users.svc").Logger()
	return &UsersService{log: log, db: db}, nil
}

func (s *UsersService) GetUsers(ctx context.Context) ([]models.Users, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (s *UsersService) GetUserByField(ctx context.Context, f string, v string) (*models.Users, error) {
	m := models.Users{}
	if err := s.db.NewSelect().Model(&m).Where(f+" = ?", v).Scan(ctx, &m); err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *UsersService) GetUser(ctx context.Context, id string) (*models.Users, error) {
	m := models.Users{ID: id}
	s.log.Info().Str("id", id).Msg("user ID")
	if err := s.db.NewSelect().Model(&m).WherePK("id").Scan(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't get user")
		return nil, huma.Error404NotFound("user not found")
	}
	return &m, nil
}

func (s *UsersService) CreateUser(ctx context.Context, email, username, password string) (*models.Users, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil
	}
	m := models.Users{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		ID:           ulid.Make().String(),
	}
	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't insert user")
		if strings.Contains(err.Error(), "users_username_key") {
			return nil, huma.Error400BadRequest("username already exists")
		}
		if strings.Contains(err.Error(), "users_email_key") {
			return nil, huma.Error400BadRequest("email already exists")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *UsersService) UpdateUser(ctx context.Context, user models.Users) (*models.Users, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (s *UsersService) DeleteUser(ctx context.Context, id string) error {
	return huma.Error501NotImplemented("not implemented")
}
