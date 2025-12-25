package dto

// CreateStudentReq defines the request for creating a student
// Follows the pattern of CreateUserReq
// Level is required, UserID is optional (for linking to a user)
type CreateStudentReq struct {
	AuthHeader
	Body struct {
		Level  string  `json:"level" doc:"Level of the student" required:"true"`
		UserID *string `json:"user_id" doc:"User ID to link the student to" required:"false"`
	}
}

type CreateStudentRes struct{ Body CreateStudentResBody }
type CreateStudentResBody struct {
	ID        string  `json:"id"`
	Level     *string `json:"level"`
	UserID    *string `json:"user_id"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}

// UpdateStudentReq for updating a student
// All fields except ID are optional
// Follows UpdateUserReq pattern
type UpdateStudentReq struct {
	AuthHeader
	Body struct {
		ID     string  `json:"id" doc:"ID of the student" required:"true"`
		Level  *string `json:"level" doc:"Level of the student" required:"false"`
		UserID *string `json:"user_id" doc:"User ID to link the student to" required:"false"`
	}
}
type UpdateStudentRes struct{ Body UpdateStudentResBody }
type UpdateStudentResBody struct {
	ID        string  `json:"id"`
	Level     *string `json:"level"`
	UserID    *string `json:"user_id"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}

type GetStudentByIDReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the student" required:"true"`
}
type GetStudentByIDRes struct{ Body StudentsModelRes }

type GetStudentResBody struct {
	ID        string  `json:"id"`
	Level     *string `json:"level"`
	UserID    *string `json:"user_id"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}

type DeleteStudentReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the student" required:"true"`
}
type DeleteStudentResBody struct {
	ID string `json:"id"`
}
type DeleteStudentRes struct {
	Body DeleteStudentResBody
}

type ListStudentsReq struct {
	AuthHeader
	ListQuery
}
type ListStudentsResBody struct {
	Students  []StudentsModelRes `json:"students"`
	Total     int                `json:"total"`
	ListQuery ListQuery          `json:"query"`
}
type ListStudentsRes struct {
	Body ListStudentsResBody
}

// StudentsModelRes represents a student with embedded user information
type StudentsModelRes struct {
	ID        string  `json:"id"`
	Level     *string `json:"level"`
	UserID    *string `json:"user_id"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
	GetUserResBody
}
