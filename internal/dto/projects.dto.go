package dto

import "github.com/ICan-TC/users/internal/models"

type CreateProjectReq struct {
	Body struct {
		Name     string                 `json:"name" doc:"Name of the project" minLength:"1" maxLength:"100"`
		Logo     string                 `json:"logo" doc:"Logo URL"`
		Tags     []string               `json:"tags,omitempty" doc:"Tags for categorization"`
		Labels   map[string]interface{} `json:"labels,omitempty" doc:"Custom labels (JSON object)"`
		Metadata map[string]interface{} `json:"metadata,omitempty" doc:"Additional metadata (JSON object)"`
	}
}

type CreateProjectRes struct {
	Body models.Project
}

type UpdateProjectReq struct {
	Body struct {
		Name     string                 `json:"name" doc:"Name of the project" minLength:"1" maxLength:"100"`
		Logo     string                 `json:"logo" doc:"Logo URL"`
		Tags     []string               `json:"tags,omitempty" doc:"Tags for categorization"`
		Labels   map[string]interface{} `json:"labels,omitempty" doc:"Custom labels (JSON object)"`
		Metadata map[string]interface{} `json:"metadata,omitempty" doc:"Additional metadata (JSON object)"`
	}
}

type UpdateProjectRes struct {
	Body models.Project
}

type DeleteProjectReq struct {
	Body struct {
		ID string `json:"id" doc:"ID of the project"`
	}
}
type DeleteProjectRes struct {
	Body struct {
		ID string `json:"id" doc:"ID of the project"`
	}
}
type GetProjectReq struct {
	ID string `query:"id" doc:"ID of the project"`
}
type GetProjectRes struct {
	Body *models.Project
}
type SearchProjectReq struct {
	Query struct {
		Query string `query:"q" doc:"Query string"`
	}
}
type SearchProjectRes struct {
	Body []*models.Project
}
type ListProjectReq struct {
	ListQuery
}
type ListProjectRes struct {
	Body struct {
		Items     []*models.Project `json:"items"`
		Page      int               `json:"page"`
		PerPage   int               `json:"per_page"`
		Total     int               `json:"total"`
		ListQuery ListQuery         `json:"query"`
	}
}
