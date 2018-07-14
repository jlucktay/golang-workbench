package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

func main() {
	type Result struct {
		Name string `xml:"FullName"`
	}

	v := Result{}

	data := `
		<Person>
			<FullName>Grace R. Emlin</FullName>
			<Company>Example Inc.</Company>
			<Email where="home">
				<Addr>gre@example.com</Addr>
			</Email>
			<Email where='work'>
				<Addr>gre@work.com</Addr>
			</Email>
			<Group>
				<Value>Friends</Value>
				<Value>Squash</Value>
			</Group>
			<City>Hanga Roa</City>
			<State>Easter Island</State>
		</Person>`

	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %q\n", v.Name)
}
