package main

import (
	uuid "github.com/satori/go.uuid"
)

type inMemoryStorage struct {
	store map[uuid.UUID]Payment
}

func (ims *inMemoryStorage) Init() error {
	ims.store = make(map[uuid.UUID]Payment)
	return nil
}

func (ims *inMemoryStorage) Create(p Payment) (uuid.UUID, error) {
	newId, errNew := uuid.NewV4()
	if errNew != nil {
		return uuid.Nil, errNew
	}
	ims.store[newId] = p
	return newId, nil
}

func (ims *inMemoryStorage) Read(id uuid.UUID) (Payment, error) {
	if p, exists := ims.store[id]; exists {
		return p, nil
	}
	return Payment{}, &NotFoundError{id}
}

func (ims *inMemoryStorage) Update(id uuid.UUID, p Payment) error {
	if _, exists := ims.store[id]; exists {
		ims.store[id] = p
		return nil
	}
	return &NotFoundError{id}
}

func (ims *inMemoryStorage) Delete(id uuid.UUID) error {
	if _, exists := ims.store[id]; exists {
		delete(ims.store, id)
		return nil
	}
	return &NotFoundError{id}
}
