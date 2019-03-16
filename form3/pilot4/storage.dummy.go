package main

import (
	uuid "github.com/satori/go.uuid"
)

type inMemoryStorage struct {
	store map[uuid.UUID]Payment
}

func (ds *inMemoryStorage) Init() error {
	ds.store = make(map[uuid.UUID]Payment)
	return nil
}

func (ds *inMemoryStorage) Create(p Payment) (uuid.UUID, error) {
	newId, errNew := uuid.NewV4()
	if errNew != nil {
		return uuid.Nil, errNew
	}
	ds.store[newId] = p
	return newId, nil
}

func (ds *inMemoryStorage) Read(id uuid.UUID) (Payment, error) {
	if p, exists := ds.store[id]; exists {
		return p, nil
	}
	return Payment{}, &NotFoundError{id}
}

func (ds *inMemoryStorage) Update(id uuid.UUID, p Payment) error {
	if _, exists := ds.store[id]; exists {
		ds.store[id] = p
		return nil
	}
	return &NotFoundError{id}
}

func (ds *inMemoryStorage) Delete(id uuid.UUID) error {
	if _, exists := ds.store[id]; exists {
		delete(ds.store, id)
		return nil
	}
	return &NotFoundError{id}
}
