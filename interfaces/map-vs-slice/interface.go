package main

type myStorage interface { //nolint:unused // The purpose of this exercise.
	Create(new widget) error
	Read(name string) (widget, error)
	Update(name string, data uint64) error
	Delete(name string) error
}
