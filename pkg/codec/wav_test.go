package codec

import (
	"encoding/binary"
	"testing"
)

func TestToWav(t *testing.T) {
	if binary.Size(WavHeader{}) != 44 {
		t.Fail()
	}
}
