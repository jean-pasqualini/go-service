package handlers

import (
	"context"
	"fmt"
	"github.com/jean-pasqualini/go-service/business/auth"
	"github.com/jean-pasqualini/go-service/business/data/user"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

type userController struct {
	user user.User
	auth *auth.Auth
}

func (uc userController) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.userController.query")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web values missing from context")
	}

	params := web.Params(r)
	_, err := strconv.Atoi(params["page"])
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid page format: %s", params["page"]), http.StatusBadRequest)
	}
	_, err = strconv.Atoi(params["rows"])
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid rows format: %s", params["rows"]), http.StatusBadRequest)
	}

	users, err := uc.user.Query(ctx, v.TraceId)
	if err != nil {
		return errors.Wrap(err, "unable to query for users")
	}

	return web.Respond(ctx, w, users, http.StatusOK)
}

func (uc userController) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web values missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	usr, err := uc.user.QueryByID(ctx, v.TraceId, claims, params["id"])
	if err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

func (uc userController) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web values missing from context")
	}

	var nu user.NewUser
	if err := web.Decode(r, &nu); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	usr, err := uc.user.Create(ctx, v.TraceId, nu, v.Now)
	if err != nil {
		return errors.Wrapf(err, "User: %+v", &usr)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

func (uc userController) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web values missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var uu user.UpdateUser
	if err := web.Decode(r, &uu); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	params := web.Params(r)
	err := uc.user.Update(ctx, v.TraceId, claims, params["id"], uu, v.Now)
	if err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s, User: %+v", params["id"], &uu)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (uc userController) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web values missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	err := uc.user.Delete(ctx, v.TraceId, claims, params["id"])
	if err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (uc userController) token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web values missing from context")
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := uc.user.Authenticate(ctx, v.TraceId, v.Now, email, pass)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating")
		}
	}

	params := web.Params(r)

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = uc.auth.GenerateToken(params["kid"], claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
