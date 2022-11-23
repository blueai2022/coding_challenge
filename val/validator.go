package val

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/blucv2022/crowdstats/models"
)

const (
	nameCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ '-."
)

var (
	nameCharsCheck [256]bool

// regex approach commented out: use name chars check outperforms it by far
// isValidPersonName = regexp.MustCompile(`^[a-zA-Z'\-\s]+$`).MatchString
)

func init() {
	for i := 0; i < len(nameCharSet); i++ {
		nameCharsCheck[nameCharSet[i]] = true
	}
}

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(strings.TrimSpace(value))
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidatePersonName(value string) error {
	if err := ValidateString(value, 1, 100); err != nil {
		return err
	}

	//skip validation for non-english name in unicode
	if len(value) != utf8.RuneCountInString(value) {
		return nil
	}

	for i := 0; i < len(value); i++ {
		if !nameCharsCheck[value[i]] {
			return fmt.Errorf("must contain letters, dots, -, ', or space: %s", value)
		}
	}

	// regex approach commented out: use name chars check outperforms it by far
	// if !isValidPersonName(value) {
	// 	return fmt.Errorf("must contain letters, -, ', or space: %s", value)
	// }

	return nil
}

func ValidateAge(age int) error {
	if age > 200 || age < 0 {
		return fmt.Errorf("must be a valid age: %d", age)
	}
	return nil
}

func ValidatePerson(person *models.Person) error {
	if err := ValidateAge(person.Age); err != nil {
		return err
	}

	if err := ValidatePersonName(person.FirstName); err != nil {
		return err
	}

	if err := ValidatePersonName(person.LastName); err != nil {
		return err
	}

	return nil
}
