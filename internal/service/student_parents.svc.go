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

type StudentParentsService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewStudentParentsService(db *bun.DB) (*StudentParentsService, error) {
	log := logging.L().With().Str("service", "student_parents.svc").Logger()
	return &StudentParentsService{log: log, db: db}, nil
}

func (s *StudentParentsService) GetStudentParents(ctx context.Context, params *dto.ListStudentParentsReq) (*dto.ListStudentParentsRes, error) {
	total, err := s.db.NewSelect().Model((*models.StudentParents)(nil)).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	var studentParents []models.StudentParents
	res := &dto.ListStudentParentsRes{
		Body: dto.ListStudentParentsResBody{
			Total:          0,
			ListQuery:      params.ListQuery,
			StudentParents: nil,
		},
	}
	res.Body.Total = total

	q := s.db.NewSelect().Model(&studentParents)
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("(student_id ILIKE ? OR parent_id ILIKE ?)", search, search)
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &studentParents); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resStudentParents := []dto.GetStudentParentResBody{}
	for _, sp := range studentParents {
		newSP := dto.GetStudentParentResBody{
			StudentID: sp.StudentID,
			ParentID:  sp.ParentID,
			CreatedAt: int(sp.CreatedAt.Unix()),
			UpdatedAt: int(sp.UpdatedAt.Unix()),
		}
		resStudentParents = append(resStudentParents, newSP)
	}
	res.Body.StudentParents = resStudentParents
	return res, nil
}

func (s *StudentParentsService) GetStudentParentByID(ctx context.Context, studentID string, parentID string) (*models.StudentParents, error) {
	// Validate IDs
	if _, err := ulid.Parse(studentID); err != nil {
		return nil, huma.Error400BadRequest("studentID is invalid", err)
	}
	if _, err := ulid.Parse(parentID); err != nil {
		return nil, huma.Error400BadRequest("parentID is invalid", err)
	}

	m := models.StudentParents{StudentID: studentID, ParentID: parentID}
	if err := s.db.NewSelect().Model(&m).WherePK().Scan(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't get student-parent relationship")
		return nil, huma.Error404NotFound("student-parent relationship not found")
	}
	return &m, nil
}

func (s *StudentParentsService) CreateStudentParent(ctx context.Context, studentID string, parentID string) (*models.StudentParents, error) {
	// Validate IDs
	if _, err := ulid.Parse(studentID); err != nil {
		return nil, huma.Error400BadRequest("studentID is invalid", err)
	}
	if _, err := ulid.Parse(parentID); err != nil {
		return nil, huma.Error400BadRequest("parentID is invalid", err)
	}

	// Verify student exists
	student := models.Students{StudentID: studentID}
	if err := s.db.NewSelect().Model(&student).WherePK("id").Scan(ctx, &student); err != nil {
		return nil, huma.Error404NotFound("student not found")
	}

	// Verify parent exists
	parent := models.Parents{ParentID: parentID}
	if err := s.db.NewSelect().Model(&parent).WherePK("id").Scan(ctx, &parent); err != nil {
		return nil, huma.Error404NotFound("parent not found")
	}

	m := models.StudentParents{
		StudentID: studentID,
		ParentID:  parentID,
	}
	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't insert student-parent relationship")
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, huma.Error409Conflict("student-parent relationship already exists")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *StudentParentsService) UpdateStudentParent(ctx context.Context, oldStudentID string, oldParentID string, newStudentID string, newParentID string) (*models.StudentParents, error) {
	// Validate all IDs
	if _, err := ulid.Parse(oldStudentID); err != nil {
		return nil, huma.Error400BadRequest("old studentID is invalid", err)
	}
	if _, err := ulid.Parse(oldParentID); err != nil {
		return nil, huma.Error400BadRequest("old parentID is invalid", err)
	}
	if _, err := ulid.Parse(newStudentID); err != nil {
		return nil, huma.Error400BadRequest("new studentID is invalid", err)
	}
	if _, err := ulid.Parse(newParentID); err != nil {
		return nil, huma.Error400BadRequest("new parentID is invalid", err)
	}

	// Verify new student exists
	student := models.Students{StudentID: newStudentID}
	if err := s.db.NewSelect().Model(&student).WherePK("id").Scan(ctx, &student); err != nil {
		return nil, huma.Error404NotFound("new student not found")
	}

	// Verify new parent exists
	parent := models.Parents{ParentID: newParentID}
	if err := s.db.NewSelect().Model(&parent).WherePK("id").Scan(ctx, &parent); err != nil {
		return nil, huma.Error404NotFound("new parent not found")
	}

	m := models.StudentParents{
		StudentID: newStudentID,
		ParentID:  newParentID,
	}

	// Delete old relationship and insert new one
	if _, err := s.db.NewDelete().Model(&models.StudentParents{StudentID: oldStudentID, ParentID: oldParentID}).WherePK().Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("student-parent relationship not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't update student-parent relationship")
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &m, nil
}

func (s *StudentParentsService) DeleteStudentParent(ctx context.Context, studentID string, parentID string) error {
	// Validate IDs
	if _, err := ulid.Parse(studentID); err != nil {
		return huma.Error400BadRequest("studentID is invalid", err)
	}
	if _, err := ulid.Parse(parentID); err != nil {
		return huma.Error400BadRequest("parentID is invalid", err)
	}

	m := models.StudentParents{StudentID: studentID, ParentID: parentID}
	if _, err := s.db.NewDelete().Model(&m).WherePK().Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return huma.Error404NotFound("student-parent relationship not found")
		}
		return huma.Error500InternalServerError(err.Error())
	}
	return nil
}
