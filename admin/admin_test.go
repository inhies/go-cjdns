package admin

import (
	"bytes"
	"testing"
)

func TestPathMarshalingUnmarshaling(t *testing.T) {
	path := new(Path)
	if err := path.UnmarshalText([]byte("0000.0114.a785.58e3")); err != nil {
		t.Error("Failed to unmarshal Path,", err)
		return
	}
	if *path == 0 {
		t.Error("unmarshaled path was empty")
		return
	}

	test, err := path.MarshalText()
	if err != nil {
		t.Error("Failed to marshal Path,", err)
		return
	}
	if !bytes.Equal([]byte("0000.0114.a785.58e3"), test) {
		t.Errorf("Path marshal and unmarshal mismatch, wanted \"0000.0114.a785.58e3\", got %q", test)
	}
}
