package main

import "os"
import "fmt"
import "image"

import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt/builder"

// TODO: flag for monospaced numbers

func main() {
	// create font builder
	fmt.Print("creating new font builder\n")
	fontBuilder := builder.New()

	// add metadata
	fmt.Print("...adding metadata\n")
	err := fontBuilder.SetName("starship")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("starship")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("An uppercase only futuristic font.")
	if err != nil { panic(err) }
	err = fontBuilder.SetFirstVerDate(ggfnt.Date{ Month: 7, Day: 10, Year: 2024 })
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 2)

	// set metrics
	fmt.Print("...setting metrics\n")
	fontBuilder.SetAscent(6)
	fontBuilder.SetExtraAscent(3)
	fontBuilder.SetUppercaseAscent(6)
	fontBuilder.SetMidlineAscent(0)
	fontBuilder.SetDescent(1)
	fontBuilder.SetHorzInterspacing(1)
	fontBuilder.SetLineGap(3)
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
	addRuneRange(fontBuilder, runeToUID, ' ', '`') // ascii table until 'a'
	addRuneRange(fontBuilder, runeToUID, '{', '~') // ascii table after 'z'
	addRunes(fontBuilder, runeToUID,
		'À', 'Á', 'Â', 'Ä',
		'È', 'É', 'Ê', 'Ë',
		'Ì', 'Í', 'Î', 'Ï',
		'Ò', 'Ó', 'Ô', 'Ö',
		'Ù', 'Ú', 'Û', 'Ü',
	) // accents and diacritics
	addRunes(fontBuilder, runeToUID,
		'¡', '¿', '´', '¨', '·', '¦', '‘', '’', '“', '”', '…',
	) // additional punctuation
	addRunes(fontBuilder, runeToUID, '−', '×', '÷', '±', 'º', '′', '″', '¬', 'π') // ++maths
	addRunes(fontBuilder, runeToUID, '€', '¢', '¥', '¤') // currency symbols
	addRunes(fontBuilder, runeToUID, 'Ñ', 'Ç') // ++latin letters
	addRunes(fontBuilder, runeToUID, '–', '‑', '—') // ++dashes
	addRunes(fontBuilder, runeToUID, '♩', '♪', '♫', '♬') // notes
	// addRunes(fontBuilder, runeToUID, '�') // special

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
	fontBuilder.SetKerningPair(runeToUID['T'], runeToUID['T'], -1)
	fontBuilder.SetKerningPair(runeToUID['L'], runeToUID['T'], -1)
	fontBuilder.SetKerningPair(runeToUID['T'], runeToUID['.'], -1)
	fontBuilder.SetKerningPair(runeToUID['T'], runeToUID[','], -1)
	fontBuilder.SetKerningPair(runeToUID['T'], runeToUID['_'], -1)

	// show size
	font, err := fontBuilder.Build()
	if err != nil { panic(err) }
	err = font.Validate(ggfnt.FmtDefault)
	if err != nil { panic(err) }
	fmt.Printf("...raw size of %d bytes\n", font.RawSize())

	// export
	const FileName = "starship-6d0-v0p2.ggfnt"
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
func rawMask(width int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0, -height + 1, width, 1))
	copy(img.Pix, mask)
	return img
}

func rawMaskXShifted(width, xShift int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0 + xShift, -height + 3, width + xShift, 3))
	copy(img.Pix, mask)
	return img
}

