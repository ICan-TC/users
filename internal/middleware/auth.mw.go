package middleware

import (
	"net/http"
	"strings"

	"github.com/ICan-TC/lib/tokens"
	"github.com/ICan-TC/users/internal/config"
	"github.com/danielgtaylor/huma/v2"
)

func AuthMiddleware(hc huma.Context, next func(huma.Context)) {
	ctx := hc.Context()
	h := hc.Header("Authorization")
	if h == "" {
		hc.SetStatus(http.StatusUnauthorized)
		return
	}
	splits := strings.Split(h, " ")
	if len(splits) != 2 {
		hc.SetStatus(http.StatusBadRequest)
		return
	}
	if splits[0] != "Bearer" {
		hc.SetStatus(http.StatusBadRequest)
		return
	}
	if splits[1] == "" {
		hc.SetStatus(http.StatusBadRequest)
		return
	}

	_, err := tokens.ParseToken(ctx, splits[1], config.Get().Auth.Secret, "access")
	if err != nil {
		hc.SetStatus(http.StatusUnauthorized)
		return
	}
	next(hc)
}
