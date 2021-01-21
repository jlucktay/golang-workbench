package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-playground/mold/v4"
)

func main() {
	tform := mold.New()
	tform.Register("set", transformMyData)

	type Test struct {
		StringField string `mold:"set"`
	}

	tt := Test{StringField: "string"}

	err := tform.Struct(context.Background(), &tt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tt: %#v\n", tt)
}

func transformMyData(_ context.Context, fl mold.FieldLevel) error {
	switch fl.Field().Interface().(type) {
	case string:
		fl.Field().SetString("prefix " + fl.Field().String() + " suffix")
		return nil
	default:
		return fmt.Errorf("nope")
	}
}
