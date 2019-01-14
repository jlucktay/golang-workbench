package main

import (
	"fmt"
)

type dummyStoreMap struct {
	w map[string]uint64
}

func (dsm *dummyStoreMap) Create(new widget) error {
	if _, exists := dsm.w[new.name]; exists {
		return fmt.Errorf("widget with name '%s' already exists", new.name)
	}
	dsm.w[new.name] = new.data
	return nil
}

func (dsm *dummyStoreMap) Read(name string) (widget, error) {
	if existing, ok := dsm.w[name]; ok {
		return widget{name, existing}, nil
	}
	return widget{}, fmt.Errorf("nothing in store named '%s'", name)
}

func (dsm *dummyStoreMap) Update(name string, data uint64) error {
	if _, exists := dsm.w[name]; exists {
		dsm.w[name] = data
		return nil
	}
	return fmt.Errorf("nothing in store named '%s'", name)
}

func (dsm *dummyStoreMap) Delete(name string) error {
	if _, exists := dsm.w[name]; exists {
		delete(dsm.w, name)
		return nil
	}
	return fmt.Errorf("nothing in store named '%s'", name)
}
