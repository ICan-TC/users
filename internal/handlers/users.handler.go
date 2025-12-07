package handlers

import (
	"context"
	"net/http"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
)

type UsersHandler struct {
	svc *service.UsersService
	log zerolog.Logger
}

func RegisterUsersRoutes(api huma.API, svc *service.UsersService) {
	h := &UsersHandler{svc: svc, log: logging.L()}
	g := huma.NewGroup(api, "/users")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Users"}
	})

	huma.Register(g, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create a user",
		Description:   "Create a user",
		DefaultStatus: http.StatusCreated,
	}, h.CreateUser)

	huma.Register(g, huma.Operation{
		OperationID:   "update-user",
		Method:        http.MethodPatch,
		Path:          "",
		Summary:       "Update a user",
		Description:   "Update a user",
		DefaultStatus: http.StatusOK,
	}, h.UpdateUser)

	huma.Register(g, huma.Operation{
		OperationID:   "get-user-by-field",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "Get a user by field",
		Description:   "Get a user by field",
		DefaultStatus: http.StatusOK,
	}, h.GetUserByField)

	huma.Register(g, huma.Operation{
		OperationID:   "get-user-by-id",
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Get a user by ID",
		Description:   "Get a user by ID",
		DefaultStatus: http.StatusOK,
	}, h.GetUserByID)

	huma.Register(g, huma.Operation{
		OperationID:   "delete-user",
		Method:        http.MethodDelete,
		Path:          "/{id}",
		Summary:       "Delete a user",
		Description:   "Delete a user",
		DefaultStatus: http.StatusOK,
	}, h.DeleteUser)

	huma.Register(g, huma.Operation{
		OperationID:   "list-users",
		Method:        http.MethodGet,
		Path:          "",
		Summary:       "List users",
		Description:   "List users",
		DefaultStatus: http.StatusOK,
	}, h.ListUsers)

}

func (h *UsersHandler) CreateUser(c context.Context, input *dto.CreateUserReq) (*dto.CreateUserRes, error) {
	user, err := h.svc.CreateUser(c, input.Body.Username, input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	h.log.Info().Str("username", input.Body.Username).Str("email", input.Body.Email).Str("id", user.ID).Any("created", user.CreatedAt).
		Msg("Created user")
	return &dto.CreateUserRes{
		Body: dto.CreateUserResBody{
			ID: user.ID, Username: user.Username, Email: user.Email, CreatedAt: int(user.CreatedAt.Unix()), UpdatedAt: int(user.CreatedAt.Unix()),
		},
	}, nil
}

func (h *UsersHandler) UpdateUser(c context.Context, input *dto.UpdateUserReq) (*dto.UpdateUserRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *UsersHandler) GetUserByField(c context.Context, input *dto.GetUserByFieldReq) (*dto.GetUserByFieldRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
func (h *UsersHandler) GetUserByID(c context.Context, input *dto.GetUserByIDReq) (*dto.GetUserByIDRes, error) {
	user, err := h.svc.GetUser(c, input.ID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("id", input.ID).Str("username", user.Username).Str("email", user.Email).Any("created", user.CreatedAt).
		Msg("Get user by ID")
	return &dto.GetUserByIDRes{
		Body: dto.GetUserResBody{
			ID: user.ID, Username: user.Username, Email: user.Email, CreatedAt: int(user.CreatedAt.Unix()), UpdatedAt: int(user.CreatedAt.Unix()),
		},
	}, nil
}
func (h *UsersHandler) DeleteUser(c context.Context, input *dto.DeleteUserReq) (*dto.DeleteUserRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
func (h *UsersHandler) ListUsers(c context.Context, input *dto.ListUsersReq) (*dto.ListUsersRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
