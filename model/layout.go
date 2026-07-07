package model




type RegisterModel struct{
	NickName string
	FirstName string
	LastName string
	Email string
	Password string
}

type LoginModel struct {
	EmailOrNickName string `json:"emailOrNickName"`
	Password         string `json:"password"`
}

type SessionModel struct{
	
}
