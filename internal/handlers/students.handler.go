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

type StudentsHandler struct {
	svc *service.StudentsService
	log zerolog.Logger
}

func RegisterStudentsRoutes(api huma.API, svc *service.StudentsService) {
	h := &StudentsHandler{svc: svc, log: logging.L()}
	g := huma.NewGroup(api, "/students")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Students"}
	})
	g.UseMiddleware(middleware.AuthMiddleware)

	huma.Register(g, huma.Operation{
		OperationID:   "create-student",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create a student",
		Description:   "Create a student",
		DefaultStatus: http.StatusCreated,
	}, h.CreateStudent)

	huma.Register(g, huma.Operation{
		OperationID:   "update-student",
		Method:        http.MethodPatch,
		Path:          "",
		Summary:       "Update a student",
		Description:   "Update a student",
		DefaultStatus: http.StatusOK,
	}, h.UpdateStudent)

	huma.Register(g, huma.Operation{
		OperationID:   "get-student-by-id",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get a student by ID",
		Description:   "Get a student by ID",
		DefaultStatus: http.StatusOK,
	}, h.GetStudentByID)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-student",
		Method:        http.MethodDelete,
		Path:          "/{id}",
		Summary:       "Delete a student",
		Description:   "Delete a student",
		DefaultStatus: http.StatusOK,
	}, h.DeleteStudent)

	huma.Register(g, huma.Operation{
		OperationID:   "list-students",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List students",
		Description:   "List students",
		DefaultStatus: http.StatusOK,
	}, h.ListStudents)
}

func (h *StudentsHandler) CreateStudent(c context.Context, input *dto.CreateStudentReq) (*dto.CreateStudentRes, error) {
	student, err := h.svc.CreateStudent(c, input.Body.Level, input.Body.UserID)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	h.log.Info().Str("id", student.StudentID).Str("level", *student.Level).Any("created", student.CreatedAt).
		Msg("Created student")
	return &dto.CreateStudentRes{
		Body: dto.CreateStudentResBody{
			ID:        student.StudentID,
			Level:     student.Level,
			UserID:    student.UserID,
			CreatedAt: int(student.CreatedAt.Unix()),
			UpdatedAt: int(student.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *StudentsHandler) UpdateStudent(c context.Context, input *dto.UpdateStudentReq) (*dto.UpdateStudentRes, error) {
	m := models.Students{StudentID: input.Body.ID}
	if input.Body.Level != nil {
		m.Level = input.Body.Level
	}
	if input.Body.UserID != nil {
		m.UserID = input.Body.UserID
	}
	student, err := h.svc.UpdateStudent(c, m)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateStudentRes{
		Body: dto.UpdateStudentResBody{
			ID:        student.StudentID,
			Level:     student.Level,
			UserID:    student.UserID,
			CreatedAt: int(student.CreatedAt.Unix()),
			UpdatedAt: int(student.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *StudentsHandler) GetStudentByID(c context.Context, input *dto.GetStudentByIDReq) (*dto.GetStudentByIDRes, error) {
	student, err := h.svc.GetStudentByID(c, input.ID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("id", input.ID).Int("created", student.CreatedAt).
		Msg("Get student by ID")
	return &dto.GetStudentByIDRes{
		Body: *student,
	}, nil
}

func (h *StudentsHandler) DeleteStudent(c context.Context, input *dto.DeleteStudentReq) (*dto.DeleteStudentRes, error) {
	if err := h.svc.DeleteStudent(c, input.ID); err != nil {
		return nil, err
	}
	return &dto.DeleteStudentRes{
		Body: dto.DeleteStudentResBody{
			ID: input.ID,
		},
	}, nil
}

func (h *StudentsHandler) ListStudents(c context.Context, input *dto.ListStudentsReq) (*dto.ListStudentsRes, error) {
	return h.svc.GetStudents(c, input)
}
