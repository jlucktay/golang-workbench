// Provide command-line arguments in the form of 'YYYY-mmm-DD' where:
// - YYYY is a four digit year
// - mmm is a three letter month abbreviation (Jan, Feb, Mar, et al)
// - DD is a two digit date
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	from, to := os.Args[1], os.Args[2]

	const shortForm = "2006-Jan-02"
	fromTime, _ := time.Parse(shortForm, from)
	toTime, _ := time.Parse(shortForm, to)

	dur := toTime.Sub(fromTime)
	fmt.Println("elapsed days:", int(dur/(time.Hour*24)))
}
