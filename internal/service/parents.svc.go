package service

import (
	"context"
	"strings"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type ParentsService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewParentsService(db *bun.DB) (*ParentsService, error) {
	log := logging.L().With().Str("service", "parents.svc").Logger()
	db.RegisterModel(&models.StudentParents{})
	return &ParentsService{log: log, db: db}, nil
}

func (s *ParentsService) GetParents(ctx context.Context, params *dto.ListParentsReq) (*dto.ListParentsRes, error) {
	total, err := s.db.NewSelect().Model((*models.Parents)(nil)).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	var parents []models.Parents
	res := &dto.ListParentsRes{
		Body: dto.ListParentsResBody{
			Total:     0,
			ListQuery: params.ListQuery,
			Parents:   nil,
		},
	}
	res.Body.Total = total

	q := s.db.NewSelect().Model(&parents).Relation("User")
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("user_id ILIKE ?", search)
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &parents); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resParents := []dto.ParentModelRes{}
	for _, par := range parents {
		newParent := dto.ParentModelRes{
			ID:     par.ParentID,
			UserID: par.UserID,
			UserModelRes: dto.UserModelRes{
				ID:          par.User.UserID,
				Username:    par.User.Username,
				Email:       par.User.Email,
				FirstName:   par.User.FirstName,
				FamilyName:  par.User.FamilyName,
				PhoneNumber: par.User.PhoneNumber,
				DateOfBirth: par.User.DateOfBirth,
			},
			CreatedAt: int(par.CreatedAt.Unix()),
			UpdatedAt: int(par.UpdatedAt.Unix()),
		}
		resParents = append(resParents, newParent)
	}
	res.Body.Parents = resParents
	return res, nil
}

func (s *ParentsService) GetParentByID(ctx context.Context, id string) (*dto.ParentModelRes, error) {
	m := models.Parents{ParentID: id}
	if err := s.db.NewSelect().Model(&m).Relation("User").Relation("Students").WherePK("id").Scan(ctx); err != nil {
		s.log.Err(err).Msg("Couldn't get parent")
		return nil, huma.Error404NotFound("parent not found")
	}

	// Convert students to DTO
	students := []dto.GetStudentResBody{}
	if m.Students != nil {
		for _, student := range m.Students {
			newStudent := dto.GetStudentResBody{
				ID:        student.StudentID,
				Level:     student.Level,
				UserID:    student.UserID,
				CreatedAt: int(student.CreatedAt.Unix()),
				UpdatedAt: int(student.UpdatedAt.Unix()),
			}
			students = append(students, newStudent)
		}
	}

	res := s.ModelToRes(&m)
	res.Students = students

	return res, nil
}

func (s *ParentsService) CreateParent(ctx context.Context, userID string) (*models.Parents, error) {
	if _, err := ulid.Parse(userID); err != nil {
		return nil, huma.Error400BadRequest("userID is invalid", err)
	}
	m := models.Parents{
		ParentID: userID,
		UserID:   userID,
	}
	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't insert parent")
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *ParentsService) UpdateParent(ctx context.Context, parent models.Parents) (*models.Parents, error) {
	m := parent
	m.ParentID = parent.ParentID
	if err := s.db.NewUpdate().Model(&m).Returning("*").OmitZero().WherePK("id").Scan(ctx, &m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("parent not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *ParentsService) DeleteParent(ctx context.Context, id string) error {
	m := models.Parents{ParentID: id}
	if _, err := s.db.NewDelete().Model(&m).WherePK("id").Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return huma.Error404NotFound("parent not found")
		}
		return huma.Error500InternalServerError(err.Error())
	}
	return nil
}

func (s *ParentsService) ModelToRes(m *models.Parents) *dto.ParentModelRes {
	if m == nil {
		return nil
	}
	user := m.User
	var userRes dto.UserModelRes
	if user != nil {
		userRes = *UsersModelToRes(user, false)
	}
	res := &dto.ParentModelRes{
		ID:           m.ParentID,
		UserID:       m.UserID,
		Students:     nil,
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
