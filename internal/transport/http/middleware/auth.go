package middleware

import (
	"context"
	"net/http"

	"realTime/internal/domain"
	"uuid"
)

type authService interface {
	ValueidateSession(context.Context, uuid.UUID) (domain.Session, error)
}

type middleWare struct {
	authSvc authService
}

func NewMiddleware(authSvc authService) *middleWare {
	return &middleWare{authSvc: authSvc}
}

func (m *middleWare) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session-id")
		if err != nil {
			http.Error(w, "missing session-id", http.StatusUnauthorized)
			return
		}
		sessionUUID, err := uuid.Parse(c.Value)
		if err != nil {
			http.Error(w, "failed to validate session id", http.StatusUnauthorized)
			return
		}
		session, err := m.authSvc.ValueidateSession(r.Context(), sessionUUID)
		if err != nil {
			http.Error(w, "failed to validate session id ", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
