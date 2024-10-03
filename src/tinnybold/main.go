package main

import "os"
import "fmt"
import "image"

import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt/builder"

func main() {
	// create font builder
	fmt.Print("creating new font builder\n")
	fontBuilder := builder.New()

	// add metadata
	fmt.Print("...adding metadata\n")
	err := fontBuilder.SetName("tinny-bold")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("tinny")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("A bold version of 'tinny', the low resolution, rounded and friendly-looking font.")
	if err != nil { panic(err) }
	err = fontBuilder.SetFirstVerDate(ggfnt.Date{ Month: 10, Day: 2, Year: 2024 })
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 1)

	// set metrics
	fmt.Print("...setting metrics\n")
	fontBuilder.SetAscent(7)
	fontBuilder.SetExtraAscent(0)
	fontBuilder.SetUppercaseAscent(6)
	fontBuilder.SetMidlineAscent(4)
	fontBuilder.SetDescent(3)
	fontBuilder.SetHorzInterspacing(1)
	fontBuilder.SetLineGap(1)
	err = fontBuilder.GetMetricsStatus()
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
	addRuneRange(fontBuilder, runeToUID, ' ', '~') // ascii table
	addRunes(fontBuilder, runeToUID,
		'À', 'Á', 'Â', 'Ä', 'à', 'á', 'â', 'ä',
		'È', 'É', 'Ê', 'Ë', 'è', 'é', 'ê', 'ë',
		'Ì', 'Í', 'Î', 'Ï', 'ì', 'í', 'î', 'ï',
		'Ò', 'Ó', 'Ô', 'Ö', 'ò', 'ó', 'ô', 'ö',
		'Ù', 'Ú', 'Û', 'Ü', 'ù', 'ú', 'û', 'ü',
	) // accents and diacritics
	addRunes(fontBuilder, runeToUID,
		'¡', '¿', '´', '¨', '·', '¦', '‘', '’', '“', '”', '…',
	) // additional punctuation
	addRunes(fontBuilder, runeToUID, '−', '×', '÷', '±', 'º', '′', '″', '¬', 'π') // ++maths
	addRunes(fontBuilder, runeToUID, '€', '£', '¢', '¥', '¤') // currency symbols
	addRunes(fontBuilder, runeToUID,
		'Ñ', 'ñ', 'Ç', 'ç', 'Ø', 'ø', 'Þ', 'þ', 'Ð', 'ð', 'Đ', 'đ', 
		'Æ', 'æ', 'Œ', 'œ', 'ƒ',
	) // ++latin letters
	addRunes(fontBuilder, runeToUID, 'Ǿ', 'ǿ', /*'Ø̈', 'ø̈' */) // ++latin accents and diacritics
	addRunes(fontBuilder, runeToUID, '–', '‑', '—') // ++dashes
	addRunes(fontBuilder, runeToUID, '♩', '♪', '♫', '♬') // notes
	addRunes(fontBuilder, runeToUID, '�') // special

	// adjust placement for a few glyphs
	fontBuilder.SetGlyphPlacement(runeToUID['j'], ggfnt.GlyphPlacement{ Advance: uint8(2) })
	fontBuilder.SetGlyphPlacement(runeToUID['ƒ'], ggfnt.GlyphPlacement{ Advance: uint8(3) })
	fontBuilder.SetGlyphPlacement(runeToUID['ì'], ggfnt.GlyphPlacement{ Advance: uint8(2) })
	fontBuilder.SetGlyphPlacement(runeToUID['ï'], ggfnt.GlyphPlacement{ Advance: uint8(2) })
	fontBuilder.SetGlyphPlacement(runeToUID['î'], ggfnt.GlyphPlacement{ Advance: uint8(2) })

	// ensure that all bitmaps are being used
	for codePoint, _ := range pkgBitmaps {
		_, found := runeToUID[codePoint]
		if !found { panic("missing code point '" + string(codePoint) + "'") }
	}

	// set kerning pairs
	for _, codePoint := range ".,;:!?-–‑—" { // slightly reduce space after punctuation
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID[' '], -1)
	}
	for _, codePoint := range "-–‑—~" { // slightly reduce space before punctuation
		fontBuilder.SetKerningPair(runeToUID[' '], runeToUID[codePoint], -1)
	}
	// none of these convinced me. none looked wrong,
	// but in case of doubt we take the simplest path
	// for _, codePoint := range "FYPT" {
		//fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['.'], -1)
		//fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID[','], -1)
		//fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['_'], -1)
	// }

	fontBuilder.SetKerningPair(runeToUID['y'], runeToUID['j'], 1)
	fontBuilder.SetKerningPair(runeToUID['q'], runeToUID['j'], 1)
	fontBuilder.SetKerningPair(runeToUID['g'], runeToUID['j'], 1)
	fontBuilder.SetKerningPair(runeToUID['y'], runeToUID['ƒ'], 1)
	fontBuilder.SetKerningPair(runeToUID['q'], runeToUID['ƒ'], 1)
	fontBuilder.SetKerningPair(runeToUID['g'], runeToUID['ƒ'], 1)
	for _, codePoint := range "âädêëfƒlíïîôöñtûü0123456789AÀÁÄÂBDFHIÌÍÏÎJMNÑOÒÓÔÖǾPQRTSUVWXYZ({[|¦]})/!?€¥¨^´`\"'′‘’“”" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ì'], +1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], +1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], +1)
	}
	for _, codePoint := range "uksz¡" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], +1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], +1)
	}
	// for _, codePoint := range "íïî" {
	// 	fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ì'], +2)
	// 	fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], +2)
	// 	fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], +2)
	// }
	// for _, codePoint := range "äbêëfƒhiíjklñôötûü0123456789AÀÁÄÂBCÇDEÈÉÊËFGHIÌÍÏÎKLMNÑOÒÓÔÖǾPQRTSUVWXYZ({[|¦]})\\$!?%&¥¨^´`\"'′‘’“”" {
	// 	fontBuilder.SetKerningPair(runeToUID['í'], runeToUID[codePoint], +1)
	// 	fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID[codePoint], +1)
	// 	fontBuilder.SetKerningPair(runeToUID['î'], runeToUID[codePoint], +1)
	// }
	// for _, codePoint := range "usz€¡" {
	// 	fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID[codePoint], +1)
	// 	fontBuilder.SetKerningPair(runeToUID['î'], runeToUID[codePoint], +1)
	// }
	// for _, codePoint := range "ìïî" {
	// 	fontBuilder.SetKerningPair(runeToUID['í'], runeToUID[codePoint], +2)
	// 	fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID[codePoint], +2)
	// 	fontBuilder.SetKerningPair(runeToUID['î'], runeToUID[codePoint], +2)
	// }
	//fontBuilder.SetKerningPair(runeToUID['L'], runeToUID['Y'], -1)
	// TODO: maybe some UPPER-lower sequences could use some kerning.
	// "F{o|a|e|...}" if a good example (though it's debatable if we
	// really want this)
	for _, codePoint := range "aàáâäoòóôöeèéêënmpqg" {
		fontBuilder.SetKerningPair(runeToUID['Y'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['F'], runeToUID[codePoint], -1)
	}

	// show size
	font, err := fontBuilder.Build()
	if err != nil { panic(err) }
	err = font.Validate(ggfnt.FmtDefault)
	if err != nil { panic(err) }
	fmt.Printf("...raw size of %d bytes\n", font.RawSize())

	// export
	const FileName = "tinny-bold-6d3-v0p1.ggfnt"
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
		uid, err := fontBuilder.AddGlyph(bitmap) // *
		// * most glyphs work ok with the default placement.
		//   the remaining ones can still be adjusted afterwards
		if err != nil { panic(err) }
		err = fontBuilder.Map(codePoint, uid)
		if err != nil { panic(err) }
		codePointsMap[codePoint] = uid
	}
}

