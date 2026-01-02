package dto

// CreateParentReq defines the request for creating a parent
// UserID is required for linking to a user
type CreateParentReq struct {
	AuthHeader
	Body struct {
		UserID string `json:"user_id" doc:"User ID to link the parent to" required:"true"`
	}
}

type CreateParentRes struct{ Body CreateParentResBody }
type CreateParentResBody struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

// UpdateParentReq for updating a parent
// All fields except ID are optional
type UpdateParentReq struct {
	AuthHeader
	Body struct {
		ID     string  `json:"id" doc:"ID of the parent" required:"true"`
		UserID *string `json:"user_id" doc:"User ID to link the parent to" required:"false"`
	}
}
type UpdateParentRes struct{ Body UpdateParentResBody }
type UpdateParentResBody struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type GetParentByIDReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the parent" required:"true"`
}
type GetParentByIDRes struct{ Body GetParentResBody }

// GetParentResBody includes the parent info along with their student children
type GetParentResBody struct {
	ParentModelRes
}

type DeleteParentReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the parent" required:"true"`
}
type DeleteParentResBody struct {
	ID string `json:"id"`
}
type DeleteParentRes struct {
	Body DeleteParentResBody
}

type ListParentsReq struct {
	AuthHeader
	ListQuery
}
type ListParentsResBody struct {
	Parents   []ParentModelRes `json:"parents"`
	Total     int              `json:"total"`
	ListQuery ListQuery        `json:"query"`
}
type ListParentsRes struct {
	Body ListParentsResBody
}

type ParentModelRes struct {
	ID        string              `json:"id"`
	UserID    string              `json:"user_id"`
	Students  []GetStudentResBody `json:"students,omitempty"`
	CreatedAt int                 `json:"created_at"`
	UpdatedAt int                 `json:"updated_at"`
	UserModelRes
}
