// The 'reflect' package in the standard library lets Go examine aspects of itself at runtime.
package main

import (
	"fmt"
	"reflect"
)

type iface interface {
	method()
}

type typeOne struct{}

func (typeOne) method() {}

type typeTwo struct{}

func (typeTwo) method() {}

func main() {
	// Static type of variable i is iface. It won’t change. Dynamic type on the other had is … well dynamic. After first assignment, dynamic type of i is typeOne. It isn’t set in stone though so the second assignment changes dynamic type of i to typeTwo. When value of interface type value is nil (which is zero value for interfaces) then dynamic type is not set.
	var i iface = typeOne{}

	// How to get dynamic type of interface type value?
	// Package reflect can be used to achieve that (source code):
	fmt.Println("PkgPath:", reflect.TypeOf(i).PkgPath())
	fmt.Println("Name:   ", reflect.TypeOf(i).Name())
	fmt.Println("String: ", reflect.TypeOf(i).String())

	i = typeTwo{}
	_ = i

	fmt.Println("PkgPath:", reflect.TypeOf(i).PkgPath())
	fmt.Println("Name:   ", reflect.TypeOf(i).Name())
	fmt.Println("String: ", reflect.TypeOf(i).String())
}
