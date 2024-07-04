package jammy

import "io"
import _ "embed"
import "github.com/tinne26/ggfnt"

//go:embed jammy-5d2-v0p3.ggfnt
var bytes []byte

var cachedFont *ggfnt.Font

const NotdefRune = '\uE000'
const Notdef = ggfnt.GlyphIndex(0)
const LowHyphenRune = '\uE001'
const LowHyphen = ggfnt.GlyphIndex(180)

const ZeroDisambiguationMarkSettingKey = ggfnt.SettingKey(0)
const ZeroDisambiguationMarkSettingName = "zero-disambiguation-mark"
const NumericStyleSettingKey = ggfnt.SettingKey(1)
const NumericStyleSettingName = "numeric-style"

func Release() { cachedFont = nil }
func Font() *ggfnt.Font {
	if cachedFont == nil {
		font, err := ggfnt.Parse(&byteSliceReader{ data: bytes })
		if err != nil { panic(err) } // (go test .)
		cachedFont = font
	}
	return cachedFont
}

type byteSliceReader struct { data []byte ; index int }
func (self *byteSliceReader) Read(buffer []byte) (int, error) {
	// determine max read size and stop if nothing to read
	maxRead := len(self.data) - self.index
	if maxRead <= 0 { return 0, io.EOF }
	if len(buffer) == 0 { return 0, nil }
	
	// determine final read size and copy the data
	readSize := min(maxRead, len(buffer))
	copy(buffer, self.data[self.index : self.index + maxRead])
	self.index += readSize
	if len(buffer) < maxRead { return readSize, nil }
	return readSize, io.EOF
}
