package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type PaymentStorage interface {
	Init() error
	Create(Payment) (uuid.UUID, error)
	Read(uuid.UUID) (Payment, error)
	Update(uuid.UUID, Payment) error
	Delete(uuid.UUID) error
}

type NotFoundError struct {
	id uuid.UUID
}

func (re *NotFoundError) Error() string {
	return fmt.Sprintf("Payment ID '%s' not found.", re.id)
}
