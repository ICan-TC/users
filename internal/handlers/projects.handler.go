package handlers

import (
	"context"
	"net/http"

	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/service"
	"github.com/danielgtaylor/huma/v2"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func RegisterProjectRoutes(api huma.API, service *service.ProjectService) {
	h := &ProjectHandler{svc: service}
	g := huma.NewGroup(api, "/projects")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"projects"}
	})

	huma.Register(g, huma.Operation{
		OperationID:   "create-project",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create project",
		Description:   "Create a new project",
		DefaultStatus: http.StatusCreated,
	}, h.Create)

	huma.Register(g, huma.Operation{
		OperationID:   "update-project",
		Method:        http.MethodPut,
		Path:          "",
		Summary:       "Update project",
		Description:   "Update an existing project",
		DefaultStatus: http.StatusOK,
	}, h.Update)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-project",
		Method:        http.MethodDelete,
		Path:          "",
		Summary:       "Delete project",
		Description:   "Delete an existing project",
		DefaultStatus: http.StatusOK,
	}, h.Delete)

	huma.Register(g, huma.Operation{
		OperationID:   "get-project",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get project",
		Description:   "Get an existing project",
		DefaultStatus: http.StatusOK,
	}, h.Get)

	huma.Register(g, huma.Operation{
		OperationID:   "search-project",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "Search project",
		Description:   "Search for projects",
		DefaultStatus: http.StatusOK,
	}, h.Search)

	huma.Register(g, huma.Operation{
		OperationID:   "list-project",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List project",
		Description:   "List projects",
		DefaultStatus: http.StatusOK,
	}, h.List)

}

func (h *ProjectHandler) Create(ctx context.Context, input *dto.CreateProjectReq) (*dto.CreateProjectRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *ProjectHandler) Update(ctx context.Context, input *dto.UpdateProjectReq) (*dto.UpdateProjectRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *ProjectHandler) Delete(ctx context.Context, input *dto.DeleteProjectReq) (*dto.DeleteProjectRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *ProjectHandler) Get(ctx context.Context, input *dto.GetProjectReq) (*dto.GetProjectRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *ProjectHandler) Search(ctx context.Context, input *dto.SearchProjectReq) (*dto.SearchProjectRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *ProjectHandler) List(ctx context.Context, input *dto.ListProjectReq) (*dto.ListProjectRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
