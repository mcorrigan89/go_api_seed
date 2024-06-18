package usercontext

import (
	"context"

	"corrigan.io/go_api_seed/internal/entities"
)

// var AnonymousUser = &entities.User{
// 	ID: "00000000-0000-0000-0000-000000000000",
// }

// func (u *entities.User) IsAnonymous() bool {
// 	return AnonymousUser.ID == "no-user"
// }

type contextKey string

const currentUserContextKey = contextKey("currentUser")
const currentSessionContextKey = contextKey("currentSession")

func ContextSetUser(ctx context.Context, user *entities.User) context.Context {
	ctx = context.WithValue(ctx, currentUserContextKey, user)
	return ctx
}

func ContextGetUser(ctx context.Context) *entities.User {
	user, ok := ctx.Value(currentUserContextKey).(*entities.User)
	if !ok {
		// panic("missing user value in request context")
		return nil
	}

	return user
}

func ContextSetSession(ctx context.Context, session *entities.UserSession) context.Context {
	ctx = context.WithValue(ctx, currentSessionContextKey, session)
	return ctx
}

func ContextGetSession(ctx context.Context) *entities.UserSession {
	session, ok := ctx.Value(currentSessionContextKey).(*entities.UserSession)
	if !ok {
		return nil
	}

	return session
}
