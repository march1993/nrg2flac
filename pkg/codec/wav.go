package codec

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

// http://soundfile.sapp.org/doc/WaveFormat/
// https://en.wikipedia.org/wiki/Compact_Disc_Digital_Audio

func ToWav(f io.Reader, size int64, fp string) error {
	wh := WavHeader{
		ChunkID:       riff,
		ChunkSize:     36 + uint32(size),
		Format:        wave,
		Subchunk1ID:   fmt_,
		Subchunk1Size: 16,
		AudioFormat:   1, // PCM
		NumChannels:   2,
		SampleRate:    44100,
		ByteRate:      44100 * 2 * 2,
		BlockAlign:    4, // NumChannels * BitsPerSample/8
		BitsPerSample: 16,
		Subchunk2ID:   data,
		Subchunk2Size: uint32(size),
	}

	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, wh)
	if _, err := io.CopyN(buf, f, size); nil != err {
		return err
	}

	return os.WriteFile(fp, buf.Bytes(), 0644)
}

// combined RIFF, fmt and data header
// little endian
type WavHeader struct {
	ChunkID       [4]uint8
	ChunkSize     uint32
	Format        [4]uint8
	Subchunk1ID   [4]uint8
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]uint8
	Subchunk2Size uint32
	// followed by data
}

var (
	riff = [4]uint8{'R', 'I', 'F', 'F'}
	wave = [4]uint8{'W', 'A', 'V', 'E'}
	fmt_ = [4]uint8{'f', 'm', 't', ' '}
	data = [4]uint8{'d', 'a', 't', 'a'}
)
