package domain


type User struct {
	ID []byte
	NickName string
	LastName string	
	FirstName string
	Email string
	Age uint8
}



type RepoUser interface {
	Get(context.Context, []byte) (User, error)
	Create(context.Context, User) (User, error)
	GetUsers(context.Context) ([]User, error)
}
