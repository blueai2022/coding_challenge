package dataloader

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/blucv2022/crowdstats/models"
)

func Digest(reader io.Reader) (*models.DataDigest, error) {
	csvReader := csv.NewReader(reader)

	digest := &models.DataDigest{}

	i := 0
	var fieldPos map[string]int

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if _, ok := err.(*csv.ParseError); ok {
				return nil, ErrInvalidDataFormat
			}
			return nil, ErrUnexpectedParseError
		}

		if len(record) != numberOfFields {
			return nil, ErrInvalidDataFormat
		}

		//remove white spaces
		fields := make([]string, numberOfFields)
		for i := range record {
			fields[i] = strings.TrimSpace(record[i])
		}

		if i == 0 {
			fieldPos, err = parseFieldPos(fields)
			if err != nil {
				return nil, err
			}
		} else {
			person, err := parsePerson(fields, fieldPos)
			if err != nil {
				return nil, err
			}

			// add to digest
			digest.TotalAgeCounts++
			digest.AgeCounts[person.Age]++

			//only create person name string if necessary
			if len(digest.AgePersonName[person.Age]) == 0 {
				digest.AgePersonName[person.Age] =
					fmt.Sprintf("%s %s", person.FirstName, person.LastName)
			}
		}

		i++
	}

	return digest, nil
}
