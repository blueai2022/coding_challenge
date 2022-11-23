package dataloader

import (
	"errors"
	"fmt"

	"github.com/blucv2022/crowdstats/models"
)

const (
	numberOfFields = 3
	firstNameField = "fname"
	lastNameField  = "lname"
	ageField       = "age"
)

// Different types of error returned by the Load function
var (
	ErrFileNotFound         = errors.New("cannot find data file")
	ErrInvalidDataFormat    = errors.New("invalid data format")
	ErrInvalidDataValue     = errors.New("invalid data value")
	ErrUnexpectedParseError = errors.New("unexpected parse error")
	ErrFileLoadFailure      = errors.New("cannot load data file")
)

type Loader interface {
	// Load(url string) ([]*models.Person, error)
	LoadDigest(url string) (*models.DataDigest, error)
}

type LoadError struct {
	Err error // The actual error
}

func (e *LoadError) Error() string {
	return fmt.Sprintf("load error: %v", e.Err)
}

func (e *LoadError) Unwrap() error {
	return e.Err
}
