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

type EnrollmentsService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewEnrollmentsService(db *bun.DB) (*EnrollmentsService, error) {
	log := logging.L().With().Str("service", "enrollments.svc").Logger()
	return &EnrollmentsService{log: log, db: db}, nil
}

func (s *EnrollmentsService) GetEnrollments(ctx context.Context, params *dto.ListEnrollmentsReq) (*dto.ListEnrollmentsRes, error) {
	total, err := s.db.NewSelect().Model((*models.Enrollments)(nil)).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	var enrollments []models.Enrollments
	res := &dto.ListEnrollmentsRes{
		Body: dto.ListEnrollmentsResBody{
			Total:       0,
			ListQuery:   params.ListQuery,
			Enrollments: nil,
		},
	}
	res.Body.Total = total

	q := s.db.NewSelect().Model(&enrollments)
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("(student_id ILIKE ? OR group_id ILIKE ?)", search, search)
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &enrollments); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resEnrollments := []dto.EnrollmentModelRes{}
	for _, enr := range enrollments {
		resEnrollments = append(resEnrollments, *s.ModelToRes(&enr))
	}
	res.Body.Enrollments = resEnrollments
	return res, nil
}

func (s *EnrollmentsService) GetEnrollmentByID(ctx context.Context, studentID, groupID string) (*dto.EnrollmentModelRes, error) {
	if _, err := ulid.Parse(studentID); err != nil {
		return nil, huma.Error400BadRequest("studentID is invalid", err)
	}
	if _, err := ulid.Parse(groupID); err != nil {
		return nil, huma.Error400BadRequest("groupID is invalid", err)
	}

	m := models.Enrollments{StudentID: studentID, GroupID: groupID}
	if err := s.db.NewSelect().Model(&m).WherePK().Scan(ctx); err != nil {
		s.log.Err(err).Msg("Couldn't get enrollment")
		return nil, huma.Error404NotFound("enrollment not found")
	}
	return s.ModelToRes(&m), nil
}

func (s *EnrollmentsService) CreateEnrollment(ctx context.Context, studentID string, groupID string, fee *float64) (*models.Enrollments, error) {
	if _, err := ulid.Parse(studentID); err != nil {
		return nil, huma.Error400BadRequest("studentID is invalid", err)
	}
	if _, err := ulid.Parse(groupID); err != nil {
		return nil, huma.Error400BadRequest("groupID is invalid", err)
	}

	// If fee not specified, fetch group's default fee
	actualFee := fee
	if fee == nil {
		group := models.Groups{GroupID: groupID}
		if err := s.db.NewSelect().Model(&group).WherePK("id").Scan(ctx); err != nil {
			s.log.Err(err).Msg("Couldn't get group for default fee")
			return nil, huma.Error404NotFound("group not found")
		}
		actualFee = &group.DefaultFee
	}

	m := models.Enrollments{
		StudentID: studentID,
		GroupID:   groupID,
		Fee:       *actualFee,
	}
	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't insert enrollment")
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *EnrollmentsService) UpdateEnrollment(ctx context.Context, enrollment models.Enrollments) (*models.Enrollments, error) {
	m := enrollment
	if err := s.db.NewUpdate().Model(&m).Returning("*").OmitZero().WherePK().Scan(ctx, &m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("enrollment not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *EnrollmentsService) DeleteEnrollment(ctx context.Context, studentID, groupID string) error {
	if _, err := ulid.Parse(studentID); err != nil {
		return huma.Error400BadRequest("studentID is invalid", err)
	}
	if _, err := ulid.Parse(groupID); err != nil {
		return huma.Error400BadRequest("groupID is invalid", err)
	}

	m := models.Enrollments{StudentID: studentID, GroupID: groupID}
	if _, err := s.db.NewDelete().Model(&m).WherePK().Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return huma.Error404NotFound("enrollment not found")
		}
		return huma.Error500InternalServerError(err.Error())
	}
	return nil
}

func (s *EnrollmentsService) ModelToRes(m *models.Enrollments) *dto.EnrollmentModelRes {
	if m == nil {
		return nil
	}
	res := &dto.EnrollmentModelRes{
		StudentID: m.StudentID,
		GroupID:   m.GroupID,
		Fee:       m.Fee,
	}
	if !m.CreatedAt.IsZero() {
		res.CreatedAt = int(m.CreatedAt.Unix())
	}
	if !m.UpdatedAt.IsZero() {
		res.UpdatedAt = int(m.UpdatedAt.Unix())
	}
	return res
}

func (s *EnrollmentsService) GetEnrollmentsByGroupID(ctx context.Context, params *dto.GetEnrollmentsByGroupIDReq) (*dto.GetEnrollmentsByGroupIDRes, error) {
	// Validate groupID
	if _, err := ulid.Parse(params.GroupID); err != nil {
		return nil, huma.Error400BadRequest("groupID is invalid", err)
	}

	total, err := s.db.NewSelect().Model((*models.Enrollments)(nil)).Where("group_id = ?", params.GroupID).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	var enrollments []models.Enrollments
	res := &dto.GetEnrollmentsByGroupIDRes{
		Body: dto.ListEnrollmentsResBody{
			Total:       total,
			ListQuery:   params.ListQuery,
			Enrollments: nil,
		},
	}

	q := s.db.NewSelect().Model(&enrollments).Where("group_id = ?", params.GroupID)
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("student_id ILIKE ?", search)
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &enrollments); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resEnrollments := []dto.EnrollmentModelRes{}
	for _, enr := range enrollments {
		resEnrollments = append(resEnrollments, *s.ModelToRes(&enr))
	}
	res.Body.Enrollments = resEnrollments
	return res, nil
}

func (s *EnrollmentsService) GetEnrollmentsByStudentID(ctx context.Context, params *dto.GetEnrollmentsByStudentIDReq) (*dto.GetEnrollmentsByStudentIDRes, error) {
	// Validate studentID
	if _, err := ulid.Parse(params.StudentID); err != nil {
		return nil, huma.Error400BadRequest("studentID is invalid", err)
	}

	total, err := s.db.NewSelect().Model((*models.Enrollments)(nil)).Where("student_id = ?", params.StudentID).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	var enrollments []models.Enrollments
	res := &dto.GetEnrollmentsByStudentIDRes{
		Body: dto.ListEnrollmentsResBody{
			Total:       total,
			ListQuery:   params.ListQuery,
			Enrollments: nil,
		},
	}

	q := s.db.NewSelect().Model(&enrollments).Where("student_id = ?", params.StudentID)
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("group_id ILIKE ?", search)
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &enrollments); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resEnrollments := []dto.EnrollmentModelRes{}
	for _, enr := range enrollments {
		resEnrollments = append(resEnrollments, *s.ModelToRes(&enr))
	}
	res.Body.Enrollments = resEnrollments
	return res, nil
}
