package utils

import "github.com/google/uuid"

func IsValidUUIDv4(id uuid.UUID) bool {
	return id.Version() == 4 && id != uuid.Nil
}

func ParseUUIDFromString(str string)( uuid.UUID, error){
	id, err := uuid.Parse(str)
	if err != nil{
		return uuid.Nil, err
	}
	return id, nil
}