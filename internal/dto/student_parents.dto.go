package dto

// CreateStudentParentReq defines the request for linking a student to a parent
type CreateStudentParentReq struct {
	AuthHeader
	Body struct {
		StudentID string `json:"student_id" doc:"ID of the student" required:"true"`
		ParentID  string `json:"parent_id" doc:"ID of the parent" required:"true"`
	}
}

type CreateStudentParentRes struct{ Body CreateStudentParentResBody }
type CreateStudentParentResBody struct {
	StudentID string `json:"student_id"`
	ParentID  string `json:"parent_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

// UpdateStudentParentReq for updating a student-parent relationship
// Currently only the IDs can be provided (immutable composite key)
type UpdateStudentParentReq struct {
	AuthHeader
	Body struct {
		StudentID string `json:"student_id" doc:"ID of the student" required:"true"`
		ParentID  string `json:"parent_id" doc:"ID of the parent" required:"true"`
	}
}
type UpdateStudentParentRes struct{ Body UpdateStudentParentResBody }
type UpdateStudentParentResBody struct {
	StudentID string `json:"student_id"`
	ParentID  string `json:"parent_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type GetStudentParentReq struct {
	AuthHeader
	StudentID string `path:"student_id" doc:"ID of the student" required:"true"`
	ParentID  string `path:"parent_id" doc:"ID of the parent" required:"true"`
}
type GetStudentParentRes struct{ Body GetStudentParentResBody }

type GetStudentParentResBody struct {
	StudentID string `json:"student_id"`
	ParentID  string `json:"parent_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type DeleteStudentParentReq struct {
	AuthHeader
	StudentID string `path:"student_id" doc:"ID of the student" required:"true"`
	ParentID  string `path:"parent_id" doc:"ID of the parent" required:"true"`
}
type DeleteStudentParentResBody struct {
	StudentID string `json:"student_id"`
	ParentID  string `json:"parent_id"`
}
type DeleteStudentParentRes struct {
	Body DeleteStudentParentResBody
}

type ListStudentParentsReq struct {
	AuthHeader
	ListQuery
}
type ListStudentParentsResBody struct {
	StudentParents []GetStudentParentResBody `json:"student_parents"`
	Total          int                       `json:"total"`
	ListQuery      ListQuery                 `json:"query"`
}
type ListStudentParentsRes struct {
	Body ListStudentParentsResBody
}
