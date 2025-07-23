package repository

import (
	"errors"
	"fmt"
)

var (
	ErrStorage = errors.New("storage error")
)

func NewStorageError(operation, message string) error {
	return fmt.Errorf("%w: %s failed: %s", ErrStorage, operation, message)
}
