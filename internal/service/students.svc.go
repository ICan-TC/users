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

type StudentsService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewStudentsService(db *bun.DB) (*StudentsService, error) {
	log := logging.L().With().Str("service", "students.svc").Logger()
	return &StudentsService{log: log, db: db}, nil
}

func (s *StudentsService) GetStudents(ctx context.Context, params *dto.ListStudentsReq) (*dto.ListStudentsRes, error) {
	var students []models.Students
	res := &dto.ListStudentsRes{
		Body: dto.ListStudentsResBody{
			Total:     0,
			ListQuery: params.ListQuery,
			Students:  nil,
		},
	}
	q := s.db.NewSelect().Model(&students).Relation("User")
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("level ILIKE ?", search)
	}

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

	if err := q.Scan(ctx, &students); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resStudents := []dto.StudentsModelRes{}
	for _, st := range students {
		newStudent := *s.ModelToRes(&st)
		resStudents = append(resStudents, newStudent)
	}
	res.Body.Students = resStudents
	return res, nil
}

func (s *StudentsService) GetStudentByID(ctx context.Context, id string) (*dto.StudentsModelRes, error) {
	m := models.Students{StudentID: id}
	if err := s.db.NewSelect().Model(&m).Relation("User").WherePK("id").Scan(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't get student")
		return nil, huma.Error404NotFound("student not found")
	}
	return s.ModelToRes(&m), nil
}

func (s *StudentsService) CreateStudent(ctx context.Context, level string, userID *string) (*models.Students, error) {
	if userID == nil {
		return nil, huma.Error400BadRequest("userID is required")
	}
	if _, err := ulid.Parse(*userID); err != nil {
		return nil, huma.Error400BadRequest("userID is invalid", err)
	}
	m := models.Students{
		StudentID: *userID,
		Level:     &level,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if userID != nil {
		m.UserID = userID
	}
	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't insert student")
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *StudentsService) UpdateStudent(ctx context.Context, student models.Students) (*models.Students, error) {
	m := student
	m.StudentID = student.StudentID
	if err := s.db.NewUpdate().Model(&m).Returning("*").OmitZero().WherePK("id").Scan(ctx, &m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("student not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *StudentsService) DeleteStudent(ctx context.Context, id string) error {
	m := models.Students{StudentID: id, DeletedAt: time.Now()}
	if _, err := s.db.NewDelete().Model(&m).WherePK("id").Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return huma.Error404NotFound("student not found")
		}
		return huma.Error500InternalServerError(err.Error())
	}
	return nil
}

func (s *StudentsService) ModelToRes(m *models.Students) *dto.StudentsModelRes {
	if m == nil {
		return nil
	}
	user := m.User
	var userRes dto.UserModelRes
	if user != nil {
		userRes = *UsersModelToRes(user, false)
	}
	res := &dto.StudentsModelRes{
		ID:           m.StudentID,
		Level:        m.Level,
		UserID:       m.UserID,
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
