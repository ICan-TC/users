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

type ParentsHandler struct {
	svc *service.ParentsService
	log zerolog.Logger
}

func RegisterParentsRoutes(api huma.API, svc *service.ParentsService) {
	h := &ParentsHandler{svc: svc, log: logging.L()}
	g := huma.NewGroup(api, "/parents")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Parents"}
	})
	g.UseMiddleware(middleware.AuthMiddleware)

	huma.Register(g, huma.Operation{
		OperationID:   "create-parent",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create a parent",
		Description:   "Create a parent",
		DefaultStatus: http.StatusCreated,
	}, h.CreateParent)

	huma.Register(g, huma.Operation{
		OperationID:   "update-parent",
		Method:        http.MethodPatch,
		Path:          "",
		Summary:       "Update a parent",
		Description:   "Update a parent",
		DefaultStatus: http.StatusOK,
	}, h.UpdateParent)

	huma.Register(g, huma.Operation{
		OperationID:   "get-parent-by-id",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get a parent by ID with their student children",
		Description:   "Get a parent by ID with their student children",
		DefaultStatus: http.StatusOK,
	}, h.GetParentByID)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-parent",
		Method:        http.MethodDelete,
		Path:          "/{id}",
		Summary:       "Delete a parent",
		Description:   "Delete a parent",
		DefaultStatus: http.StatusOK,
	}, h.DeleteParent)

	huma.Register(g, huma.Operation{
		OperationID:   "list-parents",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List parents",
		Description:   "List parents",
		DefaultStatus: http.StatusOK,
	}, h.ListParents)
}

func (h *ParentsHandler) CreateParent(c context.Context, input *dto.CreateParentReq) (*dto.CreateParentRes, error) {
	parent, err := h.svc.CreateParent(c, input.Body.UserID)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	h.log.Info().Str("id", parent.ParentID).Str("user_id", parent.UserID).Any("created", parent.CreatedAt).
		Msg("Created parent")
	return &dto.CreateParentRes{
		Body: dto.CreateParentResBody{
			ID:        parent.ParentID,
			UserID:    parent.UserID,
			CreatedAt: int(parent.CreatedAt.Unix()),
			UpdatedAt: int(parent.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *ParentsHandler) UpdateParent(c context.Context, input *dto.UpdateParentReq) (*dto.UpdateParentRes, error) {
	m := models.Parents{ParentID: input.Body.ID}
	if input.Body.UserID != nil {
		m.UserID = *input.Body.UserID
	}
	parent, err := h.svc.UpdateParent(c, m)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateParentRes{
		Body: dto.UpdateParentResBody{
			ID:        parent.ParentID,
			UserID:    parent.UserID,
			CreatedAt: int(parent.CreatedAt.Unix()),
			UpdatedAt: int(parent.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *ParentsHandler) GetParentByID(c context.Context, input *dto.GetParentByIDReq) (*dto.GetParentByIDRes, error) {
	parent, err := h.svc.GetParentByID(c, input.ID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("id", input.ID).Int("student_count", len(parent.Students)).Any("created", parent.CreatedAt).
		Msg("Get parent by ID")
	return &dto.GetParentByIDRes{
		Body: dto.GetParentResBody{
			ParentModelRes: *parent,
		},
	}, nil
}

func (h *ParentsHandler) DeleteParent(c context.Context, input *dto.DeleteParentReq) (*dto.DeleteParentRes, error) {
	if err := h.svc.DeleteParent(c, input.ID); err != nil {
		return nil, err
	}
	return &dto.DeleteParentRes{
		Body: dto.DeleteParentResBody{
			ID: input.ID,
		},
	}, nil
}

func (h *ParentsHandler) ListParents(c context.Context, input *dto.ListParentsReq) (*dto.ListParentsRes, error) {
	return h.svc.GetParents(c, input)
}
