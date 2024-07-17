package nrg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type Chunk struct {
	ChunkHead

	Data any
	List any
}

type ChunkID [4]uint8
type ChunkHead struct {
	ChunkID   ChunkID
	ChunkSize uint32 // bytes
}

func ReadChunks(f *os.File, offsetOfFirstChunk uint64) (chunks []Chunk, err error) {
	if _, err = f.Seek(int64(offsetOfFirstChunk), 0); nil != err {
		return
	}

	reading := true
	for reading {
		c := Chunk{}

		if err = binary.Read(f, binary.BigEndian, &c.ChunkHead); nil != err {
			return
		}
		switch c.ChunkID {
		case END_:
			err = readChunk[End_, End_Item](f, &c)
			reading = false
			if c.ChunkSize != 0 {
				err = ErrInvalidEnd_Chunk
			}
		case CUES:
			err = readChunk[Cues, CuesItem](f, &c)
		case CUEX:
			err = readChunk[Cuex, CuexItem](f, &c)
		case DAOI:
			err = readChunk[Daoi, DaoiItem](f, &c)
		case DAOX:
			err = readChunk[Daox, DaoxItem](f, &c)
		case SINF:
			err = readChunk[Sinf, SinfItem](f, &c)
		case MTYP:
			err = readChunk[Mtyp, MtypItem](f, &c)
		default:
			return chunks, fmt.Errorf("chunk id: %+v", c.ChunkID)
		}

		if nil != err {
			fmt.Printf("invalid chunk: %+v\n", c.ChunkHead)
			return
		}

		chunks = append(chunks, c)
	}

	return
}

func readChunk[D, I any](f *os.File, c *Chunk) error {
	d := new(D)
	if err := binary.Read(f, binary.BigEndian, d); nil != err {
		return err
	}
	c.Data = d

	cs := int(c.ChunkHead.ChunkSize)
	ds := binary.Size(new(D))
	is := binary.Size(new(I))

	if is == 0 {
		c.List = []I{}
		if cs != ds {
			return ErrInvalidChunkData
		}
	} else {
		fmt.Printf("cs: %d, ds: %d, is: %d\n", cs, ds, is)
		if int(cs-ds)%is != 0 {
			fmt.Printf("cs: %d, ds: %d, is: %d\n", cs, ds, is)
			return ErrInvalidChunkItem
		}
		n := (cs - ds) / is

		l := make([]I, n)

		if err := binary.Read(f, binary.BigEndian, l); nil != err {
			return err
		}

		c.List = l
	}

	return nil
}

func (id ChunkID) String() string {
	return string(id[:])
}

// return first found chunk with certain type
func GetChunk[D, I any](chunks []Chunk) (d D, l []I, ok bool) {
	for _, c := range chunks {
		if d, ok := c.Data.(*D); ok {
			if l, ok := c.List.([]I); ok {
				return *d, l, true
			} else {
				panic("chunk data and list type mismatch")
			}
		}
	}

	return
}

var (
	ErrInvalidChunkItem = errors.New("invalid chunk item")
	ErrInvalidChunkData = errors.New("invalid chunk data")
	ErrInvalidEnd_Chunk = errors.New("invalid END! chunk")
	ErrChunkNotFound    = errors.New("chunk not found")
)
