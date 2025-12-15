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

type GroupsService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewGroupsService(db *bun.DB) (*GroupsService, error) {
	log := logging.L().With().Str("service", "groups.svc").Logger()
	return &GroupsService{log: log, db: db}, nil
}

func (s *GroupsService) GetGroups(ctx context.Context, params *dto.ListGroupsReq) (*dto.ListGroupsRes, error) {
	total, err := s.db.NewSelect().Model((*models.Groups)(nil)).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	var groups []models.Groups
	res := &dto.ListGroupsRes{
		Body: dto.ListGroupsResBody{
			Total:     0,
			ListQuery: params.ListQuery,
			Groups:    nil,
		},
	}
	res.Body.Total = total

	q := s.db.NewSelect().Model(&groups)
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("(name ILIKE ? OR subject ILIKE ? OR level ILIKE ?)", search, search, search)
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &groups); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resGroups := []dto.GroupModelRes{}
	for _, grp := range groups {
		resGroups = append(resGroups, *s.ModelToRes(&grp))
	}
	res.Body.Groups = resGroups
	return res, nil
}

func (s *GroupsService) GetGroupByID(ctx context.Context, id string) (*dto.GroupModelRes, error) {
	m := models.Groups{GroupID: id}
	if err := s.db.NewSelect().Model(&m).WherePK("id").Scan(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't get group")
		return nil, huma.Error404NotFound("group not found")
	}
	return s.ModelToRes(&m), nil
}

func (s *GroupsService) CreateGroup(ctx context.Context, name string, description string, teacherID string, defaultFee float64, subject string, level string, metadata map[string]interface{}) (*models.Groups, error) {
	if _, err := ulid.Parse(teacherID); err != nil {
		return nil, huma.Error400BadRequest("teacherID is invalid", err)
	}

	m := models.Groups{
		GroupID:     ulid.Make().String(),
		Name:        name,
		Description: description,
		TeacherID:   teacherID,
		DefaultFee:  defaultFee,
		Subject:     subject,
		Level:       level,
		Metadata:    metadata,
	}
	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't insert group")
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *GroupsService) UpdateGroup(ctx context.Context, group models.Groups) (*models.Groups, error) {
	m := group
	m.GroupID = group.GroupID
	if err := s.db.NewUpdate().Model(&m).Returning("*").OmitZero().WherePK("id").Scan(ctx, &m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("group not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *GroupsService) DeleteGroup(ctx context.Context, id string) error {
	m := models.Groups{GroupID: id}
	if _, err := s.db.NewDelete().Model(&m).WherePK("id").Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return huma.Error404NotFound("group not found")
		}
		return huma.Error500InternalServerError(err.Error())
	}
	return nil
}

func (s *GroupsService) ModelToRes(m *models.Groups) *dto.GroupModelRes {
	if m == nil {
		return nil
	}
	res := &dto.GroupModelRes{
		ID:          m.GroupID,
		Name:        m.Name,
		Description: m.Description,
		TeacherID:   m.TeacherID,
		DefaultFee:  m.DefaultFee,
		Subject:     m.Subject,
		Level:       m.Level,
		Metadata:    m.Metadata,
	}
	if !m.CreatedAt.IsZero() {
		res.CreatedAt = int(m.CreatedAt.Unix())
	}
	if !m.UpdatedAt.IsZero() {
		res.UpdatedAt = int(m.UpdatedAt.Unix())
	}
	return res
}
