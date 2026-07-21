package middleware

import (
	"context"
	"fmt"
	"net/http"
	"realTime/crypto"
	"realTime/internal/domain"
	"sync"
)

type userService interface {
}

type authService interface {
	ValueidateSession(context.Context, crypto.UUID) (domain.Session, error)
}

type middleWare struct {
	authSvc          authService
	routes           map[string]string
	rateLimiterMap   map[string]client
	rateLimiterMutex sync.RWMutex
}

/*type route  struct {
	path string
	method string
}*/

type client struct {
	path      string
	limitedAt int64
}

func NewMiddleware(authSvc authService) middleWare {
	return middleWare{authSvc: authSvc}
}

func (m middleWare) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session-id")
		if err != nil {
			http.Error(w, "missing session-id", http.StatusUnauthorized)
			return
		}
		fmt.Println(err, c)
		sessionUUID, err := crypto.ParseUUID(c.Value)
		if err != nil {
			http.Error(w, "failed to validate session id", http.StatusUnauthorized)
			return
		}
		session, err := m.authSvc.ValueidateSession(r.Context(), sessionUUID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "failed to validate session id ", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(context.Background(), "user_id", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
