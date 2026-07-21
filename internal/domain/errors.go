package domain


const (
	UnexpectedCode = iota
	NotFoundCode
	BadFormatCode
	ConflictCode
	UnauthorizedCode
)

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string {
	return e.Message
}
