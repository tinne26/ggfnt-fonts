package main

import "os"
import "fmt"
import "image"

import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt/builder"

// NOTES: maybe letters with accents should be normal, take +3 extra ascent?

var oSwitchType uint8

func main() {
	// create font builder
	fmt.Print("creating new font builder\n")
	fontBuilder := builder.New()

	// add metadata
	fmt.Print("...adding metadata\n")
	err := fontBuilder.SetName("omen")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("omen")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("Uppercase font vaguely inspired by old norse glyphs.")
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 1)
	creationDate := ggfnt.Date{ Day: 7, Month: 7, Year: 2024 }
	err = fontBuilder.SetFirstVerDate(creationDate)
	if err != nil { panic(err) }
	err = fontBuilder.SetMajorVerDate(creationDate)
	if err != nil { panic(err) }

	// set metrics
	fmt.Print("...setting metrics\n")
	fontBuilder.SetAscent(6)
	fontBuilder.SetExtraAscent(2)
	fontBuilder.SetUppercaseAscent(6)
	fontBuilder.SetMidlineAscent(0)
	fontBuilder.SetDescent(1)
	fontBuilder.SetExtraDescent(1)
	fontBuilder.SetHorzInterspacing(1)
	fontBuilder.SetLineGap(2)
	err = fontBuilder.GetMetricsStatus()
	if err != nil { panic(err) }

	// create setting for the "O"
	settingKey, err := fontBuilder.AddSetting("o-style", "ornate", "neutral")
	if err != nil { panic(err) }
	oSwitchType, err = fontBuilder.AddSwitch(settingKey)
	if err != nil { panic(err) }

	// add notdef as the first glyph
	notdefUID, err := fontBuilder.AddGlyph(notdef)
	if err != nil { panic(err) }
	err = fontBuilder.SetGlyphName(notdefUID, "notdef")
	if err != nil { panic(err) }
	err = fontBuilder.Map('\uE000', notdefUID)
	if err != nil { panic(err) }

	// add all other glyphs
	runeToUID := make(map[rune]uint64, 128)
	addRunes(fontBuilder, runeToUID, ' ') // spacing
	addRuneRange(fontBuilder, runeToUID, ' ', '`') // until 'a'
	addRuneRange(fontBuilder, runeToUID, '{', '~') // until ASCII END
	addRunes(fontBuilder, runeToUID, '´', '¨', '·', '¡', '¿', '¦') // additional punctuation
	addRunes(fontBuilder, runeToUID,
		'À', 'Á', 'Â', 'Ä',
		'È', 'É', 'Ê', 'Ë',
		'Ì', 'Í', 'Î', 'Ï',
		'Ò', 'Ó', 'Ô', 'Ö',
		'Ù', 'Ú', 'Û', 'Ü',
	) // accents and diacritics
	addRunes(fontBuilder, runeToUID, 'Ñ', 'Ç') // ++spanish letters
	addRunes(fontBuilder, runeToUID, ' ', ' ') // thin space and hair space for padding
	addRunes(fontBuilder, runeToUID, '◀', '▶') // special symbols

	// a few kerning pairs
	fmt.Printf("...configuring kerning pairs\n")
	for _, codePoint := range ".,;:!?▶◀" { // slightly reduce space after punctuation
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID[' '], -1)
	}
	fontBuilder.SetKerningPair(runeToUID['.'], runeToUID['?'], -1) // for better ..?
	fontBuilder.SetKerningPair(runeToUID[' '], runeToUID['▶'], -1)
	for _, codePoint := range "YL" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['-'], -1)
	}
	fontBuilder.SetKerningPair(runeToUID['L'], runeToUID['·'], -1)
	for _, codePoint := range "PY" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['.'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID[','], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['_'], -1)
	}
	fontBuilder.SetKerningPair(runeToUID['L'], runeToUID['Y'], -1)
	fontBuilder.SetKerningPair(runeToUID['T'], runeToUID['T'], -1)

	// show size
	font, err := fontBuilder.Build()
	if err != nil { panic(err) }
	err = font.Validate(ggfnt.FmtDefault)
	if err != nil { panic(err) }
	fmt.Printf("...raw size of %d bytes\n", font.RawSize())

	// export
	const FileName = "omen-6d0-v0p1.ggfnt"
	file, err := os.Create(FileName)
	if err != nil { panic(err) }
	fmt.Printf("...exporting %s\n", FileName)
	err = fontBuilder.Export(file)
	if err != nil {
		_ = file.Close()
		_ = os.Remove(FileName)
		panic(err)
	}
	
	// close file
	fmt.Print("...closing exported file\n")
	err = file.Close()
	if err != nil { panic(err) }
}

