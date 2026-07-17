package middleware 

import (
	"fmt"
	"sync"
	"context"
	"net/http"
	"realTime/internal/domain"
)



type userService interface {
	
}

type authService interface {
	ValidateSession(context.Context, string) (domain.Session, error)
}

type middleWare struct {
	authSvc  authService	
	routes map[string]string
	rateLimiterMap map[string]client
	rateLimiterMutex sync.RWMutex 
}


/*type route  struct {
	path string
	method string
}*/

type client struct {
	path string	
	limitedAt int64
}


func NewMiddleware(authSvc authService) middleWare {
		return middleWare{authSvc: authSvc}
}

func (m middleWare) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
			sessionID := r.Header.Get("x-session-id") 

			if sessionID == "" {
				http.Error(w, "missing session-id", http.StatusUnauthorized)
				return
			}
			session, err := m.authSvc.ValidateSession(r.Context(), sessionID)   			
			if err != nil {
				fmt.Println(err)
				http.Error(w, "failed to validate session id ", http.StatusUnauthorized)
				return
			}	
			ctx := context.WithValue(context.Background(), "user_id", session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
	}  ) 
}


