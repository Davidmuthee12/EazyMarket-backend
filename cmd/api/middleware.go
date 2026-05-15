package main

import (
	"context"
	"fmt"
	"net"
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
		user, err = app.store.Users.GetByUUID(ctx, userUUID)
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// RequireRole secures routes by enforcing a minimum role precedence.
// It expects AuthTokenMiddleware to run before it and set the user in context.
func (app *application) RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := getUserFromCtx(r)
			if user == nil {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized"))
				return
			}

			allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
			if err != nil {
				app.internalServerError(w, r, err)
				return
			}

			if !allowed {
				app.forbiddenResponse(w, r)
				return
			}

			if user.Status == "suspended" {
				app.forbiddenResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) RequireActiveVendor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		if user == nil {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized"))
			return
		}

		vendor, err := app.store.Vendor.GetVendorByUUID(r.Context(), user.UUID)
		if err != nil {
			app.notFoundResponse(w, r, err)
			return
		}

		if vendor.Status != "approved" {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) StorefrontMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subdomain := app.storefrontSubdomain(r)
		if subdomain == "" {
			app.notFoundResponse(w, r, fmt.Errorf("storefront subdomain is missing"))
			return
		}

		vendor, err := app.store.Vendor.GetVendorBySubdomain(r.Context(), subdomain)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), storefrontVendorCtx, vendor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) storefrontSubdomain(r *http.Request) string {
	if subdomain := strings.TrimSpace(r.Header.Get("X-Store-Subdomain")); subdomain != "" {
		return strings.ToLower(subdomain)
	}

	if subdomain := strings.TrimSpace(r.URL.Query().Get("store")); subdomain != "" {
		return strings.ToLower(subdomain)
	}

	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	host = strings.ToLower(host)
	if host == "" || host == "localhost" || strings.HasPrefix(host, "127.") || strings.HasPrefix(host, "[::1]") {
		return ""
	}

	parts := strings.Split(host, ".")
	if len(parts) < 3 {
		return ""
	}

	subdomain := parts[0]
	if subdomain == "www" || subdomain == "api" {
		return ""
	}

	return subdomain
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.config.rateLimiter.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		ip := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			ip = host
		}

		if allow, retryAfter := app.ratelimiter.Allow(ip); !allow {
			app.rateLimitExceededResponse(w, r, retryAfter.String())
			return
		}

		next.ServeHTTP(w, r)
	})
}
