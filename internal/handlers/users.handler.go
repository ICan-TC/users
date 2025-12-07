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
	g.UseMiddleware(middleware.AuthMiddleware)

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
		Path:          "/{field}/{value}",
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
	user, err := h.svc.CreateUser(c, input.Body.Email, input.Body.Username, input.Body.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	h.log.Info().Str("username", input.Body.Username).Str("email", input.Body.Email).Str("id", user.UserID).Any("created", user.CreatedAt).
		Msg("Created user")
	return &dto.CreateUserRes{
		Body: dto.CreateUserResBody{
			ID: user.UserID, Username: user.Username, Email: user.Email, CreatedAt: int(user.CreatedAt.Unix()), UpdatedAt: int(user.CreatedAt.Unix()),
		},
	}, nil
}

func (h *UsersHandler) UpdateUser(c context.Context, input *dto.UpdateUserReq) (*dto.UpdateUserRes, error) {
	m := models.Users{UserID: input.Body.Id}
	if input.Body.Username != nil {
		m.Username = *input.Body.Username
	}
	if input.Body.Email != nil {
		m.Email = *input.Body.Email
	}
	if input.Body.Password != nil {
		m.PasswordHash = *input.Body.Password
	}
	if input.Body.FirstName != nil {
		m.FirstName = input.Body.FirstName
	}
	if input.Body.FamilyName != nil {
		m.FamilyName = input.Body.FamilyName
	}
	if input.Body.PhoneNumber != nil {
		m.PhoneNumber = input.Body.PhoneNumber
	}
	if input.Body.DateOfBirth != nil {
		m.DateOfBirth = input.Body.DateOfBirth
	}
	user, err := h.svc.UpdateUser(c, m)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateUserRes{
		Body: dto.UpdateUserResBody{
			ID:          user.UserID,
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			FamilyName:  user.FamilyName,
			PhoneNumber: user.PhoneNumber,
			DateOfBirth: user.DateOfBirth,
			CreatedAt:   int(user.CreatedAt.Unix()),
			UpdatedAt:   int(user.UpdatedAt.Unix()),
		},
	}, nil
}

func (h *UsersHandler) GetUserByField(c context.Context, input *dto.GetUserByFieldReq) (*dto.GetUserByFieldRes, error) {
	a, err := h.svc.GetUserByField(c, input.Field, input.Value)
	if err != nil {
		return nil, err
	}
	return &dto.GetUserByFieldRes{
		Body: dto.GetUserResBody{
			ID:       a.UserID,
			Username: a.Username,
			Email:    a.Email,
		},
	}, nil
}
func (h *UsersHandler) GetUserByID(c context.Context, input *dto.GetUserByIDReq) (*dto.GetUserByIDRes, error) {
	user, err := h.svc.GetUserByID(c, input.ID)
	if err != nil {
		return nil, err
	}
	h.log.Info().Str("id", input.ID).Str("username", user.Username).Str("email", user.Email).Any("created", user.CreatedAt).
		Msg("Get user by ID")
	return &dto.GetUserByIDRes{
		Body: dto.GetUserResBody{
			ID: user.UserID, Username: user.Username, Email: user.Email, CreatedAt: int(user.CreatedAt.Unix()), UpdatedAt: int(user.CreatedAt.Unix()),
		},
	}, nil
}
func (h *UsersHandler) DeleteUser(c context.Context, input *dto.DeleteUserReq) (*dto.DeleteUserRes, error) {
	if err := h.svc.DeleteUser(c, input.ID); err != nil {
		return nil, err
	}
	return &dto.DeleteUserRes{
		Body: dto.DeleteUserResBody{
			ID: input.ID,
		},
	}, nil
}
func (h *UsersHandler) ListUsers(c context.Context, input *dto.ListUsersReq) (*dto.ListUsersRes, error) {
	return h.svc.GetUsers(c, input)
}
