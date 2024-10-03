package main

import "os"
import "fmt"
import "image"

import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt/builder"

var SwitchZeroKey uint8

func main() {
	// create font builder
	fmt.Print("creating new font builder\n")
	fontBuilder := builder.New()

	// add metadata
	fmt.Print("...adding metadata\n")
	err := fontBuilder.SetName("strut")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("strut")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("A swaggy font with full ASCII coverage that's rich and playful without being excessive.")
	if err != nil { panic(err) }
	err = fontBuilder.SetFirstVerDate(ggfnt.Date{ Month: 7, Day: 16, Year: 2024 })
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 1)

	// set metrics
	fmt.Print("...setting metrics\n")
	fontBuilder.SetAscent(6)
	fontBuilder.SetExtraAscent(3)
	fontBuilder.SetUppercaseAscent(6)
	fontBuilder.SetMidlineAscent(4)
	fontBuilder.SetDescent(2)
	fontBuilder.SetExtraDescent(0)
	fontBuilder.SetHorzInterspacing(1)
	fontBuilder.SetLineGap(2)
	err = fontBuilder.GetMetricsStatus()
	if err != nil { panic(err) }

	// zero disambiguation mark
	settingZeroDisKey, err := fontBuilder.AddSetting("zero-disambiguation-mark", "off", "on")
	if err != nil { panic(err) }
	SwitchZeroKey, err = fontBuilder.AddSwitch(settingZeroDisKey)
	if err != nil { panic(err) }

	// add notdef as the glyph zero
	notdefUID, err := fontBuilder.AddGlyph(notdef)
	if err != nil { panic(err) }
	err = fontBuilder.SetGlyphName(notdefUID, "notdef")
	if err != nil { panic(err) }
	err = fontBuilder.Map('\uE000', notdefUID)
	if err != nil { panic(err) }

	// add glyphs and map them
	runeToUID := make(map[rune]uint64, 128)
	addRuneRange(fontBuilder, runeToUID, ' ', '~') // ascii
	addRunes(fontBuilder, runeToUID,
		'À', 'Á', 'Ä', 'Â',
		'È', 'É', 'Ë', 'Ê',
		'Ì', 'Í', 'Ï', 'Î',
		'Ò', 'Ó', 'Ö', 'Ô',
		'Ù', 'Ú', 'Ü', 'Û',
	)
	addRunes(fontBuilder, runeToUID,
		'à', 'á', 'ä', 'â',
		'è', 'é', 'ë', 'ê',
		'ì', 'í', 'ï', 'î',
		'ò', 'ó', 'ö', 'ô',
		'ù', 'ú', 'ü', 'û',
	)
	addRunes(fontBuilder, runeToUID, '¡', '¿', '´', '¨', '·', '¦', 'º', '…', '—', '¬', '‘', '’', '“', '”') // additional symbols
	addRunes(fontBuilder, runeToUID, 'Ñ', 'ñ', 'Ç', 'ç') // spanish letters
	addRunes(fontBuilder, runeToUID, '€', '¢', '¥', '¤') // some annoying currency symbols

	// kerning
	fmt.Printf("...configuring kerning pairs\n")
	for _, codePoint := range ".,;:!?" { // slightly reduce space after punctuation
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID[' '], -1)
	}
	fontBuilder.SetKerningPair(runeToUID['.'], runeToUID['?'], -2) // improve "..?"

	for _, codePoint := range "acegmnopqrsuvwxyz" { // slightly reduce space for most F + lowercase pairings
		fontBuilder.SetKerningPair(runeToUID['F'], runeToUID[codePoint], -1)
	}

	// adjust j (would greatly benefit from kerning ranges or default to -1,
	// which could indeed be set as the first value with a special control...)
	for _, codePoint := range ".-~—\"'‘’“$AÀÁÄÂBCÇDEÈÉËÊFGHIÌÍÏÎJKLMNÑOÒÓÖÔPQRSTUÙÚÜÛVWXYZaàáäâbcçdeèéëêfhiìíïîklmnñoòóöôprstuùúüûvwxz0123456789" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['j'], -1)
	}
	// adjust ï and î left sides
	for _, codePoint := range "_aàbcçdeèhijkmnoòprstuùvwxz" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], -1)
	}

	// ensure that all bitmaps are being used
	for codePoint, _ := range pkgBitmaps {
		_, found := runeToUID[codePoint]
		if !found { panic("missing code point '" + string(codePoint) + "'") }
	}

	// show size
	font, err := fontBuilder.Build()
	if err != nil { panic(err) }
	err = font.Validate(ggfnt.FmtDefault)
	if err != nil { panic(err) }
	fmt.Printf("...raw size of %d bytes\n", font.RawSize())

	// export
	const FileName = "strut-6d2-v0p1.ggfnt"
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
		bitmap, found := pkgBitmaps[rune(codePoint)]
		if !found { panic("missing bitmap for '" + string(codePoint) + "'") }
		uid, err := fontBuilder.AddGlyph(bitmap)
		if err != nil { panic(err) }
		if codePoint == '0' {
			disZeroUID, err := fontBuilder.AddGlyph(disZero)
			if err != nil { panic(err) }
			err = fontBuilder.MapWithSwitchSingles(codePoint, SwitchZeroKey, uid, disZeroUID)
		} else {
			err = fontBuilder.Map(codePoint, uid)
		}
		if err != nil { panic(err) }
		codePointsMap[codePoint] = uid
	}
}

