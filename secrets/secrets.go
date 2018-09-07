package secrets

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type secrets struct {
	Token string
}

// ReadTokenFromSecrets will read the JSON file located at 'path' and return
// the value of the 'token' key within. Super secure local token storage!
func ReadTokenFromSecrets(path string) (token string) {
	fileContents, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var tokenMap map[string]string

	if umErr := json.Unmarshal(fileContents, &tokenMap); umErr != nil {
		log.Fatal(umErr)
	}

	if t, ok := tokenMap["token"]; ok {
		token = t
	}

	return
}
