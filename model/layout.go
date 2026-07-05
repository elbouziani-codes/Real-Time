package model




type RegisterModel struct{
	Nick_name string
	First_name string
	Last_name string
	Email string
	Password_hash string
}

type LoginModel struct {
	EmailOrNickName string `json:"emailOrNickName"`
	Password         string `json:"password"`
}
type SessionModel struct{
	
}