func addRunes(fontBuilder *builder.Font, codePointsMap map[rune]uint64, runes ...rune) {
	for _, codePoint := range runes {
		bitmap, found := pkgBitmaps[rune(codePoint)]
		if !found { panic("missing bitmap for '" + string(codePoint) + "'") }
		uid, err := fontBuilder.AddGlyph(bitmap) // *
		if err != nil { panic(err) }
		err = fontBuilder.Map(codePoint, uid)
		if err != nil { panic(err) }
		codePointsMap[codePoint] = uid
	}
}

// helper for mask creation
func rawMask(width, descent int, mask []byte) *image.Alpha {
	height := len(mask)/width
	if len(mask)%width != 0 { panic("mask size doesn't match given width") }
	img := image.NewAlpha(image.Rect(0, -height + descent, width, descent))
	copy(img.Pix, mask)
	return img
}

var notdef = rawMask(4, 1, []byte{
	1, 1, 1, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1, // baseline
	1, 1, 1, 1,
})

var disZero = rawMask(4, 0, []byte{
	1, 1, 1, 1,
	1, 0, 1, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 1, 0, 1,
	1, 1, 1, 1,
})

var pkgBitmaps = map[rune]*image.Alpha{
	' ': rawMask(4, 0, []byte{
		0, 0, 0, 0,
	}),
	'!': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		0, 0,
		1, 1,
	}),
	'"': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'#': rawMask(7, 0, []byte{
		0, 1, 1, 0, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 0, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 0, 1, 1, 0,
	}),
	'$': rawMask(7, 1, []byte{
		0, 0, 0, 1, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 1, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 1, 0, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 1, 0, 0, 0,
	}),
	'%': rawMask(10, 0, []byte{
		1, 1, 1, 0, 0, 1, 1, 0, 0, 0,
		1, 0, 1, 0, 0, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 0, 1, 1, 1,
		0, 0, 0, 1, 1, 0, 0, 1, 0, 1,
		0, 0, 0, 1, 1, 0, 0, 1, 1, 1,
	}),
	'&': rawMask(6, 0, []byte{
		0, 1, 1, 0, 0, 0,
		1, 1, 0, 1, 0, 0,
		0, 1, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 0,
		1, 1, 0, 0, 1, 1,
		0, 1, 1, 1, 1, 0,
	}),
	'\'': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		0, 0,
	}),
	'(': rawMask(3, 1, []byte{
		0, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		0, 1, 1,
	}),
	')': rawMask(3, 1, []byte{
		1, 1, 0,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		1, 1, 0,
	}),
	'*': rawMask(4, 0, []byte{
		1, 0, 0, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'+': rawMask(4, 0, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 0, 0, 0,
	}),
	',': rawMask(2, 2, []byte{
		1, 1,
		1, 1,
		0, 1,
		1, 0,
	}),
	'-': rawMask(4, 0, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'.': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
	}),
	'/': rawMask(4, 0, []byte{
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
	}),
	'0': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'1': rawMask(4, 0, []byte{
		0, 1, 1, 0,
		1, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
	}),
	'2': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'3': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		1, 1, 1, 1,
	}),
	'4': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
	}),
	'5': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		1, 1, 1, 1,
	}),
	'6': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'7': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
	}),
	'8': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'9': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
	}),
	':': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		0, 0,
		1, 1,
		1, 1,
	}),
	';': rawMask(2, 2, []byte{
		1, 1,
		1, 1,
		0, 0,
		1, 1,
		1, 1,
		0, 1,
		1, 0,
	}),
	'<': rawMask(3, 0, []byte{
		0, 1, 1,
		1, 1, 0,
		1, 1, 0,
		0, 1, 1,
		0, 0, 0,
	}),
	'=': rawMask(5, 0, []byte{
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	'>': rawMask(3, 0, []byte{
		1, 1, 0,
		0, 1, 1,
		0, 1, 1,
		1, 1, 0,
		0, 0, 0,
	}),
	'?': rawMask(5, 0, []byte{
		1, 1, 1, 1, 0,
		1, 1, 0, 1, 1,
		0, 0, 0, 1, 1,
		0, 0, 1, 1, 0,
		0, 0, 0, 0, 0,
		0, 0, 1, 1, 0,
	}),
	'@': rawMask(6, 0, []byte{
		1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 0,
	}),
	'A': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'B': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'C': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'Ç': rawMask(4, 1, []byte{
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
		0, 0, 1, 0,
	}),
	'D': rawMask(4, 0, []byte{
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		// 1, 1, 1, 1, 0,
		// 1, 1, 1, 0, 1,
		// 1, 1, 0, 0, 1,
		// 1, 1, 0, 0, 1,
		// 1, 1, 1, 0, 1,
		// 1, 1, 1, 1, 0,
	}),
	'E': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'F': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
	}),
	'G': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'H': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'I': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
	}),
	'J': rawMask(4, 0, []byte{
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
	}),
	'K': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'L': rawMask(3, 0, []byte{
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
	}),
	'M': rawMask(6, 0, []byte{
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 1, 1, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 1,
	}),
	'N': rawMask(4, 0, []byte{
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'Ñ': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		1, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'O': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'P': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
	}),
	'Q': rawMask(5, 0, []byte{
		1, 1, 1, 1, 1,
		1, 1, 0, 0, 1,
		1, 1, 0, 0, 1,
		1, 1, 0, 0, 1,
		1, 1, 0, 1, 1,
		1, 1, 1, 1, 1,
	}),
	'R': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'S': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 0, 1, 1,
		1, 1, 1, 1,
	}),
	'T': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
	}),
	'U': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'V': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		0, 1, 1, 0,
	}),
	'W': rawMask(6, 0, []byte{
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 1, 0,
	}),
	'X': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'Y': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
	}),
	'Z': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 1,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'[': rawMask(3, 1, []byte{
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
	}),
	'\\': rawMask(4, 0, []byte{
		1, 1, 0, 0,
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 1,
		0, 0, 1, 1,
	}),
	']': rawMask(3, 1, []byte{
		1, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		1, 1, 1,
	}),
	'^': rawMask(3, 0, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'_': rawMask(5, 2, []byte{
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
	}),
	'`': rawMask(2, 0, []byte{
		1, 0,
		0, 1,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
	}),
	'a': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
	}),
	'b': rawMask(4, 0, []byte{
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'c': rawMask(3, 0, []byte{
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
	}),
	'ç': rawMask(3, 1, []byte{
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
		0, 1, 0,
	}),
	'd': rawMask(4, 0, []byte{
		0, 0, 1, 1,
		0, 0, 1, 1,
		1, 1, 1, 1,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
	}),
	'e': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 1, 1,
	}),
	'f': rawMask(3, 0, []byte{
		0, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
	}),
	'g': rawMask(4, 2, []byte{
		1, 1, 1, 1,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 1,
	}),
	'h': rawMask(4, 0, []byte{
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'i': rawMask(2, 0, []byte{
		1, 0,
		0, 0,
		1, 0,
		1, 1,
		1, 1,
		1, 1,
	}),
	'j': rawMask(3, 2, []byte{
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		1, 1, 1,
	}),
	'k': rawMask(4, 0, []byte{
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
	}),
	'l': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
	}),
	'm': rawMask(6, 0, []byte{
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1,
	}),
	'n': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'ñ': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'o': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'p': rawMask(4, 2, []byte{
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
	}),
	'q': rawMask(4, 2, []byte{
		1, 1, 1, 1,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
	}),
	'r': rawMask(3, 0, []byte{
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
	}),
	's': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		1, 1, 0, 0,
		0, 0, 1, 1,
		1, 1, 1, 1,
	}),
	't': rawMask(3, 0, []byte{
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
	}),
	'u': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'v': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		0, 1, 1, 0,
	}),
	'w': rawMask(6, 0, []byte{
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 1, 1, 1, 1,
	}),
	'x': rawMask(4, 0, []byte{
		1, 1, 0, 1,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'y': rawMask(4, 2, []byte{
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 1,
	}),
	'z': rawMask(4, 0, []byte{
		1, 1, 1, 1,
		0, 0, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'{': rawMask(4, 1, []byte{
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 1,
	}),
	'|': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
	}),
	'}': rawMask(4, 1, []byte{
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
	}),
	'~': rawMask(6, 0, []byte{
		1, 1, 1, 1, 0, 0,
		1, 1, 1, 1, 0, 1,
		1, 0, 1, 1, 1, 1,
		0, 0, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
	}),

	// additional punctuation symbols
	'¡': rawMask(2, 2, []byte{
		1, 1,
		0, 0,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
	}),
	'¿': rawMask(5, 2, []byte{
		0, 1, 1, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 1, 0, 0,
		1, 1, 0, 0, 0,
		1, 1, 0, 1, 1,
		0, 1, 1, 1, 1,
	}),
	'´': rawMask(2, 0, []byte{
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
	}),
	'¨': rawMask(3, 0, []byte{
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'·': rawMask(1, 0, []byte{
		1,
		0,
		0,
	}),
	'¦': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		0, 0,
		1, 1,
		1, 1,
		1, 1,
	}),
	'º': rawMask(3, 0, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'…': rawMask(8, 0, []byte{
		1, 1, 0, 1, 1, 0, 1, 1,
		1, 1, 0, 1, 1, 0, 1, 1,
	}),
	'—': rawMask(6, 0, []byte{
		1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'¬': rawMask(5, 0, []byte{
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 0,
	}),
	'‘': rawMask(2, 0, []byte{
		0, 1,
		1, 0,
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		0, 0,
	}),
	'’': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0,
	}),
	'“': rawMask(5, 0, []byte{
		0, 1, 0, 0, 1,
		1, 0, 0, 1, 0,
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'”': rawMask(5, 0, []byte{
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		0, 1, 0, 0, 1,
		1, 0, 0, 1, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),

	// uppercase letters with accents
	'À': rawMask(4, 0, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'Á': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'Ä': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'Â': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
	}),
	'È': rawMask(4, 0, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'É': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'Ë': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'Ê': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'Ì': rawMask(2, 0, []byte{
		1, 0,
		0, 1,
		0, 0,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
	}),
	'Í': rawMask(2, 0, []byte{
		0, 1,
		1, 0,
		0, 0,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
	}),
	'Ï': rawMask(3, 0, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
	}),
	'Î': rawMask(3, 0, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
	}),
	'Ò': rawMask(4, 0, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'Ó': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'Ö': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'Ô': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'Ù': rawMask(4, 0, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'Ú': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'Ü': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'Û': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),

	// lowercase letters with accents
	'à': rawMask(4, 0, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
	}),
	'á': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
	}),
	'ä': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
	}),
	'â': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
	}),
	'è': rawMask(4, 0, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 1, 1,
	}),
	'é': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 1, 1,
	}),
	'ë': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 1, 1,
	}),
	'ê': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 1, 1,
	}),
	'ì': rawMask(2, 0, []byte{
		1, 0,
		0, 1,
		0, 0,
		1, 0,
		1, 1,
		1, 1,
		1, 1,
	}),
	'í': rawMask(2, 0, []byte{
		0, 1,
		1, 0,
		0, 0,
		1, 0,
		1, 1,
		1, 1,
		1, 1,
	}),
	'ï': rawMask(3, 0, []byte{
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
	}),
	'î': rawMask(3, 0, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
	}),
	'ò': rawMask(4, 0, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'ó': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'ö': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'ô': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}),
	'ù': rawMask(4, 0, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'ú': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'ü': rawMask(4, 0, []byte{
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),
	'û': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
	}),

	// annoying currency symbols
	'€': rawMask(4, 0, []byte{
		0, 0, 1, 1,
		0, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 1, 0,
		0, 1, 0, 0,
		0, 0, 1, 1,
	}),
	'¢': rawMask(4, 0, []byte{
		0, 0, 1, 0,
		0, 1, 1, 1,
		1, 0, 1, 0,
		1, 0, 1, 0,
		0, 1, 1, 1,
		0, 0, 1, 0,
	}),
	'¥': rawMask(6, 0, []byte{
		1, 1, 0, 0, 0, 1,
		1, 1, 1, 0, 1, 0,
		0, 1, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 1,
		0, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0,
	}),
	'¤': rawMask(5, 0, []byte{
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
		0, 1, 0, 1, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
	}),
}
