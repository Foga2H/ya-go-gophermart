package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Foga2H/ya-go-gophermart/internal/auth"
	"github.com/Foga2H/ya-go-gophermart/internal/logger"
)

type Jwt struct {
	token  *auth.Token
	logger *logger.Logger
}

func NewJwt(token *auth.Token, logger *logger.Logger) *Jwt {
	return &Jwt{
		token:  token,
		logger: logger,
	}
}

type contextKey string

const userIDKey contextKey = "userID"

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok
}

func (j *Jwt) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userToken, ok := extractToken(r.Header.Get("Authorization"))
			if !ok {
				j.logger.Errorf("missing token")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			userID, err := j.token.ParseUserID(userToken)
			if err != nil {
				j.logger.Errorf("failed to parse token: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ContextWithUserID(r.Context(), userID)))
		})
	}
}

func extractToken(authorizationHeader string) (string, bool) {
	fields := strings.Fields(authorizationHeader)
	if len(fields) == 0 {
		return "", false
	}

	if len(fields) == 1 {
		return fields[0], true
	}

	if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") && fields[1] != "" {
		return fields[1], true
	}

	return "", false
}
