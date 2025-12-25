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

type TeachersHandler struct {
	svc *service.TeachersService
	log zerolog.Logger
}

func RegisterTeachersRoutes(api huma.API, svc *service.TeachersService) {
	h := &TeachersHandler{svc: svc, log: logging.L()}
	g := huma.NewGroup(api, "/teachers")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Teachers"}
	})
	g.UseMiddleware(middleware.AuthMiddleware)

	huma.Register(g, huma.Operation{
		OperationID:   "create-teacher",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create a teacher",
		Description:   "Create a teacher",
		DefaultStatus: http.StatusCreated,
	}, h.CreateTeacher)

	huma.Register(g, huma.Operation{
		OperationID:   "update-teacher",
		Method:        http.MethodPatch,
		Path:          "",
		Summary:       "Update a teacher",
		Description:   "Update a teacher",
		DefaultStatus: http.StatusOK,
	}, h.UpdateTeacher)

	huma.Register(g, huma.Operation{
		OperationID:   "get-teacher-by-id",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get a teacher by ID",
		Description:   "Get a teacher by ID",
		DefaultStatus: http.StatusOK,
	}, h.GetTeacherByID)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-teacher",
		Method:        http.MethodDelete,
		Path:          "/{id}",
		Summary:       "Delete a teacher",
		Description:   "Delete a teacher",
		DefaultStatus: http.StatusOK,
	}, h.DeleteTeacher)

	huma.Register(g, huma.Operation{
		OperationID:   "list-teachers",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List teachers",
		Description:   "List teachers",
		DefaultStatus: http.StatusOK,
	}, h.ListTeachers)
}

func (h *TeachersHandler) CreateTeacher(c context.Context, input *dto.CreateTeacherReq) (*dto.CreateTeacherRes, error) {
	teacher, err := h.svc.CreateTeacher(c, input.Body.UserID)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	h.log.Info().Str("id", teacher.TeacherID).Str("user_id", *teacher.UserID).Any("created", teacher.CreatedAt).
		Msg("Created teacher")
	return &dto.CreateTeacherRes{
		Body: dto.CreateTeacherResBody{
			ID:        teacher.TeacherID,
			UserID:    *teacher.UserID,
			CreatedAt: int(teacher.CreatedAt.Unix()),
			UpdatedAt: int(teacher.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *TeachersHandler) UpdateTeacher(c context.Context, input *dto.UpdateTeacherReq) (*dto.UpdateTeacherRes, error) {
	m := models.Teachers{TeacherID: input.Body.ID}
	if input.Body.UserID != nil {
		m.UserID = input.Body.UserID
	}
	teacher, err := h.svc.UpdateTeacher(c, m)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateTeacherRes{
		Body: dto.UpdateTeacherResBody{
			ID:        teacher.TeacherID,
			UserID:    *teacher.UserID,
			CreatedAt: int(teacher.CreatedAt.Unix()),
			UpdatedAt: int(teacher.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *TeachersHandler) GetTeacherByID(c context.Context, input *dto.GetTeacherByIDReq) (*dto.GetTeacherByIDRes, error) {
	teacher, err := h.svc.GetTeacherByID(c, input.ID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("id", input.ID).Int("created", teacher.CreatedAt).
		Msg("Get teacher by ID")
	return &dto.GetTeacherByIDRes{
		Body: *teacher,
	}, nil
}

func (h *TeachersHandler) DeleteTeacher(c context.Context, input *dto.DeleteTeacherReq) (*dto.DeleteTeacherRes, error) {
	if err := h.svc.DeleteTeacher(c, input.ID); err != nil {
		return nil, err
	}
	return &dto.DeleteTeacherRes{
		Body: dto.DeleteTeacherResBody{
			ID: input.ID,
		},
	}, nil
}

func (h *TeachersHandler) ListTeachers(c context.Context, input *dto.ListTeachersReq) (*dto.ListTeachersRes, error) {
	return h.svc.GetTeachers(c, input)
}
