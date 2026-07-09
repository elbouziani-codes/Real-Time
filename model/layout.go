package model


type RegisterModel struct {
	NickName  string `json:"nickName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginModel struct {
	EmailOrNickName string `json:"emailOrNickName"`
	Password         string `json:"password"`
}

type SessionModel struct {
}
