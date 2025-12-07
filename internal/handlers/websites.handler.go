package handlers

import (
	"context"
	"net/http"

	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/service"
	"github.com/danielgtaylor/huma/v2"
)

type WebsiteHandler struct {
	svc *service.WebsiteService
}

func RegisterWebsiteRoutes(api huma.API, service *service.WebsiteService) {
	h := &WebsiteHandler{svc: service}
	g := huma.NewGroup(api, "/websites")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"websites"}
	})

	huma.Register(g, huma.Operation{
		OperationID:   "create-website",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create website",
		Description:   "Create a new website",
		DefaultStatus: http.StatusCreated,
	}, h.Create)

	huma.Register(g, huma.Operation{
		OperationID:   "update-website",
		Method:        http.MethodPut,
		Path:          "",
		Summary:       "Update website",
		Description:   "Update an existing website",
		DefaultStatus: http.StatusOK,
	}, h.Update)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-website",
		Method:        http.MethodDelete,
		Path:          "",
		Summary:       "Delete website",
		Description:   "Delete an existing website",
		DefaultStatus: http.StatusOK,
	}, h.Delete)

	huma.Register(g, huma.Operation{
		OperationID:   "get-website",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get website",
		Description:   "Get an existing website",
		DefaultStatus: http.StatusOK,
	}, h.Get)

	huma.Register(g, huma.Operation{
		OperationID:   "search-website",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "Search website",
		Description:   "Search for websites",
		DefaultStatus: http.StatusOK,
	}, h.Search)

	huma.Register(g, huma.Operation{
		OperationID:   "list-website",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List website",
		Description:   "List websites",
		DefaultStatus: http.StatusOK,
	}, h.List)

}

func (h *WebsiteHandler) Create(ctx context.Context, input *dto.CreateWebsiteReq) (*dto.CreateWebsiteRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *WebsiteHandler) Update(ctx context.Context, input *dto.UpdateWebsiteReq) (*dto.UpdateWebsiteRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *WebsiteHandler) Delete(ctx context.Context, input *dto.DeleteWebsiteReq) (*dto.DeleteWebsiteRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *WebsiteHandler) Get(ctx context.Context, input *dto.GetWebsiteReq) (*dto.GetWebsiteRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *WebsiteHandler) Search(ctx context.Context, input *dto.SearchWebsiteReq) (*dto.SearchWebsiteRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *WebsiteHandler) List(ctx context.Context, input *dto.ListWebsiteReq) (*dto.ListWebsiteRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
