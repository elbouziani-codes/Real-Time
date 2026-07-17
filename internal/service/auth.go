package service

import (
	"fmt"
	"time"
	"errors"
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
)

type AuthService struct {
	authRepo AuthRepo
	userRepo UserRepo
}
type AuthRepo interface {
	DeleteSession(context.Context, string) error 
	SaveSession(context.Context, string, string) error
	GetByID(context.Context, string)  (domain.Session, error)
}
/*type UserRepo interface {
	GetByID(context.Context, string) (domain.User, error)
}*/

func NewAuthService(authRepo AuthRepo, UserRepo UserRepo) *AuthService {
	return &AuthService{authRepo: authRepo, userRepo: UserRepo}
}

func (a *AuthService) CreateSession(ctx context.Context, userID string) (string, error) {
	id, err := crypto.GenerateUUID()
	if err != nil {
		return "", err
	}

	err = a.authRepo.SaveSession(ctx, userID, id)
	if err != nil {
		return "", err
	}

	return id, nil
}


func (a *AuthService) ValidateSession(ctx context.Context, sessionID string) (domain.Session, error) {
	session, err := a.authRepo.GetByID(ctx, sessionID) 
	if err != nil {
		fmt.Println("exit", sessionID)
		return domain.Session{}, err
	}
	if time.Now().Unix() > session.ExpireAt  {
			err :=  a.Logout(ctx, session.UserID)
			if err != nil {
				return domain.Session{}, err
			}
			return domain.Session{}, domain.Error{Message: "invalid credentials", Code: domain.UnauthorizedCode}
	} 
	return session, nil
}


func (a *AuthService) Logout(ctx context.Context, userID string) (error) {
	err := a.authRepo.DeleteSession(ctx, userID) 
	if err != nil  {
		var buckErr domain.Error 
		if errors.As(err, buckErr) {
			if buckErr.Code == domain.NotFoundCode {
					return nil
			}	
		} 
		return err
	}
	return nil
}


func (a *AuthService) Login(ctx context.Context, creds domain.Credentials, idType int) (string, error) {
	// phas_1 finding target user
	var userID string	
	var hashedPassword string
	if (idType == domain.EmailType) {
		user, err := a.userRepo.GetByEmail(ctx, creds.Identifier)
		if err != nil {
			return "", domain.Error{Message: "invalid credentials", Code: domain.UnauthorizedCode}

		}
		userID = user.ID
		hashedPassword = user.Password
	} else {
		user, err := a.userRepo.GetByNickName(ctx, creds.Identifier) 	
		if err != nil {
			return "",  domain.Error{Message: "invalid credentials", Code: domain.UnauthorizedCode}

		}
		userID = user.ID
		hashedPassword = user.Password
	}	

	// phase_2 matching password
	err := crypto.CompareHashWithPassword(hashedPassword, creds.Password)	
	if err != nil {
			return "", domain.Error{Message: "invalid credentials", Code: domain.UnauthorizedCode}
	}

	// clearing old session
	err = a.Logout(ctx, userID) 
	if err != nil   {
		return "", err
	}

	//spawnning new one 
	sessionID, err := a.CreateSession(ctx, userID)
	if err != nil {
		return "", err 
	}
	return sessionID, nil
}
