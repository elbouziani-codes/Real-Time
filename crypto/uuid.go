package crypto 

import (
		"github.com/gofrs/uuid/v5"
)

var u1 = uuid.Must(uuid.NewV4())

func GenerateUUID() (string, error) {
	u2, err := uuid.NewV4()	
	if err != nil {
		return "", err
	}
	return u2.String(), nil
}
