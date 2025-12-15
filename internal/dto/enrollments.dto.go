package dto

// CreateEnrollmentReq defines the request for creating an enrollment
// StudentID, GroupID, and Fee are required
type CreateEnrollmentReq struct {
	AuthHeader
	Body struct {
		StudentID string  `json:"student_id" doc:"Student ID to enroll" required:"true"`
		GroupID   string  `json:"group_id" doc:"Group ID to enroll in" required:"true"`
		Fee       float64 `json:"fee" doc:"Fee for the enrollment" required:"true"`
	}
}

type CreateEnrollmentRes struct{ Body EnrollmentModelRes }

// UpdateEnrollmentReq for updating an enrollment
// All fields except StudentID and GroupID are optional
type UpdateEnrollmentReq struct {
	AuthHeader
	Body struct {
		StudentID string   `json:"student_id" doc:"Student ID" required:"true"`
		GroupID   string   `json:"group_id" doc:"Group ID" required:"true"`
		Fee       *float64 `json:"fee" doc:"Fee for the enrollment" required:"false"`
	}
}

type UpdateEnrollmentRes struct{ Body EnrollmentModelRes }

type GetEnrollmentByIDReq struct {
	AuthHeader
	StudentID string `path:"student_id" doc:"Student ID" required:"true"`
	GroupID   string `path:"group_id" doc:"Group ID" required:"true"`
}

type GetEnrollmentByIDRes struct{ Body EnrollmentModelRes }

type DeleteEnrollmentReq struct {
	AuthHeader
	StudentID string `path:"student_id" doc:"Student ID" required:"true"`
	GroupID   string `path:"group_id" doc:"Group ID" required:"true"`
}

type DeleteEnrollmentResBody struct {
	StudentID string `json:"student_id"`
	GroupID   string `json:"group_id"`
}

type DeleteEnrollmentRes struct {
	Body DeleteEnrollmentResBody
}

type ListEnrollmentsReq struct {
	AuthHeader
	ListQuery
}

type ListEnrollmentsResBody struct {
	Enrollments []EnrollmentModelRes `json:"enrollments"`
	Total       int                  `json:"total"`
	ListQuery   ListQuery            `json:"query"`
}

type ListEnrollmentsRes struct {
	Body ListEnrollmentsResBody
}

type EnrollmentModelRes struct {
	StudentID string  `json:"student_id"`
	GroupID   string  `json:"group_id"`
	Fee       float64 `json:"fee"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}
