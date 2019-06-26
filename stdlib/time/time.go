package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println(time.Now())
	fmt.Println(time.Now())
	log.Println(time.Now().Format(time.RFC3339))
	fmt.Println(time.Now().Format(time.RFC3339))
}

/*
// Parse a time value from a string in the standard Unix format.
t, err := time.Parse(time.UnixDate, "Sat Mar  7 11:06:39 PST 2015")
if err != nil { // Always check errors even if they should not happen.
    panic(err)
}

// time.Time's Stringer method is useful without any format.
fmt.Println("default format:", t)

// Predefined constants in the package implement common layouts.
fmt.Println("Unix format:", t.Format(time.UnixDate))

// The time zone attached to the time value affects its output.
fmt.Println("Same, in UTC:", t.UTC().Format(time.UnixDate))
*/
