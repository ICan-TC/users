package handlers

import (
	"context"
	"net/http"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/middleware"
	"github.com/ICan-TC/users/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
)

type StudentParentsHandler struct {
	svc *service.StudentParentsService
	log zerolog.Logger
}

func RegisterStudentParentsRoutes(api huma.API, svc *service.StudentParentsService) {
	h := &StudentParentsHandler{svc: svc, log: logging.L()}
	g := huma.NewGroup(api, "/student-parents")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"StudentParents"}
	})
	g.UseMiddleware(middleware.AuthMiddleware)

	huma.Register(g, huma.Operation{
		OperationID:   "create-student-parent",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Link a student to a parent",
		Description:   "Link a student to a parent",
		DefaultStatus: http.StatusCreated,
	}, h.CreateStudentParent)

	huma.Register(g, huma.Operation{
		OperationID:   "update-student-parent",
		Method:        http.MethodPatch,
		Path:          "",
		Summary:       "Update a student-parent relationship",
		Description:   "Update a student-parent relationship",
		DefaultStatus: http.StatusOK,
	}, h.UpdateStudentParent)

	huma.Register(g, huma.Operation{
		OperationID:   "get-student-parent",
		Method:        http.MethodGet,
		Path:          "/{student_id}/{parent_id}",
		Summary:       "Get a student-parent relationship",
		Description:   "Get a student-parent relationship",
		DefaultStatus: http.StatusOK,
	}, h.GetStudentParent)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-student-parent",
		Method:        http.MethodDelete,
		Path:          "/{student_id}/{parent_id}",
		Summary:       "Unlink a student from a parent",
		Description:   "Unlink a student from a parent",
		DefaultStatus: http.StatusOK,
	}, h.DeleteStudentParent)

	huma.Register(g, huma.Operation{
		OperationID:   "list-student-parents",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List student-parent relationships",
		Description:   "List student-parent relationships",
		DefaultStatus: http.StatusOK,
	}, h.ListStudentParents)
}

func (h *StudentParentsHandler) CreateStudentParent(c context.Context, input *dto.CreateStudentParentReq) (*dto.CreateStudentParentRes, error) {
	sp, err := h.svc.CreateStudentParent(c, input.Body.StudentID, input.Body.ParentID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("student_id", sp.StudentID).Str("parent_id", sp.ParentID).Any("created", sp.CreatedAt).
		Msg("Linked student to parent")
	return &dto.CreateStudentParentRes{
		Body: dto.CreateStudentParentResBody{
			StudentID: sp.StudentID,
			ParentID:  sp.ParentID,
			CreatedAt: int(sp.CreatedAt.Unix()),
			UpdatedAt: int(sp.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *StudentParentsHandler) UpdateStudentParent(c context.Context, input *dto.UpdateStudentParentReq) (*dto.UpdateStudentParentRes, error) {
	sp, err := h.svc.UpdateStudentParent(c, input.Body.StudentID, input.Body.ParentID, input.Body.StudentID, input.Body.ParentID)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateStudentParentRes{
		Body: dto.UpdateStudentParentResBody{
			StudentID: sp.StudentID,
			ParentID:  sp.ParentID,
			CreatedAt: int(sp.CreatedAt.Unix()),
			UpdatedAt: int(sp.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *StudentParentsHandler) GetStudentParent(c context.Context, input *dto.GetStudentParentReq) (*dto.GetStudentParentRes, error) {
	sp, err := h.svc.GetStudentParentByID(c, input.StudentID, input.ParentID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("student_id", input.StudentID).Str("parent_id", input.ParentID).Any("created", sp.CreatedAt).
		Msg("Get student-parent relationship")
	return &dto.GetStudentParentRes{
		Body: dto.GetStudentParentResBody{
			StudentID: sp.StudentID,
			ParentID:  sp.ParentID,
			CreatedAt: int(sp.CreatedAt.Unix()),
			UpdatedAt: int(sp.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *StudentParentsHandler) DeleteStudentParent(c context.Context, input *dto.DeleteStudentParentReq) (*dto.DeleteStudentParentRes, error) {
	if err := h.svc.DeleteStudentParent(c, input.StudentID, input.ParentID); err != nil {
		return nil, err
	}
	return &dto.DeleteStudentParentRes{
		Body: dto.DeleteStudentParentResBody{
			StudentID: input.StudentID,
			ParentID:  input.ParentID,
		},
	}, nil
}

func (h *StudentParentsHandler) ListStudentParents(c context.Context, input *dto.ListStudentParentsReq) (*dto.ListStudentParentsRes, error) {
	return h.svc.GetStudentParents(c, input)
}
