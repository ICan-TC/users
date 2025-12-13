package handlers

import (
	"context"
	"net/http"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/middleware"
	"github.com/ICan-TC/users/internal/models"
	"github.com/ICan-TC/users/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
)

type EmployeesHandler struct {
	svc *service.EmployeesService
	log zerolog.Logger
}

func RegisterEmployeesRoutes(api huma.API, svc *service.EmployeesService) {
	h := &EmployeesHandler{svc: svc, log: logging.L()}
	g := huma.NewGroup(api, "/employees")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Employees"}
	})
	g.UseMiddleware(middleware.AuthMiddleware)

	huma.Register(g, huma.Operation{
		OperationID:   "create-employee",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create an employee",
		Description:   "Create an employee",
		DefaultStatus: http.StatusCreated,
	}, h.CreateEmployee)

	huma.Register(g, huma.Operation{
		OperationID:   "update-employee",
		Method:        http.MethodPatch,
		Path:          "",
		Summary:       "Update an employee",
		Description:   "Update an employee",
		DefaultStatus: http.StatusOK,
	}, h.UpdateEmployee)

	huma.Register(g, huma.Operation{
		OperationID:   "get-employee-by-id",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get an employee by ID",
		Description:   "Get an employee by ID",
		DefaultStatus: http.StatusOK,
	}, h.GetEmployeeByID)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-employee",
		Method:        http.MethodDelete,
		Path:          "/{id}",
		Summary:       "Delete an employee",
		Description:   "Delete an employee",
		DefaultStatus: http.StatusOK,
	}, h.DeleteEmployee)

	huma.Register(g, huma.Operation{
		OperationID:   "list-employees",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List employees",
		Description:   "List employees",
		DefaultStatus: http.StatusOK,
	}, h.ListEmployees)
}

func (h *EmployeesHandler) CreateEmployee(c context.Context, input *dto.CreateEmployeeReq) (*dto.CreateEmployeeRes, error) {
	employee, err := h.svc.CreateEmployee(c, input.Body.UserID, input.Body.Role, input.Body.Salary)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	h.log.Info().Str("id", employee.EmployeeID).Str("user_id", employee.UserID).Str("role", employee.Role).
		Float64("salary", employee.Salary).Any("created", employee.CreatedAt).
		Msg("Created employee")
	return &dto.CreateEmployeeRes{
		Body: dto.CreateEmployeeResBody{
			ID:        employee.EmployeeID,
			UserID:    employee.UserID,
			Role:      employee.Role,
			Salary:    employee.Salary,
			CreatedAt: int(employee.CreatedAt.Unix()),
			UpdatedAt: int(employee.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *EmployeesHandler) UpdateEmployee(c context.Context, input *dto.UpdateEmployeeReq) (*dto.UpdateEmployeeRes, error) {
	m := models.Employees{EmployeeID: input.Body.ID}
	if input.Body.Role != nil {
		m.Role = *input.Body.Role
	}
	if input.Body.Salary != nil {
		m.Salary = *input.Body.Salary
	}
	employee, err := h.svc.UpdateEmployee(c, m)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateEmployeeRes{
		Body: dto.UpdateEmployeeResBody{
			ID:        employee.EmployeeID,
			UserID:    employee.UserID,
			Role:      employee.Role,
			Salary:    employee.Salary,
			CreatedAt: int(employee.CreatedAt.Unix()),
			UpdatedAt: int(employee.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *EmployeesHandler) GetEmployeeByID(c context.Context, input *dto.GetEmployeeByIDReq) (*dto.GetEmployeeByIDRes, error) {
	employee, err := h.svc.GetEmployeeByID(c, input.ID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("id", input.ID).Any("created", employee.CreatedAt).
		Msg("Get employee by ID")
	return &dto.GetEmployeeByIDRes{
		Body: dto.GetEmployeeResBody{
			ID:        employee.EmployeeID,
			UserID:    employee.UserID,
			Role:      employee.Role,
			Salary:    employee.Salary,
			CreatedAt: int(employee.CreatedAt.Unix()),
			UpdatedAt: int(employee.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *EmployeesHandler) DeleteEmployee(c context.Context, input *dto.DeleteEmployeeReq) (*dto.DeleteEmployeeRes, error) {
	if err := h.svc.DeleteEmployee(c, input.ID); err != nil {
		return nil, err
	}
	return &dto.DeleteEmployeeRes{
		Body: dto.DeleteEmployeeResBody{
			ID: input.ID,
		},
	}, nil
}

func (h *EmployeesHandler) ListEmployees(c context.Context, input *dto.ListEmployeesReq) (*dto.ListEmployeesRes, error) {
	return h.svc.GetEmployees(c, input)
}
