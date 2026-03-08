package secrets

import (
	"encoding/json"
	"log"
	"os"
)

// ReadTokenFromSecrets will read the JSON file located at 'path' and return the value of the 'token' key within. Super secure local token storage!
func ReadTokenFromSecrets(path string) string {
	fileContents, readErr := os.ReadFile(path)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var tokenMap map[string]string

	if umErr := json.Unmarshal(fileContents, &tokenMap); umErr != nil {
		log.Fatal(umErr)
	}

	if t, ok := tokenMap["token"]; ok {
		return t
	}

	return ""
}
