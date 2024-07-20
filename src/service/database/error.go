package database

import "fmt"

type ConflictError struct{}

func (e *ConflictError) Error() string {
	return "attempted to create a record with an existing key"
}

type NotFoundError struct {
	Entity string
	ID     string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("unable to find %s with id %s", e.Entity, e.ID)
}
