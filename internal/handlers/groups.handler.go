package handlers

import (
	"context"
	"net/http"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/middleware"
	"github.com/ICan-TC/users/internal/models"
	"github.com/ICan-TC/users/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
)

type GroupsHandler struct {
	svc *service.GroupsService
	log zerolog.Logger
}

func RegisterGroupsRoutes(api huma.API, svc *service.GroupsService) {
	h := &GroupsHandler{svc: svc, log: logging.L()}
	g := huma.NewGroup(api, "/groups")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Groups"}
	})
	g.UseMiddleware(middleware.AuthMiddleware)

	huma.Register(g, huma.Operation{
		OperationID:   "create-group",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create a group",
		Description:   "Create a group",
		DefaultStatus: http.StatusCreated,
	}, h.CreateGroup)

	huma.Register(g, huma.Operation{
		OperationID:   "update-group",
		Method:        http.MethodPatch,
		Path:          "",
		Summary:       "Update a group",
		Description:   "Update a group",
		DefaultStatus: http.StatusOK,
	}, h.UpdateGroup)

	huma.Register(g, huma.Operation{
		OperationID:   "get-group-by-id",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get a group by ID",
		Description:   "Get a group by ID",
		DefaultStatus: http.StatusOK,
	}, h.GetGroupByID)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-group",
		Method:        http.MethodDelete,
		Path:          "/{id}",
		Summary:       "Delete a group",
		Description:   "Delete a group",
		DefaultStatus: http.StatusOK,
	}, h.DeleteGroup)

	huma.Register(g, huma.Operation{
		OperationID:   "list-groups",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List groups",
		Description:   "List groups",
		DefaultStatus: http.StatusOK,
	}, h.ListGroups)
}

func (h *GroupsHandler) CreateGroup(c context.Context, input *dto.CreateGroupReq) (*dto.CreateGroupRes, error) {
	group, err := h.svc.CreateGroup(c, input.Body.Name, input.Body.Description, input.Body.TeacherID, input.Body.DefaultFee, input.Body.Subject, input.Body.Level,
		input.Body.Metadata)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("id", group.GroupID).Str("name", group.Name).Str("teacher_id", group.TeacherID).
		Str("subject", group.Subject).Str("level", group.Level).Any("created", group.CreatedAt).
		Msg("Created group")
	return &dto.CreateGroupRes{
		Body: *h.svc.ModelToRes(group),
	}, nil
}

func (h *GroupsHandler) UpdateGroup(c context.Context, input *dto.UpdateGroupReq) (*dto.UpdateGroupRes, error) {
	m := models.Groups{GroupID: input.Body.ID}
	if input.Body.Name != nil {
		m.Name = *input.Body.Name
	}
	if input.Body.Description != nil {
		m.Description = *input.Body.Description
	}
	if input.Body.TeacherID != nil {
		m.TeacherID = *input.Body.TeacherID
	}
	if input.Body.DefaultFee != nil {
		m.DefaultFee = *input.Body.DefaultFee
	}
	if input.Body.Subject != nil {
		m.Subject = *input.Body.Subject
	}
	if input.Body.Level != nil {
		m.Level = *input.Body.Level
	}
	if input.Body.Metadata != nil {
		m.Metadata = input.Body.Metadata
	}
	group, err := h.svc.UpdateGroup(c, m)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateGroupRes{
		Body: *h.svc.ModelToRes(group),
	}, nil
}

func (h *GroupsHandler) GetGroupByID(c context.Context, input *dto.GetGroupByIDReq) (*dto.GetGroupByIDRes, error) {
	group, err := h.svc.GetGroupByID(c, input.ID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("id", input.ID).Int("created_at", group.CreatedAt).
		Msg("Get group by ID")
	return &dto.GetGroupByIDRes{
		Body: *group,
	}, nil
}

func (h *GroupsHandler) DeleteGroup(c context.Context, input *dto.DeleteGroupReq) (*dto.DeleteGroupRes, error) {
	if err := h.svc.DeleteGroup(c, input.ID); err != nil {
		return nil, err
	}
	return &dto.DeleteGroupRes{
		Body: dto.DeleteGroupResBody{
			ID: input.ID,
		},
	}, nil
}

func (h *GroupsHandler) ListGroups(c context.Context, input *dto.ListGroupsReq) (*dto.ListGroupsRes, error) {
	return h.svc.GetGroups(c, input)
}
