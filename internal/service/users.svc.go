package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
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

func (s *UsersService) GetUsers(ctx context.Context, params *dto.ListUsersReq) (*dto.ListUsersRes, error) {
	total, err := s.db.NewSelect().Model((*models.Users)(nil)).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	var users []models.Users
	res := &dto.ListUsersRes{
		Body: dto.ListUsersResBody{
			Total:     0,
			ListQuery: params.ListQuery,
			Users:     nil,
		},
	}
	res.Body.Total = total

	q := s.db.NewSelect().Model(&users)
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where(
			"(username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR family_name ILIKE ?)",
			search, search, search, search,
		)
	}
	if params.Filters != "" {
		s.log.Warn().Str("filters", params.Filters).Msg("filters are not implemented yet")
	}
	if params.Includes != "" {
		s.log.Warn().Str("includes", params.Includes).Msg("includes are not implemented yet")
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &users); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resUsers := []dto.GetUserResBody{}
	for _, u := range users {
		newUser := dto.GetUserResBody{
			ID:       u.UserID,
			Username: u.Username,
			Email:    u.Email,
		}
		if u.FirstName != nil {
			newUser.FirstName = u.FirstName
		}
		if u.FamilyName != nil {
			newUser.FamilyName = u.FamilyName
		}
		if u.PhoneNumber != nil {
			newUser.PhoneNumber = u.PhoneNumber
		}
		if u.DateOfBirth != nil {
			newUser.DateOfBirth = u.DateOfBirth
		}
		newUser.CreatedAt = int(u.CreatedAt.Unix())
		newUser.UpdatedAt = int(u.UpdatedAt.Unix())
		resUsers = append(resUsers, newUser)

	}
	res.Body.Users = resUsers
	return res, nil
}

func (s *UsersService) GetUserByID(ctx context.Context, id string) (*models.Users, error) {
	m := models.Users{UserID: id}
	if err := s.db.NewSelect().Model(&m).WherePK("id").Scan(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't get user")
		return nil, huma.Error404NotFound("user not found")
	}
	return &m, nil
}

func (s *UsersService) GetUserByField(ctx context.Context, f string, v string) (*models.Users, error) {
	m := models.Users{}
	q := s.db.NewSelect().Model(&m).Where(fmt.Sprintf("%s = ?", f), v)
	s.log.Debug().Str("query", q.String()).Msg("Couldn't get user")
	if err := q.Scan(ctx, &m); err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *UsersService) CreateUser(ctx context.Context, data *dto.CreateUserReqBody) (*models.Users, error) {
	// TODO: fix this, should probably create a new struct for this function's input
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil
	}
	userID := ulid.Make().String()
	m := models.Users{
		Username:     data.Username,
		Email:        data.Email,
		PasswordHash: string(hash),
		UserID:       userID,
	}
	if data.DateOfBirth != nil {
		parsedTime, err := time.Parse(time.RFC3339, *data.DateOfBirth)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid date format", err)
		}
		m.DateOfBirth = &parsedTime
	}
	if data.FirstName != nil {
		m.FirstName = data.FirstName
	}
	if data.FamilyName != nil {
		m.FamilyName = data.FamilyName
	}
	if data.PhoneNumber != nil {
		m.PhoneNumber = data.PhoneNumber
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
	m := user
	m.UserID = user.UserID
	if err := s.db.NewUpdate().Model(&m).Returning("*").OmitZero().WherePK("id").Scan(ctx, &m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("user not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *UsersService) DeleteUser(ctx context.Context, id string) error {
	m := models.Users{UserID: id}
	if _, err := s.db.NewDelete().Model(&m).WherePK("id").Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return huma.Error404NotFound("user not found")
		}
		return huma.Error500InternalServerError(err.Error())
	}
	return nil
}
