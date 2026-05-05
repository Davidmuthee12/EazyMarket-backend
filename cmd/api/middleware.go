package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userUUID, ok := claims["sub"].(string)
		if !ok {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid token subject"))
			return
		}

		if _, err := uuid.Parse(userUUID); err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.getUser(ctx, userUUID)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUser(ctx context.Context, userUUID string) (*store.User, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetByUUID(ctx, userUUID)
	}

	app.logger.Infow("Cache hit", "key", "user", "id", userUUID)

	user, err := app.cacheStorage.Users.Get(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		app.logger.Infow("Fetching from DB", "id", userUUID)
		user, err := app.store.Users.GetByUUID(ctx, userUUID)
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
