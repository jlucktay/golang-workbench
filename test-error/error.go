package testerr

import "errors"

func ReturnSomeErr(input int) error {
	if input > 0 {
		return nil
	}

	return errors.New("this is error message")
}
