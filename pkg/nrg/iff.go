package nrg

// https://en.wikipedia.org/wiki/NRG_%28file_format%29

var (
	CUES = ChunkID{'C', 'U', 'E', 'S'}
	CUEX = ChunkID{'C', 'U', 'E', 'X'}
	DAOI = ChunkID{'D', 'A', 'O', 'I'}
	DAOX = ChunkID{'D', 'A', 'O', 'X'}
	CDTX = ChunkID{'C', 'D', 'T', 'X'} // TODO:
	ETNF = ChunkID{'E', 'T', 'N', 'F'} // TODO:
	ETN2 = ChunkID{'E', 'T', 'N', '2'} // TODO:
	SINF = ChunkID{'S', 'I', 'N', 'F'}
	MTYP = ChunkID{'M', 'T', 'Y', 'P'}
	DINF = ChunkID{'D', 'I', 'N', 'F'} // TODO:
	TOCT = ChunkID{'T', 'O', 'C', 'T'} // TODO:
	RELO = ChunkID{'R', 'E', 'L', 'O'} // TODO:
	END_ = ChunkID{'E', 'N', 'D', '!'}
)

type End_ struct{}
type End_Item struct{}

type Cues struct{}
type CuesItem struct{}
type Cuex struct{}
type CuexItem struct {
	Mode                 uint8  // values found: 0x01 for audio; 0x21 for non copyright-protected audio; 0x41 for data
	TrackNumber          uint8  // BCD coded; 0xAA for the lead-out area
	IndexNumber          uint8  // probably BCD coded
	Padding              uint8  // always zero found
	LbaPositionInSectors uint32 // signed integer value
}

type Daoi struct {
	ChunkSize  uint32 // bytes little endian
	UPC        [14]uint8
	TocType    uint32
	FirstTrack uint8
	LastTrack  uint8
}
type DaoiItem struct {
	ISRC       [12]byte
	SectorSize uint32
	Mode       uint32
	Index0     uint32 // pre gap
	Index1     uint32 // start of tract
	Index2     uint32 // end of track + 1
}

type Daox struct {
	ChunkSize1 uint32 // bytes big endian already encountered; maybe also little endian on some machines
	UPC        [13]uint8
	Padding    uint8  // always NULL found
	TocType    uint16 // values already found: 0x0000 for audio; 0x0001 for data; 0x2001 for Mode 2/form 1 data
	FirstTrack uint8  // in the session
	LastTrack  uint8  // in the session
}

type DaoxItem struct {
	Text       [12]uint8 // or NULLs
	SectorSize uint16    // in the image file
	Mode       uint16    // Mode of the data in the image file (values already found: 0x0700 for audio; 0x1000 for audio with sub-channel; 0x0000 for data; 0x0500 for raw data; 0x0f00 for raw data with sub-channel; 0x0300 for Mode 2 Form 1 data; 0x0600 for raw Mode 2/form 1 data; 0x1100 for raw Mode 2/form 1 data with sub-channel)
	Unknown    uint16    // always 0x0001 found
	Index0     uint64    // pre-gap bytes
	Index1     uint64    // start of track (bytes)
	Index2     uint64    // end of track + 1 (bytes)
}

type Sinf struct {
	NTracks int32 // in session
}
type SinfItem struct{}

type Mtyp struct {
	Unknown uint32
}
type MtypItem struct{}
