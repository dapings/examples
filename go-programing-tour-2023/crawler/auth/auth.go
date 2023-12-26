package auth

import (
	"context"
	"errors"
	"strings"

	"go-micro.dev/v4"
	"go-micro.dev/v4/auth"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
)

func NewAuthWrapper(service micro.Service) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// fetch metadata from context (request headers).
			md, b := metadata.FromContext(ctx)
			if !b {
				return errors.New("not found metadata")
			}

			// get auth header.
			authHeader, ok := md["Authorization"]
			if !ok || !strings.HasPrefix(authHeader, auth.BearerScheme) {
				return errors.New("no provided auth token")
			}

			// extract auth token.
			token := strings.TrimPrefix(authHeader, auth.BearerScheme)

			// extract account from token.
			a := service.Options().Auth
			if _, err := a.Inspect(token); err != nil {
				return errors.New("invalid auth token")
			}

			return h(ctx, req, rsp)
		}
	}
}