var notdef = rawMask(4, []byte{
	1, 1, 1, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1, // baseline
	1, 1, 1, 1,
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
	}),
	'!': rawMask(1, []byte{
		1,
		1,
		1,
		1,
		0,
		1, // baseline
		0,
	}),
	'"': rawMask(3, []byte{
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'#': rawMask(5, []byte{
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'$': rawMask(5, []byte{
		0, 0, 1, 0, 0,
		1, 1, 1, 1, 1,
		1, 0, 1, 0, 0,
		1, 1, 1, 1, 1,
		0, 0, 1, 0, 1,
		1, 1, 1, 1, 1, // baseline
		0, 0, 1, 0, 0,
	}),
	'%': rawMask(7, []byte{
		0, 1, 0, 0, 1, 0, 0,
		1, 0, 1, 0, 1, 0, 0,
		0, 1, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 1, 0,
		0, 0, 1, 0, 1, 0, 1,
		0, 0, 1, 0, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0,
	}),
	'&': rawMask(4, []byte{
		0, 1, 0, 0,
		1, 0, 1, 0,
		1, 0, 1, 0,
		0, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'\'': rawMask(1, []byte{
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,
	}),
	'(': rawMask(2, []byte{
		0, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		0, 1,// baseline
		0, 0,
	}),
	')': rawMask(2, []byte{
		1, 0,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		1, 0, // baseline
		0, 0,
	}),
	'*': rawMask(3, []byte{
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'+': rawMask(3, []byte{
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	',': rawMask(1, []byte{
		1, // baseline
		1,
	}),
	'-': rawMask(2, []byte{
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
	}),
	'.': rawMask(1, []byte{
		1, // baseline
		0,
	}),
	'/': rawMask(3, []byte{
		0, 0, 1,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
	}),
	'0': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'1': rawMask(2, []byte{
		0, 1,
		1, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 0,
	}),
	'2': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'3': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'4': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'5': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'6': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'7': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'8': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'9': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	':': rawMask(1, []byte{
		1,
		0,
		0,
		1, // baseline
		0,
	}),
	';': rawMask(1, []byte{
		1,
		0,
		0,
		1, // baseline
		1,
	}),
	'<': rawMask(2, []byte{
		0, 1,
		1, 0,
		0, 1,
		0, 0, // baseline
		0, 0,
	}),
	'=': rawMask(3, []byte{
		1, 1, 1,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'>': rawMask(2, []byte{
		1, 0,
		0, 1,
		1, 0,
		0, 0, // baseline
		0, 0,
	}),
	'?': rawMask(3, []byte{
		1, 1, 1,
		1, 0, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 0,
		0, 1, 0, // baseline
		0, 0, 0,
	}),
	'@': rawMask(6, []byte{
		1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0, 0,
	}),
	'A': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'B': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 0,
		1, 0, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'C': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'D': rawMask(4, []byte{
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'E': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'F': rawMask(3, []byte{
		1, 1, 1,
		1, 0, 0,
		0, 0, 0,
		1, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
	}),
	'G': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		// 1, 1, 0, 1,
		// 1, 0, 0, 1,
		// 1, 0, 0, 0,
		// 1, 0, 1, 1,
		// 1, 0, 0, 1,
		// 1, 1, 1, 1, // baseline
		// 0, 0, 0, 0,
	}),
	'H': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'I': rawMask(1, []byte{
		1,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
	}),
	'J': rawMask(3, []byte{
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'K': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		// 1, 0, 0, 1,
		// 1, 0, 0, 1,
		// 1, 0, 0, 1,
		// 1, 1, 0, 0,
		// 1, 0, 1, 0,
		// 1, 0, 0, 1, // baseline
		// 0, 0, 0, 0,
	}),
	'L': rawMask(3, []byte{
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'M': rawMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
	}),
	'N': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'O': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'P': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 0, 0, 0,
		1, 0, 0, 0, // baseline
		0, 0, 0, 0,
	}),
	'Q': rawMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 0,
		1, 0, 0, 1, 0,
		1, 1, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
	}),
	'R': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'S': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'T': rawMask(5, []byte{
		1, 1, 1, 1, 1,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'U': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'V': rawMask(5, []byte{
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'W': rawMask(5, []byte{
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 1, 1, 1, 1, // baseline
		0, 0, 0, 0, 0,
	}),
	'X': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,// baseline
		0, 0, 0, 0,
	}),
	'Y': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 1, 0,
		0, 1, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'Z': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'[': rawMask(2, []byte{
		1, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 1, // baseline
		0, 0,
	}),
	'\\': rawMask(3, []byte{
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 1, // baseline
		0, 0, 0,
	}),
	']': rawMask(2, []byte{
		1, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		1, 1, // baseline
		0, 0,
	}),
	'^': rawMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'_': rawMask(3, []byte{
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'`': rawMask(2, []byte{
		1, 0,
		0, 1,
		0, 0,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
	}),
	'{': rawMask(3, []byte{
		0, 1, 1,
		0, 1, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 1, // baseline
		0, 0, 0,
	}),
	'|': rawMask(1, []byte{
		1,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
	}),
	'}': rawMask(3, []byte{
		1, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 1,
		0, 1, 0,
		1, 1, 0, // baseline
		0, 0, 0,
	}),
	'~': rawMask(5, []byte{
		0, 1, 0, 0, 0,
		1, 0, 1, 0, 1, 
		0, 0, 0, 1, 0,  
		0, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),

	// --- additional letters for completeness ---
	'À': rawMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'Á': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'Â': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ä': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'È': rawMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'É': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ê': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ë': rawMask(4, []byte{
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ì': rawMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
	}),
	'Í': rawMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
	}),
	'Î': rawMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
	}),
	'Ï': rawMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
	}),
	'Ò': rawMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ó': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ô': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ö': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ù': rawMask(4, []byte{
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ú': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Û': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ü': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
	}),

	'Ñ': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
	}),
	'Ç': rawMask(4, []byte{
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline
		0, 1, 0, 0,
	}),

	// additional symbols
	'€': rawMask(4, []byte{
		0, 0, 1, 1,
		0, 1, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 1, // baseline
		0, 0, 0, 0,
	}),
	'¢': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 1, 1,
		1, 0, 1, 0,
		1, 0, 1, 0,
		0, 1, 1, 1,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
	}),
	'¥': rawMask(5, []byte{
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'¤': rawMask(5, []byte{
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
		0, 1, 0, 1, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		0, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),

	'¡': rawMask(1, []byte{
		1,
		0,
		1,
		1,
		1, // baseline
		1,
	}),
	'¿': rawMask(3, []byte{
		0, 1, 0,
		0, 0, 0,
		1, 1, 0,
		1, 0, 0,
		1, 0, 1, // baseline
		1, 1, 1,
	}),
	'¦': rawMask(1, []byte{
		1,
		1,
		0,
		0,
		1,
		1, // baseline
		0,
	}),
	'´': rawMask(2, []byte{
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
	}),
	'¨': rawMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'·': rawMask(1, []byte{
		1,
		0,
		0,// baseline
		0,
	}),
	'‘': rawMask(1, []byte{ // opening single quote
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,
	}),
	'’': rawMask(1, []byte{ // closing single quote / apostrophe
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,
	}),
	'′': rawMask(1, []byte{ // prime
		1,
		1,
		0,
		0,
		0,
		0, // baseline
		0,
	}),
	'“': rawMask(3, []byte{ // opening double quote
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'”': rawMask(3, []byte{ // closing double quot
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'″': rawMask(3, []byte{ // double prime
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'…': rawMask(5, []byte{
		1, 0, 1, 0, 1, // baseline
		0, 0, 0, 0, 0,
	}),
	'º': rawMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'×': rawMask(3, []byte{
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'÷': rawMask(3, []byte{
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		0, 1, 0, // baseline
		0, 0, 0,
	}),
	'±': rawMask(3, []byte{
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'π': rawMask(5, []byte{
		1, 1, 1, 1, 1,
		0, 0, 0, 1, 0,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'¬': rawMask(4, []byte{
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 0, // baseline
		0, 0, 0, 0,
	}),
	'–': rawMask(3, []byte{ // en dash
		1, 1, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'‑': rawMask(3, []byte{ // non-breaking hyphen
		1, 1, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'−': rawMask(3, []byte{ // minus sign
		1, 1, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
	}),
	'—': rawMask(4, []byte{ // em dash
		1, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0, // baseline
		0, 0, 0, 0,
	}),

	// --- notes ---
	'♩': rawMask(3, []byte{
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		1, 1, 1,
		1, 1, 1, // baseline
		0, 0, 0,
	}),
	'♪': rawMask(5, []byte{
		0, 0, 1, 1, 0,
		0, 0, 1, 0, 1,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		1, 1, 1, 0, 0,
		1, 1, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),
	'♫': rawMask(7, []byte{
		0, 0, 0, 0, 1, 1, 1,
		0, 0, 1, 1, 0, 0, 1,
		0, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, // baseline
		1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'♬': rawMask(7, []byte{
		0, 0, 1, 1, 1, 1, 1,
		0, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 1, 1, 1, 1,
		0, 0, 1, 0, 0, 0, 1,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, // baseline
		0, 0, 0, 0, 0, 0, 0,
	}),
}
