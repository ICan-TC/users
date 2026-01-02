package service

import (
	"context"
	"strings"
	"time"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type TeachersService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewTeachersService(db *bun.DB) (*TeachersService, error) {
	log := logging.L().With().Str("service", "teachers.svc").Logger()
	return &TeachersService{log: log, db: db}, nil
}

func (s *TeachersService) GetTeachers(ctx context.Context, params *dto.ListTeachersReq) (*dto.ListTeachersRes, error) {
	total, err := s.db.NewSelect().Model((*models.Teachers)(nil)).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	var teachers []models.Teachers
	res := &dto.ListTeachersRes{
		Body: dto.ListTeachersResBody{
			Total:     0,
			ListQuery: params.ListQuery,
			Teachers:  nil,
		},
	}
	res.Body.Total = total

	q := s.db.NewSelect().Model(&teachers).Relation("User")
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("user_id ILIKE ?", search)
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &teachers); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resTeachers := []dto.TeachersModelRes{}
	for _, tch := range teachers {
		newTeacher := *s.ModelToRes(&tch)
		resTeachers = append(resTeachers, newTeacher)
	}
	res.Body.Teachers = resTeachers
	return res, nil
}

func (s *TeachersService) GetTeacherByID(ctx context.Context, id string) (*dto.TeachersModelRes, error) {
	m := models.Teachers{TeacherID: id}
	if err := s.db.NewSelect().Model(&m).Relation("User").WherePK("id").Scan(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't get teacher")
		return nil, huma.Error404NotFound("teacher not found")
	}
	return s.ModelToRes(&m), nil
}

func (s *TeachersService) CreateTeacher(ctx context.Context, userID string) (*models.Teachers, error) {
	if _, err := ulid.Parse(userID); err != nil {
		return nil, huma.Error400BadRequest("userID is invalid", err)
	}
	m := models.Teachers{
		TeacherID: userID,
		UserID:    &userID,
	}
	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't insert teacher")
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *TeachersService) UpdateTeacher(ctx context.Context, teacher models.Teachers) (*models.Teachers, error) {
	m := teacher
	m.TeacherID = teacher.TeacherID
	m.UpdatedAt = time.Now()
	if err := s.db.NewUpdate().Model(&m).Returning("*").OmitZero().WherePK("id").Scan(ctx, &m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("teacher not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *TeachersService) DeleteTeacher(ctx context.Context, id string) error {
	m := models.Teachers{TeacherID: id, DeletedAt: time.Now()}
	if _, err := s.db.NewDelete().Model(&m).WherePK("id").Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return huma.Error404NotFound("teacher not found")
		}
		return huma.Error500InternalServerError(err.Error())
	}
	return nil
}

func (s *TeachersService) ModelToRes(m *models.Teachers) *dto.TeachersModelRes {
	if m == nil {
		return nil
	}
	user := m.User
	var userRes dto.UserModelRes
	if user != nil {
		userRes = dto.UserModelRes{
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			FamilyName:  user.FamilyName,
			PhoneNumber: user.PhoneNumber,
			DateOfBirth: user.DateOfBirth,
		}
		if !user.CreatedAt.IsZero() {
			userRes.CreatedAt = int(user.CreatedAt.Unix())
		}
		if !user.UpdatedAt.IsZero() {
			userRes.UpdatedAt = int(user.UpdatedAt.Unix())
		}
	}
	res := &dto.TeachersModelRes{
		ID:           m.TeacherID,
		UserID:       *m.UserID,
		UserModelRes: userRes,
	}
	if !m.CreatedAt.IsZero() {
		res.CreatedAt = int(m.CreatedAt.Unix())
	}
	if !m.UpdatedAt.IsZero() {
		res.UpdatedAt = int(m.UpdatedAt.Unix())
	}
	return res
}
