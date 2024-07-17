package codec

import (
	"bytes"

	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/frame"
	"github.com/mewkiz/flac/meta"
)

func SaveFlac() []byte {
	buf := bytes.NewBuffer(nil)

	si := &meta.StreamInfo{
		BlockSizeMin:  16,
		BlockSizeMax:  16,
		BitsPerSample: 16,
		SampleRate:    44100,
		NChannels:     2,
	}
	if enc, err := flac.NewEncoder(buf, si); nil != err {
		panic(err.Error())
	} else {
		f := &frame.Frame{
			Header: frame.Header{
				SampleRate: 44100,
				Channels:   frame.ChannelsLR,
			},
			Subframes: []*frame.Subframe{{}},
		}
		enc.WriteFrame(f)
	}

	return buf.Bytes()
}
