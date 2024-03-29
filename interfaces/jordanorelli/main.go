// Playing with JSON marshaling combined with parsing timestamp fields.
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// start with a string representation of our JSON data
var input = `
{
	"created_at": "Thu May 31 00:00:01 +0000 2012"
}
`

func main() {
	// our target will be of type map[string]interface{}, which is a pretty generic type
	// that will give us a hashtable whose keys are strings, and whose values are of
	// type interface{}
	// var val map[string]interface{}
	var val map[string]Timestamp

	if err := json.Unmarshal([]byte(input), &val); err != nil {
		panic(err)
	}

	fmt.Println(val)

	for k, v := range val {
		fmt.Println(k, ":", reflect.TypeOf(v))
		fmt.Printf("%s : %T\n", k, v)
	}

	fmt.Println(time.Time(val["created_at"]))
}

// Timestamp is a proxy type for wrangling Twitter things
type Timestamp time.Time

// UnmarshalJSON on a Timestamp does some parsing for us
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	v, err := time.Parse(time.RubyDate, string(b[1:len(b)-1]))
	if err != nil {
		return err
	}

	*t = Timestamp(v)

	return nil
}
