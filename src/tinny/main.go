package main

import "os"
import "fmt"
import "image"

import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt/builder"

// TODO:
// https://en.wikipedia.org/wiki/Latin_script_in_Unicode
// (most latin-1 supplement should be added, and some
// of latin extended-A)

func main() {
	// create font builder
	fmt.Print("creating new font builder\n")
	fontBuilder := builder.New()

	// add metadata
	fmt.Print("...adding metadata\n")
	err := fontBuilder.SetName("tinny")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("tinny")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("A very low resolution, rounded and friendly-looking font with a fairly decent latin character set. This was the second ggfnt font ever created, and the first to make use of kerning features.")
	if err != nil { panic(err) }
	err = fontBuilder.SetFirstVerDate(ggfnt.Date{ Month: 6, Day: 18, Year: 2024 })
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 2)

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
	fontBuilder.SetGlyphPlacement(runeToUID['j'], ggfnt.GlyphPlacement{ Advance: uint8(1) })
	fontBuilder.SetGlyphPlacement(runeToUID['ƒ'], ggfnt.GlyphPlacement{ Advance: uint8(2) })
	fontBuilder.SetGlyphPlacement(runeToUID['ì'], ggfnt.GlyphPlacement{ Advance: uint8(1) })
	fontBuilder.SetGlyphPlacement(runeToUID['í'], ggfnt.GlyphPlacement{ Advance: uint8(1) })
	fontBuilder.SetGlyphPlacement(runeToUID['ï'], ggfnt.GlyphPlacement{ Advance: uint8(1) })
	fontBuilder.SetGlyphPlacement(runeToUID['î'], ggfnt.GlyphPlacement{ Advance: uint8(1) })

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
	for _, codePoint := range "âädêëfƒiìôöjñtûü0123456789AÀÁÄÂBDFHIÌÍÏÎJMNÑOÒÓÔÖǾPQRTSUVWXYZ({[|¦]})/!?€¥¨^´`\"'′‘’“”" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ì'], +1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], +1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], +1)
	}
	for _, codePoint := range "uksz¡" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], +1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], +1)
	}
	for _, codePoint := range "íïî" {
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ì'], +2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], +2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], +2)
	}
	for _, codePoint := range "äbêëfƒhiíjklñôötûü0123456789AÀÁÄÂBCÇDEÈÉÊËFGHIÌÍÏÎKLMNÑOÒÓÔÖǾPQRTSUVWXYZ({[|¦]})\\$!?%&¥¨^´`\"'′‘’“”" {
		fontBuilder.SetKerningPair(runeToUID['í'], runeToUID[codePoint], +1)
		fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID[codePoint], +1)
		fontBuilder.SetKerningPair(runeToUID['î'], runeToUID[codePoint], +1)
	}
	for _, codePoint := range "usz€¡" {
		fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID[codePoint], +1)
		fontBuilder.SetKerningPair(runeToUID['î'], runeToUID[codePoint], +1)
	}
	for _, codePoint := range "ìïî" {
		fontBuilder.SetKerningPair(runeToUID['í'], runeToUID[codePoint], +2)
		fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID[codePoint], +2)
		fontBuilder.SetKerningPair(runeToUID['î'], runeToUID[codePoint], +2)
	}
	fontBuilder.SetKerningPair(runeToUID['L'], runeToUID['Y'], -1)
	// TODO: maybe some UPPER-lower sequences could use some kerning.
	// "F{o|a|e|...}" if a good example (though it's debatable if we
	// really want this)

	// show size
	font, err := fontBuilder.Build()
	if err != nil { panic(err) }
	err = font.Validate(ggfnt.FmtDefault)
	if err != nil { panic(err) }
	fmt.Printf("...raw size of %d bytes\n", font.RawSize())

	// export
	const FileName = "tinny-6d3-v0p2.ggfnt"
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
		if err != nil { panic(err) }
		err = fontBuilder.Map(codePoint, uid)
		if err != nil { panic(err) }
		codePointsMap[codePoint] = uid
	}
}

// helper for mask creation
func rawAlphaMaskToWhiteMask(width int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0, -height + 3, width, 3))
	for i := 0; i < len(mask); i++ {
		img.Pix[i] = 255*mask[i]
	}
	return img
}

func rawAlphaMaskToWhiteMaskXShifted(width, xShift int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0 + xShift, -height + 3, width + xShift, 3))
	for i := 0; i < len(mask); i++ {
		img.Pix[i] = 255*mask[i]
	}
	return img
}

