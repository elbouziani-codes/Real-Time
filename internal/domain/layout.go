package domain



type LoginModel struct {
	EmailOrNickName string `json:"emailOrNickName"`
	Password         string `json:"password"`
}

type SessionModel struct {
}
