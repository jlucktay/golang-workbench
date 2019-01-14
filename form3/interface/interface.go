package main

type myStorage interface {
	Create(new widget) error
	Read(name string) (widget, error)
	Update(name string, data uint64) error
	Delete(name string) error
}
