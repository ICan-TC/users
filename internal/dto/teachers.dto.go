package dto

// CreateTeacherReq defines the request for creating a teacher
// Follows the pattern of CreateStudentReq
// UserID is required for linking to a user
type CreateTeacherReq struct {
	AuthHeader
	Body struct {
		UserID string `json:"user_id" doc:"User ID to link the teacher to" required:"true"`
	}
}

type CreateTeacherRes struct{ Body CreateTeacherResBody }
type CreateTeacherResBody struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

// UpdateTeacherReq for updating a teacher
// All fields except ID are optional
// Follows UpdateStudentReq pattern
type UpdateTeacherReq struct {
	AuthHeader
	Body struct {
		ID     string  `json:"id" doc:"ID of the teacher" required:"true"`
		UserID *string `json:"user_id" doc:"User ID to link the teacher to" required:"false"`
	}
}
type UpdateTeacherRes struct{ Body UpdateTeacherResBody }
type UpdateTeacherResBody struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type GetTeacherByIDReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the teacher" required:"true"`
}
type GetTeacherByIDRes struct{ Body TeachersModelRes }

type GetTeacherResBody struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type DeleteTeacherReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the teacher" required:"true"`
}
type DeleteTeacherResBody struct {
	ID string `json:"id"`
}
type DeleteTeacherRes struct {
	Body DeleteTeacherResBody
}

type ListTeachersReq struct {
	AuthHeader
	ListQuery
}
type ListTeachersResBody struct {
	Teachers  []TeachersModelRes `json:"teachers"`
	Total     int                `json:"total"`
	ListQuery ListQueryRes       `json:"query"`
}
type ListTeachersRes struct {
	Body ListTeachersResBody
}

// TeachersModelRes represents a teacher with embedded user information
type TeachersModelRes struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	UserModelRes
}