func addRuneRange(fontBuilder *builder.Font, codePointsMap map[rune]uint64, start, end rune) {
	for codePoint := start; codePoint <= end; codePoint++ {
		bitmap, found := pkgBitmaps[codePoint]
		if !found { panic("missing bitmap for '" + string(codePoint) + "'") }
		uid, err := fontBuilder.AddGlyph(bitmap)
		if err != nil { panic(err) }
		altBitmap, hasAltBitmap := altBitmaps[codePoint]
		if hasAltBitmap {
			altUID, err := fontBuilder.AddGlyph(altBitmap)
			if err != nil { panic(err) }
			err = fontBuilder.MapWithSwitchSingles(codePoint, oSwitchType, uid, altUID)
		} else {
			err = fontBuilder.Map(codePoint, uid)
		}
		if err != nil { panic(err) }
		codePointsMap[codePoint] = uid
	}
}

func addRunes(fontBuilder *builder.Font, codePointsMap map[rune]uint64, runes ...rune) {
	for _, codePoint := range runes {
		bitmap, found := pkgBitmaps[codePoint]
		if !found { panic("missing bitmap for '" + string(codePoint) + "'") }
		uid, err := fontBuilder.AddGlyph(bitmap) // *
		if err != nil { panic(err) }
		altBitmap, hasAltBitmap := altBitmaps[codePoint]
		if hasAltBitmap {
			altUID, err := fontBuilder.AddGlyph(altBitmap)
			if err != nil { panic(err) }
			err = fontBuilder.MapWithSwitchSingles(codePoint, oSwitchType, uid, altUID)
		} else {
			err = fontBuilder.Map(codePoint, uid)
		}
		if err != nil { panic(err) }
		codePointsMap[codePoint] = uid
	}
}

// helper for mask creation
func rawAlphaMaskToWhiteMask(width int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0, -8, width, -8 + height))
	for i := 0; i < len(mask); i++ {
		img.Pix[i] = 255*mask[i]
	}
	return img
}

var notdef = rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1, // extra ascent
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		1, 1, 1, 1,
})

