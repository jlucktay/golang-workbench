package main

import (
	"fmt"
)

type dummyStoreSlice struct {
	w []widget
}

func (dss *dummyStoreSlice) Create(new widget) error {
	if dss.exists(new.name) >= 0 {
		return fmt.Errorf("widget with name '%s' already exists", new.name)
	}
	dss.w = append(dss.w, new)
	return nil
}

func (dss *dummyStoreSlice) Read(name string) (widget, error) {
	e := dss.exists(name)
	if e >= 0 {
		return dss.w[e], nil
	}
	return widget{}, fmt.Errorf("nothing in store named '%s'", name)
}

func (dss *dummyStoreSlice) Update(name string, data uint64) error {
	e := dss.exists(name)
	if e >= 0 {
		dss.w[e].data = data
		return nil
	}
	return fmt.Errorf("nothing in store named '%s'", name)
}

func (dss *dummyStoreSlice) Delete(name string) error {
	e := dss.exists(name)
	if e >= 0 {
		dss.w = append(dss.w[:e], dss.w[e+1:]...)
		return nil
	}
	return fmt.Errorf("nothing in store named '%s'", name)
}

func (dss *dummyStoreSlice) exists(name string) int {
	for index := range dss.w {
		if dss.w[index].name == name {
			return index
		}
	}
	return -1
}
