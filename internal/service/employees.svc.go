package service

import (
	"context"
	"strings"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type EmployeesService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewEmployeesService(db *bun.DB) (*EmployeesService, error) {
	log := logging.L().With().Str("service", "employees.svc").Logger()
	return &EmployeesService{log: log, db: db}, nil
}

func (s *EmployeesService) GetEmployees(ctx context.Context, params *dto.ListEmployeesReq) (*dto.ListEmployeesRes, error) {
	total, err := s.db.NewSelect().Model((*models.Employees)(nil)).Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	var employees []models.Employees
	res := &dto.ListEmployeesRes{
		Body: dto.ListEmployeesResBody{
			Total:     0,
			ListQuery: params.ListQuery,
			Employees: nil,
		},
	}
	res.Body.Total = total

	q := s.db.NewSelect().Model(&employees)
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where("(role ILIKE ? OR user_id ILIKE ?)", search, search)
	}
	q = q.Order(params.SortBy + " " + params.SortDir)
	q = q.Limit(params.PerPage)
	q = q.Offset(params.PerPage * (params.Page - 1))

	if err := q.Scan(ctx, &employees); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return res, nil
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resEmployees := []dto.GetEmployeeResBody{}
	for _, emp := range employees {
		newEmployee := dto.GetEmployeeResBody{
			ID:        emp.EmployeeID,
			UserID:    emp.UserID,
			Role:      emp.Role,
			Salary:    emp.Salary,
			CreatedAt: int(emp.CreatedAt.Unix()),
			UpdatedAt: int(emp.UpdatedAt.Unix()),
		}
		resEmployees = append(resEmployees, newEmployee)
	}
	res.Body.Employees = resEmployees
	return res, nil
}

func (s *EmployeesService) GetEmployeeByID(ctx context.Context, id string) (*models.Employees, error) {
	m := models.Employees{EmployeeID: id}
	if err := s.db.NewSelect().Model(&m).WherePK("id").Scan(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't get employee")
		return nil, huma.Error404NotFound("employee not found")
	}
	return &m, nil
}

func (s *EmployeesService) CreateEmployee(ctx context.Context, userID string, role string, salary float64) (*models.Employees, error) {
	if _, err := ulid.Parse(userID); err != nil {
		return nil, huma.Error400BadRequest("userID is invalid", err)
	}
	m := models.Employees{
		EmployeeID: userID,
		UserID:     userID,
		Role:       role,
		Salary:     salary,
	}
	if _, err := s.db.NewInsert().Model(&m).Returning("*").Exec(ctx, &m); err != nil {
		s.log.Err(err).Msg("Couldn't insert employee")
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *EmployeesService) UpdateEmployee(ctx context.Context, employee models.Employees) (*models.Employees, error) {
	m := employee
	m.EmployeeID = employee.EmployeeID
	if err := s.db.NewUpdate().Model(&m).Returning("*").OmitZero().WherePK("id").Scan(ctx, &m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, huma.Error404NotFound("employee not found")
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &m, nil
}

func (s *EmployeesService) DeleteEmployee(ctx context.Context, id string) error {
	m := models.Employees{EmployeeID: id}
	if _, err := s.db.NewDelete().Model(&m).WherePK("id").Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return huma.Error404NotFound("employee not found")
		}
		return huma.Error500InternalServerError(err.Error())
	}
	return nil
}
