package id

import "github.com/google/uuid"

type ID struct{ uuid.UUID }

func NewID(s string) (ID, error) {
	uuid, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}

	return ID{uuid}, nil
}

func MustParse(s string) ID {
	return ID{uuid.MustParse(s)}
}

func Generate() ID {
	return ID{uuid.Must(uuid.NewRandom())}
}
