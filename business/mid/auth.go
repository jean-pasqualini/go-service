package mid

import (
	"context"
	"github.com/jean-pasqualini/go-service/business/auth"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"
)

var ErrForbidden = web.NewRequestError(
	errors.New("you are not authorized for that action"),
	http.StatusForbidden,
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Parse the authorization header.
			// Expected header is of the format `Authorization: Bearer <token>`
			authHeader := r.Header.Get("Authorization")

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			// Validate the token is signed by us
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			// Add claims to the context so they can be retrieve later
			ctx = context.WithValue(ctx, auth.Key, claims)

			return handler(ctx, w, r)
		}
	}
}

// Authorize validates that an authenticated user has at least one role from specified list.
// This method constructs the actual function that is used.
func Authorize(log *log.Logger, roles ...string) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, return failure.
			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("claims missing from context")
			}

			if !claims.Authorize(roles...) {
				log.Printf("mid: authorize:\n \texpect: %+v \troles: %+v", roles, claims.Roles)
				return ErrForbidden
			}

			return handler(ctx, w, r)
		}
	}
}
