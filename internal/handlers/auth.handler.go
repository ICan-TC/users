package handlers

import (
	"context"
	"net/http"
	"strings"

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

	huma.Register(g, huma.Operation{
		OperationID:   "verify",
		Method:        http.MethodPost,
		Path:          "/verify",
		Summary:       "Verify",
		Description:   "Verify a Token",
		DefaultStatus: http.StatusOK,
	}, h.Verify)
}

func (h *AuthHandler) Login(c context.Context, input *dto.LoginReq) (*dto.LoginRes, error) {
	_, t, err := h.svc.Login(c, *input.Body.Username, input.Body.Password)
	if err != nil {
		return nil, err
	}
	return &dto.LoginRes{
		Body: dto.Tokens{
			AccessAndExp: dto.AccessAndExp{
				AccessToken:          t.AccessToken.String(),
				AccessTokenExpiresAt: uint(t.AccessExp.Unix()),
			},
			RefreshAndExp: dto.RefreshAndExp{
				RefreshToken:          t.RefreshToken.String(),
				RefreshTokenExpiresAt: uint(t.RefreshExp.Unix()),
			},
		},
	}, nil
}

func (h *AuthHandler) Signup(c context.Context, input *dto.SignupReq) (*dto.SignupRes, error) {
	_, t, err := h.svc.Signup(c, input.Body.Email, input.Body.Username, input.Body.Password)
	if err != nil {
		return nil, err
	}
	return &dto.SignupRes{
		Body: dto.Tokens{
			AccessAndExp: dto.AccessAndExp{
				AccessToken:          t.AccessToken.String(),
				AccessTokenExpiresAt: uint(t.AccessExp.Unix()),
			},
			RefreshAndExp: dto.RefreshAndExp{
				RefreshToken:          t.RefreshToken.String(),
				RefreshTokenExpiresAt: uint(t.RefreshExp.Unix()),
			},
		},
	}, nil
}

func (h *AuthHandler) Refresh(c context.Context, input *dto.RefreshReq) (*dto.RefreshRes, error) {
	token, exp, err := h.svc.Refresh(c, input.Body.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &dto.RefreshRes{
		Body: dto.AccessAndExp{
			AccessToken:          token,
			AccessTokenExpiresAt: exp,
		},
	}, nil
}
func (h *AuthHandler) Logout(c context.Context, input *dto.LogoutReq) (*dto.LogoutRes, error) {
	if input.Authorization == "" {
		return &dto.LogoutRes{Body: struct{}{}}, huma.NewError(http.StatusBadRequest, "Authorization header is required")
	}
	sp := strings.Split(input.Authorization, " ")
	if len(sp) != 2 {
		return &dto.LogoutRes{Body: struct{}{}}, huma.NewError(http.StatusBadRequest, "Authorization header is invalid")
	}
	if sp[0] != "Bearer" {
		return &dto.LogoutRes{Body: struct{}{}}, huma.NewError(http.StatusBadRequest, "Authorization header is invalid")
	}
	token := sp[1]
	err := h.svc.Logout(c, token)
	return &dto.LogoutRes{Body: struct{}{}}, err
}

func (h *AuthHandler) Verify(c context.Context, input *dto.VerifyReq) (*dto.VerifyRes, error) {
	claims, err := h.svc.Verify(c, input.Body.Token)
	if err != nil {
		return nil, err
	}
	return &dto.VerifyRes{Body: claims}, err
}
