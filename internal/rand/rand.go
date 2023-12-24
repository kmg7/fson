package rand

import "github.com/google/uuid"

func UUID() (string, error) {
	if id, err := uuid.NewRandom(); err != nil {
		return "", err
	} else {
		return id.String(), nil
	}
}
