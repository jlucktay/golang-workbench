package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type PaymentStorage interface {
	Init() error
	Create(Payment) (uuid.UUID, error)
	Read(uuid.UUID) (Payment, ReadError)
	Update(uuid.UUID, Payment) error
	Delete(uuid.UUID) error
}

type ReadError struct {
	id uuid.UUID
}

func (re *ReadError) Error() string {
	return fmt.Sprintf("Payment ID '%s' not found.", re.id)
}