func addRunes(fontBuilder *builder.Font, codePointsMap map[rune]uint64, runes ...rune) {
	for _, codePoint := range runes {
		bitmap, found := pkgBitmaps[rune(codePoint)]
		if !found { panic("missing bitmap for '" + string(codePoint) + "'") }
		uid, err := fontBuilder.AddGlyph(bitmap) // *
		if err != nil { panic(err.Error() + " (" + string(codePoint) + ")") }
		err = fontBuilder.Map(codePoint, uid)
		if err != nil { panic(err) }
		codePointsMap[codePoint] = uid
	}
}

// helper for mask creation
func rawMask(width int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0, -height + 3, width, 3))
	copy(img.Pix, mask)
	return img
}

func rawMaskXShifted(width, xShift int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0 + xShift, -height + 3, width + xShift, 3))
	copy(img.Pix, mask)
	return img
}

var notdef = rawMask(3, []byte{
	0, 0, 0,
	1, 1, 1,
	1, 0, 1,
	1, 0, 1,
	1, 0, 1,
	1, 0, 1,
	1, 0, 1, // baseline
	1, 1, 1,
	0, 0, 0,
	0, 0, 0,
})

var pkgBitmaps = map[rune]*image.Alpha{
	// --- ascii table ---
	' ': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'!': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		0,
		1, // baseline
		0,
		0,
		0,
	}),
	'"': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'#': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 1, 0,
		1, 1, 1, 1, 1, 1,
		0, 1, 1, 0, 1, 0,
		1, 1, 1, 1, 1, 1,
		0, 1, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'$': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 0,
		1, 0, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 0, 1,
		0, 1, 1, 1, 1, 0, // baseline
		0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'%': rawMask(8, []byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 1, 1, 0, 0,
		1, 0, 1, 0, 1, 1, 0, 0,
		0, 1, 0, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 1, 0, 1, 0,
		0, 0, 1, 1, 0, 1, 0, 1,
		0, 0, 1, 1, 0, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}),
	'&': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 0, 0,
		1, 1, 0, 0, 1, 0,
		0, 1, 1, 1, 0, 0,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 0, 1, 0,
		0, 1, 1, 1, 0, 1, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'\'': rawMask(1, []byte{
		0,
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,
		0,
		0,
	}),
	'(': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,// baseline
		0, 1, 1,
		0, 0, 0,
		0, 0, 0,
	}),
	')': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1, // baseline
		1, 1, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'*': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'+': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	',': rawMask(1, []byte{
		0,
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		1,
		0,
		0,
	}),
	'-': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'.': rawMask(1, []byte{
		0,
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		0,
		0,
		0,
	}),
	'/': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'0': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'1': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 1,
		1, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'2': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 1, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'3': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 1, 1,
		0, 1, 1, 1,
		0, 0, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'4': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'5': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 0, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'6': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'7': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'8': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'9': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	':': rawMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		0,
		1, // baseline
		0,
		0,
		0,
	}),
	';': rawMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		0,
		1, // baseline
		1,
		0,
		0,
	}),
	'<': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		0, 1,
		1, 0,
		0, 1,
		0, 0, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'=': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'>': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 0,
		0, 1,
		1, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'?': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'@': rawMask(9, []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 0, 0, 0, 1, 1, 0,
		1, 1, 0, 0, 1, 1, 0, 1, 1,
		1, 1, 0, 1, 0, 1, 0, 1, 1,
		1, 1, 0, 0, 1, 1, 1, 1, 0,
		0, 1, 1, 0, 0, 0, 0, 0, 0, // baseline
		0, 0, 1, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	'A': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'B': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'C': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'D': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'E': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'F': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'G': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 1, 0, 0, 0,
		1, 1, 0, 0, 0,
		1, 1, 0, 1, 1,
		1, 1, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'H': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'I': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'J': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'K': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 0,
		1, 1, 1, 0,
		1, 1, 1, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'L': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'M': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 1,
		1, 1, 1, 0, 1, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'N': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 1, 0, 0, 1,
		1, 1, 1, 0, 1,
		1, 1, 1, 0, 1,
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		1, 1, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'O': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'P': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Q': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 1, 0, 0, 1,
		1, 1, 0, 0, 1,
		1, 1, 0, 1, 1,
		1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'R': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'S': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 1, 0, 0,
		1, 1, 1, 0,
		0, 1, 1, 1,
		0, 0, 1, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'T': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'U': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'V': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 1,
		0, 1, 1, 0, 1, 0,
		0, 1, 1, 0, 1, 0,
		0, 1, 1, 0, 1, 0,
		0, 0, 1, 1, 0, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'W': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 1, 1, 1, 1,
		0, 1, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'X': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 1,
		0, 1, 1, 0, 1, 0,
		0, 0, 1, 1, 0, 0,
		0, 1, 1, 0, 1, 0,
		1, 1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'Y': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 1,
		0, 1, 1, 0, 1, 0,
		0, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'Z': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'[': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0, // baseline
		1, 1, 1,
		0, 0, 0,
		0, 0, 0,
	}),
	'\\': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 1,
		0, 0, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	']': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1, // baseline
		1, 1, 1,
		0, 0, 0,
		0, 0, 0,
	}),
	'^': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'_': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		1, 1, 1,
		0, 0, 0,
		0, 0, 0,
	}),
	'`': rawMask(2, []byte{
		0, 0,
		0, 0,
		1, 0,
		0, 1,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'a': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 1, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'b': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'c': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 1, 1, 0,
		1, 1, 1, 0,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'd': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'e': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 0, 0,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'f': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'g': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 1, // baseline
		0, 0, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 0,
	}),
	'h': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'i': rawMask(2, []byte{
		0, 0,
		0, 0,
		1, 0,
		0, 0,
		1, 0,
		1, 1,
		1, 1, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'j': rawMaskXShifted(3, -1, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 1, 1,
		0, 1, 1, // baseline
		0, 1, 1,
		0, 1, 1,
		1, 1, 0,
	}),
	'k': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 1, 0,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'l': rawMask(2, []byte{
		0, 0,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'm': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 1, 0,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'n': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'o': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'p': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0, // baseline
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
	}),
	'q': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 1, // baseline
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
	}),
	'r': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	's': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 0,
		0, 1, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	't': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		0, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'u': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'v': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'w': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 1, 1, 1, 1,
		0, 1, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'x': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		0, 1, 1, 1,
		1, 1, 1, 0,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'y': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 1, // baseline
		0, 0, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 0,
	}),
	'z': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 1, 1, 1,
		1, 1, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'{': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 1, 1, 0, // baseline
		0, 0, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'|': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		1,
		1, // baseline
		1,
		0,
		0,
	}),
	'}': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0, // baseline
		1, 1, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'~': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		1, 0, 1, 0, 1, 
		0, 0, 0, 1, 0,  
		0, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),

	// --- additional letters for completeness ---
	'à': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'á': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'â': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ä': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'è': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 0, 0,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'é': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 0, 0,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'ê': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 1, 0, 1, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 0, 0,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'ë': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 0, 1, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 0, 0,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'ì': rawMaskXShifted(3, -1, []byte{
		0, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 1, 1,
		0, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'í': rawMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		0, 0,
		1, 0,
		1, 1,
		1, 1, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'î': rawMaskXShifted(3, -1, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 1,
		0, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'ï': rawMaskXShifted(3, -1, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 1,
		0, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'ò': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ó': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ô': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ö': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ù': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ú': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'û': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ü': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'À': rawMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Á': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Â': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ä': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'È': rawMask(5, []byte{
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 0,
		1, 1, 0, 0, 0,
		1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'É': rawMask(5, []byte{
		0, 0, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 0,
		1, 1, 0, 0, 0,
		1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'Ê': rawMask(5, []byte{
		0, 0, 1, 0, 0,
		0, 1, 0, 1, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 0,
		1, 1, 0, 0, 0,
		1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'Ë': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 0, 1, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 0,
		1, 1, 0, 0, 0,
		1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'Ì': rawMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Í': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Î': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ï': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ò': rawMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ó': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ô': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ö': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ù': rawMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ú': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Û': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ü': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),

	'Ñ': rawMask(5, []byte{
		0, 0, 1, 0, 1,
		0, 1, 0, 1, 0,
		0, 0, 0, 0, 0,
		1, 1, 0, 0, 1,
		1, 1, 1, 0, 1,
		1, 1, 0, 1, 1,
		1, 1, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'ñ': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 1,
		1, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ç': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		1, 1, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ç': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 1, 1, 0,
		1, 1, 1, 0,
		0, 1, 1, 1, // baseline
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
	}),
	'ø': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 1, 0, 1, 1,
		1, 1, 1, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'Ø': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0,
		1, 1, 0, 0, 1, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 1, 0, 0, 1,
		0, 1, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'Ǿ': rawMask(6, []byte{
		0, 0, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 0,
		0, 1, 1, 1, 1, 0,
		1, 1, 0, 0, 1, 1,
		1, 1, 0, 1, 0, 1,
		1, 1, 1, 0, 0, 1,
		0, 1, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	// 'Ø̈': rawMask(5, []byte{
	// 	0, 1, 0, 1, 0,
	// 	0, 0, 0, 0, 0,
	// 	0, 1, 1, 1, 0,
	// 	1, 0, 0, 1, 1,
	// 	1, 0, 1, 0, 1,
	// 	1, 1, 0, 0, 1,
	// 	0, 1, 1, 1, 0, // baseline
	// 	0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0,
	// }),
	'ǿ': rawMask(5, []byte{
		0, 0, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 1, 0, 1, 1,
		1, 1, 1, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	// 'ø̈': rawMask(4, []byte{
	// 	0, 0, 0, 0,
	// 	1, 0, 0, 1,
	// 	0, 0, 0, 0,
	// 	0, 1, 1, 0,
	// 	1, 0, 1, 1,
	// 	1, 1, 0, 1,
	// 	0, 1, 1, 0, // baseline
	// 	0, 0, 0, 0,
	// 	0, 0, 0, 0,
	// 	0, 0, 0, 0,
	// }),
	'Þ': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0, // baseline
		1, 1, 0, 0,
		1, 1, 0, 0,
		0, 0, 0, 0,
	}),
	'þ': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 0, // baseline
		1, 1, 0, 0,
		1, 1, 0, 0,
		0, 0, 0, 0,
	}),
	'Ð': rawMask(6, []byte{ // capital eth
		0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0,
		0, 1, 1, 0, 0, 1,
		0, 1, 1, 0, 0, 1,
		1, 1, 1, 1, 0, 1,
		0, 1, 1, 0, 0, 1,
		0, 1, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'ð': rawMask(4, []byte{ // small eth
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Đ': rawMask(6, []byte{ // capital Đ
		0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0,
		0, 1, 1, 0, 0, 1,
		0, 1, 1, 0, 0, 1,
		1, 1, 1, 1, 0, 1,
		0, 1, 1, 0, 0, 1,
		0, 1, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'đ': rawMask(4, []byte{ // small đ
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 1,
		1, 0, 1, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Æ': rawMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 0, 0,
		0, 1, 0, 1, 1, 0, 0,
		1, 1, 0, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 0, 0,
		1, 1, 0, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'æ': rawMask(8, []byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 0, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}),
	'Œ': rawMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 1, 1, 0, 0,
		1, 1, 0, 1, 1, 0, 0,
		1, 1, 0, 1, 1, 1, 0,
		1, 1, 0, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'œ': rawMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 1, 1, 0,
		1, 1, 0, 1, 1, 1, 1,
		1, 1, 0, 1, 1, 0, 0,
		0, 1, 1, 0, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'ƒ': rawMaskXShifted(5, -1, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 1, 0,
		0, 1, 1, 0, 1,
		0, 1, 1, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 1, 0, 0,
		0, 1, 1, 0, 0, // baseline
		0, 1, 1, 0, 0,
		1, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),

	// additional symbols
	'€': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 1, 0,
		0, 1, 0, 0,
		0, 0, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'¢': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 1, 1,
		1, 0, 1, 0,
		1, 0, 1, 0,
		0, 1, 1, 1,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'¥': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'£': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 1, 0,
		0, 1, 0, 0, 1,
		1, 1, 1, 0, 0,
		0, 1, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'¤': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
		0, 1, 0, 1, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		0, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),

	'¡': rawMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		1,
		1, // baseline
		1,
		1,
		0,
	}),
	'¿': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 1, // baseline
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'¦': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		0,
		1,
		1, // baseline
		1,
		0,
		0,
	}),
	'´': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'¨': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'·': rawMask(1, []byte{
		0,
		0,
		0,
		0,
		1,
		0,
		0,// baseline
		0,
		0,
		0,
	}),
	'‘': rawMask(1, []byte{ // opening single quote
		0,
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,
		0,
		0,
	}),
	'’': rawMask(1, []byte{ // closing single quote / apostrophe
		0,
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,
		0,
		0,
	}),
	'′': rawMask(1, []byte{ // prime
		0,
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,
		0,
		0,
	}),
	'“': rawMask(3, []byte{ // opening double quote
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'”': rawMask(3, []byte{ // closing double quote
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'″': rawMask(3, []byte{ // double prime
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'…': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 0, 1, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'º': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'×': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'÷': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'±': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'π': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1,
		0, 1, 1, 0, 1, 0,
		0, 1, 1, 0, 1, 0,
		0, 1, 1, 0, 0, 1, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'¬': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'–': rawMask(3, []byte{ // en dash
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'‑': rawMask(3, []byte{ // non-breaking hyphen
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'−': rawMask(3, []byte{ // minus sign
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'—': rawMask(4, []byte{ // em dash
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),

	// --- notes ---
	'♩': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		1, 1, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'♪': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 1, 0,
		0, 0, 1, 0, 1,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		1, 1, 1, 0, 0,
		1, 1, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'♫': rawMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 1,
		0, 0, 1, 1, 0, 0, 1,
		0, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, // baseline
		1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'♬': rawMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1,
		0, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 1, 1, 1, 1,
		0, 0, 1, 0, 0, 0, 1,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, // baseline
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),

	// --- special ---
	'�': rawMask(7, []byte{ // em dash
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 1, 0, 1, 1, 0,
		1, 1, 0, 1, 0, 1, 1,
		1, 1, 1, 1, 0, 1, 1,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 0, 1, 1, 0, // baseline
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	// '\uEB17': rawMask(6, []byte{
	// 	0, 0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 1, 0,
	// 	0, 0, 0, 0, 1, 1,
	// 	0, 0, 1, 1, 0, 0,
	// 	0, 1, 1, 1, 0, 0,
	// 	1, 1, 1, 0, 0, 0,
	// 	1, 1, 0, 0, 0, 0, // baseline
	// 	0, 0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0, 0,
	// }),
}
