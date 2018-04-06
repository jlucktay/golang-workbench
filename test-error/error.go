package main

import "errors"

func returnSomeErr(input int) error {
	if input > 0 {
		return nil
	}

	return errors.New("this is error message")
}
