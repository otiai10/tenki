package tenki

import (
	"testing"

	. "github.com/otiai10/mint"
)

func TestLoad(t *testing.T) {
	err := Load([]byte("{\"key\": \"value\"}"))
	Expect(t, err).ToBe(nil)
}
