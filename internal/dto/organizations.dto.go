package dto

import "github.com/ICan-TC/users/internal/models"

type CreateOrganizationReq struct {
	Body struct {
		Name     string                 `json:"name" doc:"Name of the organization" minLength:"1" maxLength:"100"`
		Logo     string                 `json:"logo" doc:"Logo URL"`
		Tags     []string               `json:"tags,omitempty" doc:"Tags for categorization"`
		Labels   map[string]interface{} `json:"labels,omitempty" doc:"Custom labels (JSON object)"`
		Metadata map[string]interface{} `json:"metadata,omitempty" doc:"Additional metadata (JSON object)"`
	}
}

type CreateOrganizationRes struct {
	Body models.Organization
}

type UpdateOrganizationReq struct {
	Body struct {
		Name     string                 `json:"name" doc:"Name of the organization" minLength:"1" maxLength:"100"`
		Logo     string                 `json:"logo" doc:"Logo URL"`
		Tags     []string               `json:"tags,omitempty" doc:"Tags for categorization"`
		Labels   map[string]interface{} `json:"labels,omitempty" doc:"Custom labels (JSON object)"`
		Metadata map[string]interface{} `json:"metadata,omitempty" doc:"Additional metadata (JSON object)"`
	}
}

type UpdateOrganizationRes struct {
	Body models.Organization
}

type DeleteOrganizationReq struct {
	Body struct {
		ID string `json:"id" doc:"ID of the organization"`
	}
}
type DeleteOrganizationRes struct {
	Body struct {
		ID string `json:"id" doc:"ID of the organization"`
	}
}
type GetOrganizationReq struct {
	ID string `query:"id" doc:"ID of the organization"`
}
type GetOrganizationRes struct {
	Body *models.Organization
}
type SearchOrganizationReq struct {
	Query struct {
		Query string `query:"q" doc:"Query string"`
	}
}
type SearchOrganizationRes struct {
	Body []*models.Organization
}
type ListOrganizationReq struct {
	ListQuery
}
type ListOrganizationRes struct {
	Body struct {
		Items     []*models.Organization `json:"items"`
		Page      int                    `json:"page"`
		PerPage   int                    `json:"per_page"`
		Total     int                    `json:"total"`
		ListQuery ListQuery              `json:"query"`
	}
}
