package id

import "github.com/google/uuid"

type ID string

func (id ID) String() string {
	return string(id)
}

func NewID(s string) (ID, error) {
	uuid, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}

	return ID(uuid.String()), nil
}

func GenerateID() ID {
	return ID(uuid.Must(uuid.NewRandom()).String())
}
