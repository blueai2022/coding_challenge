package dataloader

import (
	"log"
	"strconv"

	"github.com/blucv2022/crowdstats/models"
	"github.com/blucv2022/crowdstats/val"
)

func parseFieldPos(fields []string) (map[string]int, error) {
	fieldPos := make(map[string]int)

	for i, field := range fields {
		switch field {
		case firstNameField, lastNameField, ageField:
			fieldPos[field] = i
		default:
			log.Println("unexpected csv column name:", field)
			return nil, ErrInvalidDataFormat
		}
	}

	return fieldPos, nil
}

func parsePerson(values []string, fieldPos map[string]int) (*models.Person, error) {
	person := &models.Person{}

	pos, ok := fieldPos[firstNameField]
	if !ok || len(values[pos]) == 0 {
		return nil, ErrInvalidDataFormat
	}
	person.FirstName = values[pos]

	pos, ok = fieldPos[lastNameField]
	if !ok || len(values[pos]) == 0 {
		return nil, ErrInvalidDataFormat
	}
	person.LastName = values[pos]

	pos, ok = fieldPos[ageField]
	if !ok {
		return nil, ErrInvalidDataFormat
	}

	age, err := strconv.Atoi(values[pos])
	if err != nil {
		return nil, ErrInvalidDataFormat
	}
	person.Age = age

	if err = val.ValidatePerson(person); err != nil {
		return nil, ErrInvalidDataValue
	}

	return person, nil
}
