package service

import (
	"context"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
)

type UsersService struct {
	log zerolog.Logger
}

func NewUsersService() (*UsersService, error) {
	log := logging.L().With().Str("service", "users.svc").Logger()
	return &UsersService{log: log}, nil
}

func (s *UsersService) GetUsers(ctx context.Context) ([]models.Users, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (s *UsersService) GetUser(ctx context.Context, id string) (models.Users, error) {
	return models.Users{}, huma.Error501NotImplemented("not implemented")
}

func (s *UsersService) CreateUser(ctx context.Context, user models.Users) (models.Users, error) {
	return models.Users{}, huma.Error501NotImplemented("not implemented")
}

func (s *UsersService) UpdateUser(ctx context.Context, user models.Users) (models.Users, error) {
	return models.Users{}, huma.Error501NotImplemented("not implemented")
}

func (s *UsersService) DeleteUser(ctx context.Context, id string) error {
	return huma.Error501NotImplemented("not implemented")
}
