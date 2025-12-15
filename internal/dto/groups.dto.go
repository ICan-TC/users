package dto

// CreateGroupReq defines the request for creating a group
// All core fields are required except description and metadata
type CreateGroupReq struct {
	AuthHeader
	Body struct {
		Name        string                 `json:"name" doc:"Name of the group" required:"true"`
		Description string                 `json:"description" doc:"Description of the group" required:"false"`
		TeacherID   string                 `json:"teacher_id" doc:"Teacher ID for the group" required:"true"`
		DefaultFee  float64                `json:"default_fee" doc:"Default fee for the group" required:"true"`
		Subject     string                 `json:"subject" doc:"Subject of the group" required:"true"`
		Level       string                 `json:"level" doc:"Level of the group" required:"true"`
		Metadata    map[string]interface{} `json:"metadata" doc:"Additional metadata" required:"false"`
	}
}

type CreateGroupRes struct{ Body CreateGroupResBody }
type CreateGroupResBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	TeacherID   string                 `json:"teacher_id"`
	DefaultFee  float64                `json:"default_fee"`
	Subject     string                 `json:"subject"`
	Level       string                 `json:"level"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   int                    `json:"created_at"`
	UpdatedAt   int                    `json:"updated_at"`
}

// UpdateGroupReq for updating a group
// All fields except ID are optional
type UpdateGroupReq struct {
	AuthHeader
	Body struct {
		ID          string                 `json:"id" doc:"ID of the group" required:"true"`
		Name        *string                `json:"name" doc:"Name of the group" required:"false"`
		Description *string                `json:"description" doc:"Description of the group" required:"false"`
		TeacherID   *string                `json:"teacher_id" doc:"Teacher ID for the group" required:"false"`
		DefaultFee  *float64               `json:"default_fee" doc:"Default fee for the group" required:"false"`
		Subject     *string                `json:"subject" doc:"Subject of the group" required:"false"`
		Level       *string                `json:"level" doc:"Level of the group" required:"false"`
		Metadata    map[string]interface{} `json:"metadata" doc:"Additional metadata" required:"false"`
	}
}
type UpdateGroupRes struct{ Body UpdateGroupResBody }
type UpdateGroupResBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	TeacherID   string                 `json:"teacher_id"`
	DefaultFee  float64                `json:"default_fee"`
	Subject     string                 `json:"subject"`
	Level       string                 `json:"level"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   int                    `json:"created_at"`
	UpdatedAt   int                    `json:"updated_at"`
}

type GetGroupByIDReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the group" required:"true"`
}
type GetGroupByIDRes struct{ Body GetGroupResBody }

type GetGroupResBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	TeacherID   string                 `json:"teacher_id"`
	DefaultFee  float64                `json:"default_fee"`
	Subject     string                 `json:"subject"`
	Level       string                 `json:"level"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   int                    `json:"created_at"`
	UpdatedAt   int                    `json:"updated_at"`
}

type DeleteGroupReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the group" required:"true"`
}
type DeleteGroupResBody struct {
	ID string `json:"id"`
}
type DeleteGroupRes struct {
	Body DeleteGroupResBody
}

type ListGroupsReq struct {
	AuthHeader
	ListQuery
}
type ListGroupsResBody struct {
	Groups    []GetGroupResBody `json:"groups"`
	Total     int               `json:"total"`
	ListQuery ListQuery         `json:"query"`
}
type ListGroupsRes struct {
	Body ListGroupsResBody
}
