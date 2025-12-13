package dto

// CreateEmployeeReq defines the request for creating an employee
// Follows the pattern of CreateStudentReq
// UserID, Role, and Salary are required
type CreateEmployeeReq struct {
	AuthHeader
	Body struct {
		UserID string  `json:"user_id" doc:"User ID to link the employee to" required:"true"`
		Role   string  `json:"role" doc:"Role of the employee" required:"true"`
		Salary float64 `json:"salary" doc:"Salary of the employee" required:"true"`
	}
}

type CreateEmployeeRes struct{ Body CreateEmployeeResBody }
type CreateEmployeeResBody struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Role      string  `json:"role"`
	Salary    float64 `json:"salary"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}

// UpdateEmployeeReq for updating an employee
// All fields except ID are optional
// Follows UpdateStudentReq pattern
type UpdateEmployeeReq struct {
	AuthHeader
	Body struct {
		ID     string   `json:"id" doc:"ID of the employee" required:"true"`
		Role   *string  `json:"role" doc:"Role of the employee" required:"false"`
		Salary *float64 `json:"salary" doc:"Salary of the employee" required:"false"`
	}
}
type UpdateEmployeeRes struct{ Body UpdateEmployeeResBody }
type UpdateEmployeeResBody struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Role      string  `json:"role"`
	Salary    float64 `json:"salary"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}

type GetEmployeeByIDReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the employee" required:"true"`
}
type GetEmployeeByIDRes struct{ Body GetEmployeeResBody }

type GetEmployeeResBody struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Role      string  `json:"role"`
	Salary    float64 `json:"salary"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}

type DeleteEmployeeReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the employee" required:"true"`
}
type DeleteEmployeeResBody struct {
	ID string `json:"id"`
}
type DeleteEmployeeRes struct {
	Body DeleteEmployeeResBody
}

type ListEmployeesReq struct {
	AuthHeader
	ListQuery
}
type ListEmployeesResBody struct {
	Employees []GetEmployeeResBody `json:"employees"`
	Total     int                  `json:"total"`
	ListQuery ListQuery            `json:"query"`
}
type ListEmployeesRes struct {
	Body ListEmployeesResBody
}
