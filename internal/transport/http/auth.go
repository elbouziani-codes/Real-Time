package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"realTime/crypto"
	"realTime/internal/domain"
)

type AuthRepo interface {
	SaveSession(context.Context, crypto.UUID, crypto.UUID) error
}

type AuthService interface {
	Login(context.Context, domain.Credentials, int) (crypto.UUID, error)
	CreateSession(context.Context, crypto.UUID) (crypto.UUID, error)
	Logout(context.Context, crypto.UUID) error
}

type UserService interface {
	CreateUser(context.Context, *domain.User) error
	GetUser(context.Context, crypto.UUID, crypto.UUID) (*domain.UserProfile, error)
	GetUsers(context.Context, crypto.UUID, int, crypto.UUID) ([]domain.UserContact, error)
}

type AuthHandler struct {
	authSvc    AuthService
	userSvc    UserService
	middleware middleWare
}

type TestHandler struct {
}

func NewTestHandler(authSvc AuthService, userSvc UserService) TestHandler {
	return TestHandler{}
}

func (t *TestHandler) Test(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if err := json.NewEncoder(w).Encode(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewAuthHandler(authSvc AuthService, userSvc UserService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, userSvc: userSvc}
}


func (a *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)
	
	user, err := a.userSvc.GetUser(r.Context(), userID, userID)
	if err != nil {
		Error(err, w)
		return
	}
	
	
	
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}

// GetUsers lists the other users, ordered so the person this user talked with
// most recently comes first and people they have never messaged come last.
func (a *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)

	query := r.URL.Query()
	limit, cursor, err := domain.ValidateUserListPaging(query.Get("limit"), query.Get("cursor"))
	if err != nil {
		Error(err, w)
		return
	}

	users, err := a.userSvc.GetUsers(r.Context(), userID, limit, cursor)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	registerRequest := domain.RegisterRequest{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&registerRequest); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w) // must beh
		return
	}

	user, err := domain.ValueidateUserInfo(registerRequest)
	if err != nil {
		Error(err, w)
		return
	}

	err = a.userSvc.CreateUser(r.Context(), &user)
	if err != nil {
		Error(err, w)
		return
	}

	sessionID, err := a.authSvc.CreateSession(r.Context(), user.ID)
	if err != nil {
		Error(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	a.setCookie(w, sessionID)
	if err := json.NewEncoder(w).Encode("done"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)

	if err := a.authSvc.Logout(r.Context(), userID); err != nil {
		Error(err, w)
		return
	}

	a.clearCookie(w)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("done"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds domain.Credentials
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&creds); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w) // must beh
		return
	}
	idType, err := domain.LoginValueidation(creds)
	if err != nil {
		Error(err, w)
		//instead of repeating code I will add helper for that
		return
	}
	sessionID, err := a.authSvc.Login(r.Context(), creds, idType)
	if err != nil {
		Error(err, w)
		return
	}
	a.setCookie(w, sessionID)
	if err := json.NewEncoder(w).Encode("done"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// I may need to pass pointer
func (a *AuthHandler) setCookie(w http.ResponseWriter, sessionID crypto.UUID) {
	cookie := &http.Cookie{
		Name:     "session-id",
		Value:    sessionID.Value.String(),
		Path:     "/",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// clearCookie expires the session cookie. The attributes must mirror setCookie
// so the browser matches the cookie it already stored and drops it.
func (a *AuthHandler) clearCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session-id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
func Error(err error, w http.ResponseWriter) {
	var valErr domain.Error
	if errors.As(err, &valErr) {
		switch valErr.Code {
		case domain.ConflictCode:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		case domain.BadFormatCode:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case domain.UnauthorizedCode:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		case domain.NotFoundCode:
			http.Error(w, err.Error(), http.StatusNotFound)
			return

		default:
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
			return
		}
	}
}
