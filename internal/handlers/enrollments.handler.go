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

type EnrollmentsHandler struct {
	svc *service.EnrollmentsService
	log zerolog.Logger
}

func RegisterEnrollmentsRoutes(api huma.API, svc *service.EnrollmentsService) {
	h := &EnrollmentsHandler{svc: svc, log: logging.L()}
	g := huma.NewGroup(api, "/enrollments")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Enrollments"}
	})
	g.UseMiddleware(middleware.AuthMiddleware)

	huma.Register(g, huma.Operation{
		OperationID:   "create-enrollment",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create an enrollment",
		Description:   "Create an enrollment",
		DefaultStatus: http.StatusCreated,
	}, h.CreateEnrollment)

	huma.Register(g, huma.Operation{
		OperationID:   "update-enrollment",
		Method:        http.MethodPatch,
		Path:          "/{student_id}/{group_id}",
		Summary:       "Update an enrollment",
		Description:   "Update an enrollment",
		DefaultStatus: http.StatusOK,
	}, h.UpdateEnrollment)

	huma.Register(g, huma.Operation{
		OperationID:   "get-enrollment-by-id",
		Method:        http.MethodGet,
		Path:          "/{student_id}/{group_id}",
		Summary:       "Get an enrollment by student and group ID",
		Description:   "Get an enrollment by student and group ID",
		DefaultStatus: http.StatusOK,
	}, h.GetEnrollmentByID)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-enrollment",
		Method:        http.MethodDelete,
		Path:          "/{student_id}/{group_id}",
		Summary:       "Delete an enrollment",
		Description:   "Delete an enrollment",
		DefaultStatus: http.StatusOK,
	}, h.DeleteEnrollment)

	huma.Register(g, huma.Operation{
		OperationID:   "list-enrollments",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List enrollments",
		Description:   "List enrollments",
		DefaultStatus: http.StatusOK,
	}, h.ListEnrollments)

	huma.Register(g, huma.Operation{
		OperationID:   "get-enrollments-by-group",
		Method:        http.MethodGet,
		Path:          "/group/{group_id}",
		Summary:       "Get all enrollments for a group",
		Description:   "Get all enrollments for a specific group",
		DefaultStatus: http.StatusOK,
	}, h.GetEnrollmentsByGroupID)

	huma.Register(g, huma.Operation{
		OperationID:   "get-enrollments-by-student",
		Method:        http.MethodGet,
		Path:          "/student/{student_id}",
		Summary:       "Get all enrollments for a student",
		Description:   "Get all enrollments for a specific student",
		DefaultStatus: http.StatusOK,
	}, h.GetEnrollmentsByStudentID)
}

func (h *EnrollmentsHandler) CreateEnrollment(c context.Context, input *dto.CreateEnrollmentReq) (*dto.CreateEnrollmentRes, error) {
	enrollment, err := h.svc.CreateEnrollment(c, input.Body.StudentID, input.Body.GroupID, input.Body.Fee)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("student_id", enrollment.StudentID).Str("group_id", enrollment.GroupID).
		Float64("fee", enrollment.Fee).Any("created", enrollment.CreatedAt).
		Msg("Created enrollment")
	return &dto.CreateEnrollmentRes{
		Body: *h.svc.ModelToRes(enrollment),
	}, nil
}

func (h *EnrollmentsHandler) UpdateEnrollment(c context.Context, input *dto.UpdateEnrollmentReq) (*dto.UpdateEnrollmentRes, error) {
	m := models.Enrollments{StudentID: input.Body.StudentID, GroupID: input.Body.GroupID}
	if input.Body.Fee != nil {
		m.Fee = *input.Body.Fee
	}
	enrollment, err := h.svc.UpdateEnrollment(c, m)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateEnrollmentRes{
		Body: *h.svc.ModelToRes(enrollment),
	}, nil
}

func (h *EnrollmentsHandler) GetEnrollmentByID(c context.Context, input *dto.GetEnrollmentByIDReq) (*dto.GetEnrollmentByIDRes, error) {
	enrollment, err := h.svc.GetEnrollmentByID(c, input.StudentID, input.GroupID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("student_id", input.StudentID).Str("group_id", input.GroupID).
		Int("created_at", enrollment.CreatedAt).
		Msg("Get enrollment by ID")
	return &dto.GetEnrollmentByIDRes{
		Body: *enrollment,
	}, nil
}

func (h *EnrollmentsHandler) DeleteEnrollment(c context.Context, input *dto.DeleteEnrollmentReq) (*dto.DeleteEnrollmentRes, error) {
	if err := h.svc.DeleteEnrollment(c, input.StudentID, input.GroupID); err != nil {
		return nil, err
	}
	return &dto.DeleteEnrollmentRes{
		Body: dto.DeleteEnrollmentResBody{
			StudentID: input.StudentID,
			GroupID:   input.GroupID,
		},
	}, nil
}

func (h *EnrollmentsHandler) ListEnrollments(c context.Context, input *dto.ListEnrollmentsReq) (*dto.ListEnrollmentsRes, error) {
	return h.svc.GetEnrollments(c, input)
}

func (h *EnrollmentsHandler) GetEnrollmentsByGroupID(c context.Context, input *dto.GetEnrollmentsByGroupIDReq) (*dto.GetEnrollmentsByGroupIDRes, error) {
	result, err := h.svc.GetEnrollmentsByGroupID(c, input)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("group_id", input.GroupID).Int("count", len(result.Body.Enrollments)).
		Msg("Get enrollments by group ID")
	return result, nil
}

func (h *EnrollmentsHandler) GetEnrollmentsByStudentID(c context.Context, input *dto.GetEnrollmentsByStudentIDReq) (*dto.GetEnrollmentsByStudentIDRes, error) {
	result, err := h.svc.GetEnrollmentsByStudentID(c, input)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("student_id", input.StudentID).Int("count", len(result.Body.Enrollments)).
		Msg("Get enrollments by student ID")
	return result, nil
}
