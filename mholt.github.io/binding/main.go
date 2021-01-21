package main

import (
	"fmt"
	"net/http"

	"github.com/mholt/binding"
)

// First define a type to hold the data
// (If the data comes from JSON, see: http://mholt.github.io/json-to-go)
type ContactForm struct {
	User struct {
		ID int
	}
	Email   string
	Message string
}

// Then provide a field mapping (pointer receiver is vital)
func (cf *ContactForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&cf.User.ID: "user_id",
		&cf.Email:   "email",
		&cf.Message: binding.Field{
			Form:     "message",
			Required: true,
		},
	}
}

// Now your handlers can stay clean and simple.
func handler(resp http.ResponseWriter, req *http.Request) {
	contactForm := new(ContactForm)
	errs := binding.Bind(req, contactForm)
	if errs.Handle(resp) {
		return
	}
	fmt.Fprintf(resp, "From:    %d\n", contactForm.User.ID)
	fmt.Fprintf(resp, "Message: %s\n", contactForm.Message)
}

func main() {
	http.HandleFunc("/contact", handler)
	http.ListenAndServe(":3000", nil)
}