var pkgBitmaps = map[rune]*image.Alpha{
	' ': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	' ': rawAlphaMaskToWhiteMask(2, []byte{ // thin space
		0, 0, // extra ascent
		0, 0, // extra ascent
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
	}),
	' ': rawAlphaMaskToWhiteMask(1, []byte{ // hair space
		0, // extra ascent
		0, // extra ascent
		0,
		0,
		0,
		0,
		0,
		0, // baseline
		0,
	}),
	'!': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		1,
		1,
		1,
		1,
		0,
		1, // baseline
		0,
	}),
	'"': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'#': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'$': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 1, 0, 0, // extra ascent
		1, 1, 1, 1, 1,
		1, 0, 1, 0, 0,
		1, 0, 1, 1, 1,
		1, 1, 1, 0, 1,
		0, 0, 1, 0, 1,
		1, 1, 1, 1, 1, // baseline
		0, 0, 1, 0, 0,
	}),
	'%': rawAlphaMaskToWhiteMask(9, []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, 0, 0, 0, 0, // extra ascent
		0, 1, 0, 0, 0, 1, 0, 0, 0,
		1, 0, 1, 0, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 0, 1, 0, 1,
		0, 0, 0, 1, 0, 0, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	'&': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 0, 0,
		1, 0, 1, 0,
		1, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 1, 0,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'\'': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,	
	}),
	'(': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		0, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		0, 1, // baseline
		0, 0,
	}),
	')': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		1, 0,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		1, 0, // baseline
		0, 0,
	}),
	'*': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,	
	}),
	'+': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		0, 0, 0,
		0, 0, 0,
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	',': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		1,
	}),
	'-': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
	}),
	'.': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		0,
	}),
	'/': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		0, 0, 1,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
	}),
	'0': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'1': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 0, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'2': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'3': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 1, 0,
		0, 0, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'4': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 0, 1,
		0, 1, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'5': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 0,
		0, 1, 1, 0,
		0, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'6': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 0, 0,
		1, 0, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'7': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'8': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'9': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	':': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		0,
		0,
		1,
		0,
		0,
		1, // baseline
		0,		
	}),
	';': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		0,
		0,
		1,
		0,
		0,
		1, // baseline
		1,		
	}),
	'<': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		0, 0,
		0, 0,
		0, 1,
		1, 0,
		0, 1,
		0, 0, // baseline
		0, 0,
	}),
	'=': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'>': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		0, 0,
		0, 0,
		1, 0,
		0, 1,
		1, 0,
		0, 0, // baseline
		0, 0,
	}),
	'?': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 0, 0, // baseline
		0, 0, 0, 0,
	}),
	'@': rawAlphaMaskToWhiteMask(6, []byte{
		0, 0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, 0, // extra ascent
		1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0, 0,
	}),
	'A': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'B': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 1, 0,
		1, 1, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,		
	}),
	'C': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ç': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 1, 0, 0,
		1, 0, 0, 0,
	}),
	'D': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 0, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 0,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,		
	}),
	'E': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 1, 1,
		1, 0, 0,
		1, 0, 1,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,		
	}),
	'F': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 1, 1,
		1, 0, 0,
		1, 0, 1,
		1, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,		
	}),
	'G': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 0, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,		
	}),
	'H': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'I': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		1,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
		// 0, 0, 0, // extra ascent
		// 0, 0, 0, // extra ascent
		// 1, 1, 1,
		// 0, 1, 0,
		// 0, 1, 0,
		// 0, 1, 0,
		// 0, 1, 0,
		// 1, 1, 1, // baseline
		// 0, 0, 0,
	}),
	'J': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 1,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,		
	}),
	'K': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 1, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'L': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,		
	}),
	'M': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,		
	}),
	'N': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 0, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'Ñ': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 1, // extra ascent
		1, 0, 1, 0, // extra ascent
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 0, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'O': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline1
		0, 0, 0, 0, 0,
	}),
	'P': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 0,
		1, 1, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0, // baseline
		0, 0, 0, 0,
		// 0, 0, 0, // extra ascent
		// 0, 0, 0, // extra ascent
		// 1, 1, 1,
		// 1, 0, 1,
		// 1, 1, 1,
		// 1, 0, 0,
		// 1, 0, 0,
		// 1, 0, 0, // baseline
		// 0, 0, 0,
	}),
	'Q': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,		
	}),
	'R': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 1, 0,
		1, 0, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'S': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 0, 0,
		0, 0, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		// 0, 0, 0, 0, // extra ascent
		// 0, 0, 0, 0, // extra ascent
		// 1, 1, 1, 1,
		// 1, 0, 0, 0,
		// 1, 0, 1, 1,
		// 1, 1, 0, 1,
		// 0, 0, 0, 1,
		// 1, 1, 1, 1, // baseline
		// 0, 0, 0, 0,
	}),
	'T': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		// 0, 0, 0, 0, 0, // extra ascent
		// 0, 0, 0, 0, 0, // extra ascent
		// 1, 1, 1, 1, 1,
		// 0, 1, 0, 1, 0,
		// 0, 0, 1, 0, 0,
		// 0, 0, 1, 0, 0,
		// 0, 0, 1, 0, 0,
		// 0, 0, 1, 0, 0, // baseline
		// 0, 0, 0, 0, 0,
	}),
	'U': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		// 0, 0, 0, // extra ascent
		// 0, 0, 0, // extra ascent
		// 1, 0, 1,
		// 1, 0, 1,
		// 1, 0, 1,
		// 1, 0, 1,
		// 1, 0, 1,
		// 1, 1, 1, // baseline
		// 0, 0, 0,
		// 0, 0, 0, 0, 0, // extra ascent
		// 0, 0, 0, 0, 0, // extra ascent
		// 0, 1, 0, 1, 0,
		// 1, 1, 0, 1, 1,
		// 1, 1, 0, 1, 1,
		// 1, 1, 0, 1, 1,
		// 1, 1, 0, 1, 1,
		// 0, 1, 1, 1, 0, // baseline
		// 0, 0, 0, 0, 0,
	}),
	'V': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,		
	}),
	'W': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0, 0, 0, // extra ascent
		// 0, 0, 0, 0, 0, 0, 0, // extra ascent
		// 1, 0, 1, 0, 1, 0, 1,
		// 1, 0, 1, 0, 1, 0, 1,
		// 0, 1, 0, 1, 0, 1, 0,
		// 0, 1, 0, 1, 0, 1, 0,
		// 0, 0, 1, 0, 1, 0, 0,
		// 0, 0, 1, 0, 1, 0, 0, // baseline
		// 0, 0, 0, 0, 0, 0, 0,
	}),
	'X': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'Y': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'Z': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		1, 1, 1, 1,
		1, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'[': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		1, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 1, // baseline
		0, 0,
	}),
	'\\': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 1, // baseline
		0, 0, 0,
	}),
	']': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		1, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		1, 1, // baseline
		0, 0,
	}),
	'^': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'_': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'`': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		1, 0,
		0, 1,
		0, 0,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
	}),
	'{': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		0, 1, 1,
		0, 1, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 1, // baseline
		0, 0, 0,
	}),
	'|': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		1,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
	}),
	'}': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 1,
		0, 1, 0,
		1, 1, 0, // baseline
		0, 0, 0,
	}),
	'~': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 0, 0,
		1, 0, 1, 0, 1,
		0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),

	// extra, common punctuation
	'´': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0, // extra ascent
		0, 0, // extra ascent
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
	}),
	'¨': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		0, 0, 0, // extra ascent
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,		
	}),
	'·': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		0,
		0,
		0,
		1,
		0,
		0, // baseline
		0,		
	}),
	'¡': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		1,
		0,
		1,
		1,
		1,
		1, // baseline
		0,		
	}),
	'¿': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'¦': rawAlphaMaskToWhiteMask(1, []byte{
		0, // extra ascent
		0, // extra ascent
		1,
		1,
		0,
		0,
		1,
		1, // baseline
		0,		
	}),

	// accented letters
	'À': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 0, // extra ascent
		0, 0, 1, 0, // extra ascent
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'Á': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0, // extra ascent
		0, 1, 0, 0, // extra ascent
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'Â': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 1, 0, // extra ascent
		1, 0, 0, 1, // extra ascent
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'Ä': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		1, 0, 0, 1, // extra ascent
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,		
	}),
	'È': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0, // extra ascent
		0, 1, 0, // extra ascent
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'É': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1, // extra ascent
		0, 1, 0, // extra ascent
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'Ê': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0, // extra ascent
		1, 0, 1, // extra ascent
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'Ë': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		1, 0, 1, // extra ascent
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'Ì': rawAlphaMaskToWhiteMask(2, []byte{
		1, 0, // extra ascent
		0, 1, // extra ascent
		0, 0,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 0,
	}),
	'Í': rawAlphaMaskToWhiteMask(2, []byte{
		0, 1, // extra ascent
		1, 0, // extra ascent
		0, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		0, 0,
	}),
	'Î': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0, // extra ascent
		1, 0, 1, // extra ascent
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
	}),
	'Ï': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0, // extra ascent
		1, 0, 1, // extra ascent
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
	}),
	'Ò': rawAlphaMaskToWhiteMask(5, []byte{
		0, 1, 0, 0, 0, // extra ascent
		0, 0, 1, 0, 0, // extra ascent
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'Ó': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 1, 0, // extra ascent
		0, 0, 1, 0, 0, // extra ascent
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'Ô': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 1, 0, 0, // extra ascent
		0, 1, 0, 1, 0, // extra ascent
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'Ö': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 1, 0, 1, 0, // extra ascent
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'Ù': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 0, // extra ascent
		0, 0, 1, 0, // extra ascent
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ú': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0, // extra ascent
		0, 1, 0, 0, // extra ascent
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Û': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 1, 0, // extra ascent
		1, 0, 0, 1, // extra ascent
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ü': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		1, 0, 0, 1, // extra ascent
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),

	// special symbols
	'◀': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 1,
		0, 0, 1, 1, 1,
		1, 1, 1, 1, 1,
		0, 0, 1, 1, 1,
		0, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
	}),
	'▶': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 0, 0,
		1, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),
}

var altBitmaps = map[rune]*image.Alpha{
	'O': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		0, 0, 0, 0, // extra ascent
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,	
	}),
	'Ò': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 0, // extra ascent
		0, 0, 1, 0, // extra ascent
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'Ó': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0, // extra ascent
		0, 1, 0, 0, // extra ascent
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'Ô': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 1, 0, // extra ascent
		1, 0, 0, 1, // extra ascent
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'Ö': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0, // extra ascent
		1, 0, 0, 1, // extra ascent
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
}
