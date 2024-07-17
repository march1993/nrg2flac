package nrg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

var (
	ErrInvalidNrgFile = errors.New("invalid nrg file")
)

var (
	NERO = [4]uint8{'N', 'E', 'R', 'O'}
	NER5 = [4]uint8{'N', 'E', 'R', '5'}
)

type NeroFooter struct {
	Version            int      // v1 or v2
	ChunkID            [4]uint8 // NERO
	OffsetOfFirstChunk uint64
}

func ReadFooter(f *os.File) (NeroFooter, error) {
	nf := NeroFooter{}

	f.Seek(-12, 2)
	f.Read(nf.ChunkID[:])
	if nf.ChunkID == NER5 {
		// v2
		nf.Version = 2
		buf := make([]byte, 8)
		if _, err := f.Read(buf); nil != err {
			return nf, err
		}
		nf.OffsetOfFirstChunk = binary.BigEndian.Uint64(buf)
		return nf, nil
	}

	f.Seek(-8, 2)
	f.Read(nf.ChunkID[:])
	if nf.ChunkID == NERO {
		// v1
		nf.Version = 1
		buf := make([]byte, 4)
		if _, err := f.Read(buf); nil != err {
			return nf, err
		}
		nf.OffsetOfFirstChunk = uint64(binary.BigEndian.Uint32(buf))
		return nf, nil
	}

	f.Seek(-12, 2)
	buf := make([]byte, 12)
	f.Read(buf)
	fmt.Printf("last 12 bytes: %+v\n", buf)
	return nf, ErrInvalidNrgFile
}
