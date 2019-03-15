package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type AlreadyExistsError struct {
	id uuid.UUID
}

func (re *AlreadyExistsError) Error() string {
	return fmt.Sprintf("Payment ID '%s' already exists.", re.id)
}

type NotFoundError struct {
	id uuid.UUID
}

func (re *NotFoundError) Error() string {
	return fmt.Sprintf("Payment ID '%s' not found.", re.id)
}
