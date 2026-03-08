package readbyte

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestReadByte(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	defer func() {
		log.SetOutput(os.Stderr)
	}()

	readByte()
	t.Log(buf.String())
}
