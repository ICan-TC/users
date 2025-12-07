package dto

import "github.com/ICan-TC/users/internal/models"

type CreateWebsiteReq struct {
	Body struct {
		Name     string                 `json:"name" doc:"Name of the website" minLength:"1" maxLength:"100"`
		Logo     string                 `json:"logo" doc:"Logo URL"`
		Tags     []string               `json:"tags,omitempty" doc:"Tags for categorization"`
		Labels   map[string]interface{} `json:"labels,omitempty" doc:"Custom labels (JSON object)"`
		Metadata map[string]interface{} `json:"metadata,omitempty" doc:"Additional metadata (JSON object)"`
	}
}

type CreateWebsiteRes struct {
	Body models.Website
}

type UpdateWebsiteReq struct {
	Body struct {
		Name     string                 `json:"name" doc:"Name of the website" minLength:"1" maxLength:"100"`
		Logo     string                 `json:"logo" doc:"Logo URL"`
		Tags     []string               `json:"tags,omitempty" doc:"Tags for categorization"`
		Labels   map[string]interface{} `json:"labels,omitempty" doc:"Custom labels (JSON object)"`
		Metadata map[string]interface{} `json:"metadata,omitempty" doc:"Additional metadata (JSON object)"`
	}
}

type UpdateWebsiteRes struct {
	Body models.Website
}

type DeleteWebsiteReq struct {
	Body struct {
		ID string `json:"id" doc:"ID of the website"`
	}
}
type DeleteWebsiteRes struct {
	Body struct {
		ID string `json:"id" doc:"ID of the website"`
	}
}
type GetWebsiteReq struct {
	ID string `query:"id" doc:"ID of the website"`
}
type GetWebsiteRes struct {
	Body *models.Website
}
type SearchWebsiteReq struct {
	Query struct {
		Query string `query:"q" doc:"Query string"`
	}
}
type SearchWebsiteRes struct {
	Body []*models.Website
}
type ListWebsiteReq struct {
	ListQuery
}
type ListWebsiteRes struct {
	Body struct {
		Items     []*models.Website `json:"items"`
		Total     int               `json:"total"`
		ListQuery ListQuery         `json:"query"`
	}
}
