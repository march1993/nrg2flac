package main

import (
	"fmt"
	"nrg2flac/pkg/codec"
	"nrg2flac/pkg/nrg"
	"os"
	"path"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <nrg file>\n", os.Args[0])
		return
	}
	fp := os.Args[1]
	if len(fp) < 4 || strings.ToLower(fp[len(fp)-4:]) != ".nrg" {
		panic("invalid nrg file")
	}
	folder := fp[:len(fp)-4]
	if err := os.MkdirAll(folder, 0755); nil != err {
		panic(err.Error())
	}

	f, err := os.Open(fp)
	if nil != err {
		fmt.Println("os.Open error:", err)
		return
	}
	defer f.Close()

	nf, err := nrg.ReadFooter(f)
	fmt.Printf("nf: %+v\nerr: %+v\n", nf, err)

	chunks, err := nrg.ReadChunks(f, nf.OffsetOfFirstChunk)
	fmt.Printf("chunks: %+v\nerr: %+v\n", chunks, err)

	{
		d, l, ok := nrg.GetChunk[nrg.Daox, nrg.DaoxItem](chunks)
		fmt.Printf("ok: %t\nd: %+v\n", ok, d)
		for idx, t := range l {
			sx, ex := t.Index1, t.Index2
			fmt.Printf("[%d] sx: %d, ex: %d, mode: 0x%x, gap: %d\n", idx, sx, ex, t.Mode, t.Index0)

			f.Seek(int64(sx), 0)

			if err := codec.ToWav(f, int64(ex-sx), path.Join(folder, fmt.Sprintf("track-%d.wav", idx+1))); nil != err {
				panic(err.Error())
			}
		}
	}

}
