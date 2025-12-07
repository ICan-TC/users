package handlers

import (
	"context"
	"net/http"

	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/service"
	"github.com/danielgtaylor/huma/v2"
)

type OrganizationHandler struct {
	svc *service.OrganizationService
}

func RegisterOrganizationRoutes(api huma.API, service *service.OrganizationService) {
	h := &OrganizationHandler{svc: service}
	g := huma.NewGroup(api, "/organizations")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"organizations"}
	})

	huma.Register(g, huma.Operation{
		OperationID:   "create-organization",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create organization",
		Description:   "Create a new organization",
		DefaultStatus: http.StatusCreated,
	}, h.Create)

	huma.Register(g, huma.Operation{
		OperationID:   "update-organization",
		Method:        http.MethodPut,
		Path:          "",
		Summary:       "Update organization",
		Description:   "Update an existing organization",
		DefaultStatus: http.StatusOK,
	}, h.Update)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-organization",
		Method:        http.MethodDelete,
		Path:          "",
		Summary:       "Delete organization",
		Description:   "Delete an existing organization",
		DefaultStatus: http.StatusOK,
	}, h.Delete)

	huma.Register(g, huma.Operation{
		OperationID:   "get-organization",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get organization",
		Description:   "Get an existing organization",
		DefaultStatus: http.StatusOK,
	}, h.Get)

	huma.Register(g, huma.Operation{
		OperationID:   "search-organization",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "Search organization",
		Description:   "Search for organizations",
		DefaultStatus: http.StatusOK,
	}, h.Search)

	huma.Register(g, huma.Operation{
		OperationID:   "list-organization",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List organization",
		Description:   "List organizations",
		DefaultStatus: http.StatusOK,
	}, h.List)

}

func (h *OrganizationHandler) Create(ctx context.Context, input *dto.CreateOrganizationReq) (*dto.CreateOrganizationRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *OrganizationHandler) Update(ctx context.Context, input *dto.UpdateOrganizationReq) (*dto.UpdateOrganizationRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *OrganizationHandler) Delete(ctx context.Context, input *dto.DeleteOrganizationReq) (*dto.DeleteOrganizationRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *OrganizationHandler) Get(ctx context.Context, input *dto.GetOrganizationReq) (*dto.GetOrganizationRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *OrganizationHandler) Search(ctx context.Context, input *dto.SearchOrganizationReq) (*dto.SearchOrganizationRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *OrganizationHandler) List(ctx context.Context, input *dto.ListOrganizationReq) (*dto.ListOrganizationRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
