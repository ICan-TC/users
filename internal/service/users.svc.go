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
	var users []models.Users
	res := &dto.ListUsersRes{
		Body: dto.ListUsersResBody{
			Total:     0,
			ListQuery: dto.ListQueryRes{},
			Users:     nil,
		},
	}

	q := s.db.NewSelect().
		Model(&users).
		Relation("Teacher").
		Relation("Student").
		Relation("Employee").
		Relation("Parent")

	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where(
			"(username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR family_name ILIKE ?)",
			search, search, search, search,
		)
	}
	// if params.Filters != "" {
	// 	s.log.Warn().Str("filters", params.Filters).Msg("filters are not implemented yet")
	// }
	// if params.Includes != "" {
	// 	s.log.Warn().Str("includes", params.Includes).Msg("includes are not implemented yet")
	// }

	filters, err := dto.ParseFilters(params.Filters)
	if err != nil {
		s.log.Err(err).Msg("Couldn't parse filters")
	}
	q = dto.ApplyFilters(filters, q)

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	res.Body.Total = total

	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &users); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resUsers := []dto.UserModelRes{}
	for _, u := range users {
		resUsers = append(resUsers, *s.ModelToRes(&u, false))
	}
	res.Body.Users = resUsers
	res.Body.ListQuery = dto.ListQueryRes{
		Page: params.Page, PerPage: params.PerPage, SortBy: params.SortBy, SortDir: params.SortDir, Search: params.Search, Includes: params.Includes, Filters: filters,
	}
	return res, nil
}

func (s *UsersService) GetUserByID(ctx context.Context, id string) (*dto.UserModelRes, error) {
	m := models.Users{UserID: id}
	if err := s.db.NewSelect().Model(&m).WherePK("id").Scan(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't get user")
		return nil, huma.Error404NotFound("user not found")
	}
	return s.ModelToRes(&m, false), nil
}

func (s *UsersService) GetUserByField(ctx context.Context, f string, v string, include_hash bool) (*dto.UserModelRes, error) {
	m := models.Users{}
	q := s.db.NewSelect().Model(&m).Where(fmt.Sprintf("%s = ?", f), v)
	s.log.Debug().Str("query", q.String()).Msg("Couldn't get user")
	if err := q.Scan(ctx, &m); err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return s.ModelToRes(&m, include_hash), nil
}

func (s *UsersService) CreateUser(ctx context.Context, data *dto.CreateUserReqBody) (*dto.UserModelRes, error) {
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
	return s.ModelToRes(&m, false), nil
}

func (s *UsersService) UpdateUser(ctx context.Context, user models.Users) (*dto.UserModelRes, error) {
	m := user
	m.UserID = user.UserID
	if err := s.db.NewUpdate().Model(&m).Returning("*").OmitZero().WherePK("id").Scan(ctx, &m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("user not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return s.ModelToRes(&m, false), nil
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

func (s *UsersService) ModelToRes(m *models.Users, include_hash bool) *dto.UserModelRes {
	return UsersModelToRes(m, include_hash)
}

func UsersModelToRes(m *models.Users, include_hash bool) *dto.UserModelRes {
	if m == nil {
		return nil
	}
	res := &dto.UserModelRes{}
	res.ID = m.UserID
	res.Username = m.Username
	res.Email = m.Email

	if include_hash {
		res.PasswordHash = &m.PasswordHash
	}

	if m.FirstName != nil {
		res.FirstName = m.FirstName
	}
	if m.FamilyName != nil {
		res.FamilyName = m.FamilyName
	}
	if m.PhoneNumber != nil {
		res.PhoneNumber = m.PhoneNumber
	}
	if m.DateOfBirth != nil {
		res.DateOfBirth = m.DateOfBirth
	}
	if m.Student != nil {
		res.StudentID = &m.Student.StudentID
	}
	if m.Teacher != nil {
		res.TeacherID = &m.Teacher.TeacherID
	}
	if m.Employee != nil {
		res.EmployeeID = &m.Employee.EmployeeID
	}
	if m.Parent != nil {
		res.ParentID = &m.Parent.ParentID
	}
	if !m.CreatedAt.IsZero() {
		res.CreatedAt = int(m.CreatedAt.Unix())
	}
	if !m.UpdatedAt.IsZero() {
		res.UpdatedAt = int(m.UpdatedAt.Unix())
	}
	return res
}