var notdef = rawAlphaMaskToWhiteMask(3, []byte{
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
	' ': rawAlphaMaskToWhiteMask(3, []byte{
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
	'!': rawAlphaMaskToWhiteMask(1, []byte{
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
	'"': rawAlphaMaskToWhiteMask(3, []byte{
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
	'#': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'$': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 1, 1, 1, 0,
		1, 0, 1, 0, 0,
		0, 1, 1, 1, 0,
		0, 0, 1, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 1, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'%': rawAlphaMaskToWhiteMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0,
		1, 0, 1, 0, 1, 0, 0,
		0, 1, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 1, 0,
		0, 0, 1, 0, 1, 0, 1,
		0, 0, 1, 0, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'&': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 1, 0, 0,
		1, 0, 0, 1, 0,
		0, 1, 1, 0, 0,
		1, 0, 1, 0, 1,
		1, 0, 0, 1, 0,
		0, 1, 1, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'\'': rawAlphaMaskToWhiteMask(1, []byte{
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
	'(': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0,// baseline
		0, 1,
		0, 0,
		0, 0,
	}),
	')': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		1, 0,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		1, 0,
		0, 0,
		0, 0,
	}),
	'*': rawAlphaMaskToWhiteMask(3, []byte{
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
	'+': rawAlphaMaskToWhiteMask(3, []byte{
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
	',': rawAlphaMaskToWhiteMask(1, []byte{
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
	'-': rawAlphaMaskToWhiteMask(2, []byte{
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
	'.': rawAlphaMaskToWhiteMask(1, []byte{
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
	'/': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'0': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 0, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'1': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 1,
		1, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'2': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'3': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 1, 0,
		0, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'4': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'5': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		0, 1, 1, 0,
		0, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'6': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'7': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 1, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'8': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'9': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 1,
		0, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	':': rawAlphaMaskToWhiteMask(1, []byte{
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
	';': rawAlphaMaskToWhiteMask(1, []byte{
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
	'<': rawAlphaMaskToWhiteMask(2, []byte{
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
	'=': rawAlphaMaskToWhiteMask(3, []byte{
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
	'>': rawAlphaMaskToWhiteMask(2, []byte{
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
	'?': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'@': rawAlphaMaskToWhiteMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		1, 0, 0, 1, 1, 0, 1,
		1, 0, 1, 0, 1, 0, 1,
		1, 0, 0, 1, 1, 1, 0,
		0, 1, 0, 0, 0, 0, 0, // baseline
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'A': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'B': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'C': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'D': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'E': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'F': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'G': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'H': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'I': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'J': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'K': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 0, 1, 0,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'L': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'M': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'N': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'O': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'P': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Q': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'R': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'S': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 0,
		0, 1, 1, 0,
		0, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'T': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'U': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'V': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'W': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'X': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 1, 0, 1, 0,
		1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'Y': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'Z': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'[': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		1, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		1, 1,
		0, 0,
		0, 0,
	}),
	'\\': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	']': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		1, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		1, 1,
		0, 0,
		0, 0,
	}),
	'^': rawAlphaMaskToWhiteMask(3, []byte{
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
	'_': rawAlphaMaskToWhiteMask(3, []byte{
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
	'`': rawAlphaMaskToWhiteMask(2, []byte{
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
	'a': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'b': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'c': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 1, 1,
		1, 0, 0,
		1, 0, 0,
		0, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'd': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'e': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'f': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		1, 0,
		1, 1,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'g': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'h': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'i': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		0,
		1,
		0,
		1,
		1,
		1, // baseline
		0,
		0,
		0,
	}),
	'j': rawAlphaMaskToWhiteMaskXShifted(2, -1, []byte{
		0, 0,
		0, 0,
		0, 1,
		0, 0,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 1,
		0, 1,
		1, 0,
	}),
	'k': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 1,
		1, 1, 0,
		1, 1, 0,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'l': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		0, 1, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'm': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 0, 1, 0,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'n': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'o': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'p': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
	}),
	'q': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
	}),
	'r': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	's': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		0, 1, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	't': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		1, 0,
		1, 1,
		1, 0,
		1, 0,
		1, 0,
		0, 1, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'u': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'v': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'w': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'x': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 1, 1,
		1, 1, 0,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'y': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'z': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'{': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 1,
		0, 0, 0,
		0, 0, 0,
	}),
	'|': rawAlphaMaskToWhiteMask(1, []byte{
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
	'}': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0, // baseline
		1, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'~': rawAlphaMaskToWhiteMask(5, []byte{
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
	'à': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'á': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'â': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ä': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'è': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'é': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ê': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ë': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ì': rawAlphaMaskToWhiteMaskXShifted(2, -1, []byte{
		0, 0,
		1, 0,
		0, 1,
		0, 0,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'í': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		0, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
		0, 0,
	}),
	'î': rawAlphaMaskToWhiteMaskXShifted(3, -1, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'ï': rawAlphaMaskToWhiteMaskXShifted(3, -1, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'ò': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ó': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ô': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ö': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ù': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ú': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'û': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ü': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'À': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Á': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Â': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ä': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'È': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'É': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ê': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ë': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ì': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'Í': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'Î': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'Ï': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'Ò': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ó': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ô': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ö': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ù': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ú': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Û': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ü': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),

	'Ñ': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 1,
		1, 0, 1, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 0, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ñ': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 1,
		1, 0, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ç': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 1, // baseline
		1, 1, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'ç': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 1, 1,
		1, 0, 0,
		1, 0, 0,
		0, 1, 1, // baseline
		1, 1, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'ø': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 1, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ø': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 1, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'Ǿ': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 1, 1,
		1, 0, 1, 0, 1,
		1, 1, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	// 'Ø̈': rawAlphaMaskToWhiteMask(5, []byte{
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
	'ǿ': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 1, 1,
		1, 1, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// 'ø̈': rawAlphaMaskToWhiteMask(4, []byte{
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
	'Þ': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		1, 0, 0, 0,
		1, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'þ': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		1, 0, 0, 0,
		1, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Ð': rawAlphaMaskToWhiteMask(5, []byte{ // capital eth
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 1,
		0, 1, 0, 0, 1,
		1, 1, 1, 0, 1,
		0, 1, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'ð': rawAlphaMaskToWhiteMask(4, []byte{ // small eth
		0, 0, 0, 0,
		0, 0, 0, 1,
		0, 0, 1, 1,
		0, 0, 0, 1,
		0, 1, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Đ': rawAlphaMaskToWhiteMask(5, []byte{ // capital Đ
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 1,
		0, 1, 0, 0, 1,
		1, 1, 1, 0, 1,
		0, 1, 0, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'đ': rawAlphaMaskToWhiteMask(4, []byte{ // small đ
		0, 0, 0, 0,
		0, 0, 0, 1,
		0, 0, 1, 1,
		0, 0, 0, 1,
		0, 1, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'Æ': rawAlphaMaskToWhiteMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1,
		0, 1, 0, 1, 0, 0,
		0, 1, 0, 1, 0, 0,
		1, 0, 0, 1, 1, 0,
		1, 1, 1, 1, 0, 0,
		1, 0, 0, 1, 1, 1, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'æ': rawAlphaMaskToWhiteMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		1, 0, 0, 1, 1, 1, 1,
		1, 0, 0, 1, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'Œ': rawAlphaMaskToWhiteMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1,
		1, 0, 0, 1, 0, 0,
		1, 0, 0, 1, 0, 0,
		1, 0, 0, 1, 1, 0,
		1, 0, 0, 1, 0, 0,
		0, 1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'œ': rawAlphaMaskToWhiteMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 1, 1, 0,
		1, 0, 0, 1, 1, 1, 1,
		1, 0, 0, 1, 0, 0, 0,
		0, 1, 1, 0, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'ƒ': rawAlphaMaskToWhiteMaskXShifted(4, -1, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 1,
		0, 1, 0, 0,
		0, 1, 1, 0,
		0, 1, 0, 0,
		0, 1, 0, 0, // baseline
		0, 1, 0, 0,
		1, 0, 0, 0,
		0, 0, 0, 0,
	}),

	// additional symbols
	'€': rawAlphaMaskToWhiteMask(4, []byte{
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
	'¢': rawAlphaMaskToWhiteMask(4, []byte{
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
	'¥': rawAlphaMaskToWhiteMask(5, []byte{
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
	'£': rawAlphaMaskToWhiteMask(5, []byte{
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
	'¤': rawAlphaMaskToWhiteMask(5, []byte{
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

	'¡': rawAlphaMaskToWhiteMask(1, []byte{
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
	'¿': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 1, // baseline
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'¦': rawAlphaMaskToWhiteMask(1, []byte{
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
	'´': rawAlphaMaskToWhiteMask(2, []byte{
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
	'¨': rawAlphaMaskToWhiteMask(3, []byte{
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
	'·': rawAlphaMaskToWhiteMask(1, []byte{
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
	'‘': rawAlphaMaskToWhiteMask(1, []byte{ // opening single quote
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
	'’': rawAlphaMaskToWhiteMask(1, []byte{ // closing single quote / apostrophe
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
	'′': rawAlphaMaskToWhiteMask(1, []byte{ // prime
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
	'“': rawAlphaMaskToWhiteMask(3, []byte{ // opening double quote
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
	'”': rawAlphaMaskToWhiteMask(3, []byte{ // closing double quote
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
	'″': rawAlphaMaskToWhiteMask(3, []byte{ // double prime
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
	'…': rawAlphaMaskToWhiteMask(5, []byte{
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
	'º': rawAlphaMaskToWhiteMask(3, []byte{
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
	'×': rawAlphaMaskToWhiteMask(3, []byte{
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
	'÷': rawAlphaMaskToWhiteMask(3, []byte{
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
	'±': rawAlphaMaskToWhiteMask(3, []byte{
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
	'π': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0,
		0, 1, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'¬': rawAlphaMaskToWhiteMask(4, []byte{
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
	'–': rawAlphaMaskToWhiteMask(3, []byte{ // en dash
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
	'‑': rawAlphaMaskToWhiteMask(3, []byte{ // non-breaking hyphen
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
	'−': rawAlphaMaskToWhiteMask(3, []byte{ // minus sign
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
	'—': rawAlphaMaskToWhiteMask(4, []byte{ // em dash
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
	'♩': rawAlphaMaskToWhiteMask(3, []byte{
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
	'♪': rawAlphaMaskToWhiteMask(5, []byte{
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
	'♫': rawAlphaMaskToWhiteMask(7, []byte{
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
	'♬': rawAlphaMaskToWhiteMask(7, []byte{
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
	'�': rawAlphaMaskToWhiteMask(7, []byte{ // em dash
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
	// '\uEB17': rawAlphaMaskToWhiteMask(6, []byte{
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
