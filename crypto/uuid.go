package crypto

import (
	"github.com/gofrs/uuid/v5"
)

var Nil = UUID{Value: uuid.Nil}

// cant absracting using interface beacuse I would need tp scan value using sql
type UUID struct {
	Value uuid.UUID
}

func GenerateUUID() (UUID, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return UUID{}, err
	}
	return UUID{Value: u}, nil
}

// may be not needed
func ParseUUID(s string) (UUID, error) {
	u := uuid.Nil
	err := u.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID{Value: u}, nil
}
