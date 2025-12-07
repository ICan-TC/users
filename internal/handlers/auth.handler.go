package handlers

import (
	"context"
	"net/http"

	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	svc *service.AuthService
	log zerolog.Logger
}

func RegisterAuthRoutes(api huma.API, svc *service.AuthService) {
	h := &AuthHandler{svc: svc}
	g := huma.NewGroup(api, "/auth")
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Auth"}
	})

	huma.Register(g, huma.Operation{
		OperationID:   "login",
		Method:        http.MethodPost,
		Path:          "/login",
		Summary:       "Login",
		Description:   "Login",
		DefaultStatus: http.StatusOK,
	}, h.Login)

	huma.Register(g, huma.Operation{
		OperationID:   "signup",
		Method:        http.MethodPost,
		Path:          "/signup",
		Summary:       "Signup",
		Description:   "Create a new Account",
		DefaultStatus: http.StatusOK,
	}, h.Signup)

	huma.Register(g, huma.Operation{
		OperationID:   "refresh",
		Method:        http.MethodPost,
		Path:          "/refresh",
		Summary:       "Refresh",
		Description:   "Refresh your Access Token",
		DefaultStatus: http.StatusOK,
	}, h.Refresh)

	huma.Register(g, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodPost,
		Path:          "/logout",
		Summary:       "Logout",
		Description:   "Logout and revoke your Tokens",
		DefaultStatus: http.StatusOK,
	}, h.Logout)
}

func (h *AuthHandler) Login(c context.Context, input *dto.LoginReq) (*dto.LoginRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
func (h *AuthHandler) Signup(c context.Context, input *dto.SignupReq) (*dto.SignupRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}

func (h *AuthHandler) Refresh(c context.Context, input *dto.RefreshReq) (*dto.RefreshRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
func (h *AuthHandler) Logout(c context.Context, input *dto.LogoutReq) (*dto.LogoutRes, error) {
	return nil, huma.Error501NotImplemented("not implemented")
}
