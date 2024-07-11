package main

import "os"
import "fmt"
import "image"
import "image/color"

import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt/builder"

func main() {
	// create font builder
	fmt.Print("creating new font builder\n")
	fontBuilder := builder.New()

	// add metadata
	fmt.Print("...adding metadata\n")
	err := fontBuilder.SetName("candy")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("candy")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("Comfy and happy display font.")
	if err != nil { panic(err) }
	err = fontBuilder.SetFirstVerDate(ggfnt.Date{ Month: 7, Day: 10, Year: 2024 })
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 1)

	// set metrics
	fmt.Print("...setting metrics\n")
	fontBuilder.SetAscent(10)
	fontBuilder.SetExtraAscent(2)
	fontBuilder.SetUppercaseAscent(10)
	fontBuilder.SetMidlineAscent(0)
	fontBuilder.SetDescent(2)
	fontBuilder.SetHorzInterspacing(1) // it's perfectly usable with 2 too, depends on preferences
	fontBuilder.SetLineGap(4)
	err = fontBuilder.GetMetricsStatus()
	if err != nil { panic(err) }

	// set main dye and color palette
	err = fontBuilder.AddDye("main", 255)
	if err != nil { panic(err) }
	var RGB = func(r, g, b uint8) color.RGBA { return color.RGBA{r, g, b, 255} }
	err = fontBuilder.AddPalette("candy",
		RGB(186, 31, 147), RGB(255, 43, 202), RGB(255, 112, 219), // magentas, from dark to light
		RGB(48, 181, 181), RGB(61, 204, 204), RGB(82, 221, 221), RGB(219, 244, 255), // cyans, from dark to light
	)
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
	// addRunes(fontBuilder, runeToUID,
	// 	'À', 'Á', 'Â', 'Ä',
	// 	'È', 'É', 'Ê', 'Ë',
	// 	'Ì', 'Í', 'Î', 'Ï',
	// 	'Ò', 'Ó', 'Ô', 'Ö',
	// 	'Ù', 'Ú', 'Û', 'Ü',
	// ) // accents and diacritics
	addRunes(fontBuilder, runeToUID,
		'¡', '¿', '´', '¨', '·', '¦', '‘', '’', '“', '”', '•', '…', '\uED01', // ED01 is custom fat dot
	) // additional punctuation
	// addRunes(fontBuilder, runeToUID, '−', '×', '÷', '±', 'º', '′', '″', '¬', 'π') // ++maths
	// addRunes(fontBuilder, runeToUID, '€', '¢', '¥', '¤') // currency symbols
	addRunes(fontBuilder, runeToUID, 'Ñ', 'Ç') // ++latin letters
	addRunes(fontBuilder, runeToUID, '–', '‑', '—') // ++dashes
	//addRunes(fontBuilder, runeToUID, '♩', '♪', '♫', '♬') // notes
	//addRunes(fontBuilder, runeToUID, '❤', '💔', '�') // special
	
	err = fontBuilder.SetGlyphName(runeToUID['\uED01'], "fat-dot")
	if err != nil { panic(err) }

	// map candy icons
	candyUID, err := fontBuilder.AddGlyph(candy)
	if err != nil { panic(err) }
	err = fontBuilder.SetGlyphName(candyUID, "candy")
	if err != nil { panic(err) }
	err = fontBuilder.Map('🍬', candyUID) // '\u1F36C'
	if err != nil { panic(err) }
	monocandyUID, err := fontBuilder.AddGlyph(monocandy)
	if err != nil { panic(err) }
	err = fontBuilder.SetGlyphName(monocandyUID, "dye-candy")
	if err != nil { panic(err) }
	err = fontBuilder.Map('\uEDCA', monocandyUID)
	if err != nil { panic(err) }

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
	fontBuilder.SetKerningPair(runeToUID['.'], runeToUID['!'], -1)
	fontBuilder.SetKerningPair(runeToUID['.'], runeToUID['?'], -2)
	fontBuilder.SetKerningPair(runeToUID['…'], runeToUID['!'], -1)
	fontBuilder.SetKerningPair(runeToUID['…'], runeToUID['?'], -2)
	fontBuilder.SetKerningPair(runeToUID['L'], runeToUID['Y'], -1)
	fontBuilder.SetKerningPair(runeToUID['L'], runeToUID['?'], -1)
	fontBuilder.SetKerningPair(runeToUID['Y'], runeToUID['.'], -1)
	fontBuilder.SetKerningPair(runeToUID['T'], runeToUID['.'], -1)
	fontBuilder.SetKerningPair(runeToUID['Y'], runeToUID[','], -1)
	fontBuilder.SetKerningPair(runeToUID['T'], runeToUID[','], -1)
	fontBuilder.SetKerningPair(runeToUID['Y'], runeToUID['…'], -1)
	fontBuilder.SetKerningPair(runeToUID['T'], runeToUID['…'], -1)
	fontBuilder.SetKerningPair(runeToUID['•'], runeToUID[' '], -2)

	// show size
	font, err := fontBuilder.Build()
	if err != nil { panic(err) }
	err = font.Validate(ggfnt.FmtDefault)
	if err != nil { panic(err) }
	fmt.Printf("...raw size of %d bytes\n", font.RawSize())

	// export
	const FileName = "candy-10d2-v0p1.ggfnt"
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
func rawMask(width, descent int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0, -height + descent, width, descent))
	copy(img.Pix, mask)
	return img
}

var notdef = rawMask(6, 0, []byte{
	1, 1, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 1, 
	1, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 1,
	1, 1, 1, 1, 1, 1,
})

var candy = rawMask(12, 1, []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 4, 3, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 4, 3, 3, 2, 0,
	0, 0, 0, 0, 7, 7, 7, 6, 3, 3, 3, 2,
	0, 0, 0, 7, 8, 8, 7, 7, 6, 2, 2, 2,
	0, 0, 7, 8, 7, 7, 7, 7, 6, 6, 0, 0,
	0, 0, 7, 7, 7, 7, 7, 7, 6, 6, 0, 0,
	0, 0, 6, 7, 7, 7, 7, 6, 6, 5, 0, 0,
	4, 4, 3, 6, 6, 6, 6, 6, 5, 0, 0, 0,
	4, 3, 3, 2, 5, 5, 5, 5, 0, 0, 0, 0,
	0, 3, 3, 2, 2, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0,
})

var monocandy = rawMask(12, 1, []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0,
	0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1,
	0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0,
	0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
	0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
})

var pkgBitmaps = map[rune]*image.Alpha{
	// --- ascii table ---
	' ': rawMask(6, 0, []byte{
		0, 0, 0, 0, 0, 0,
	}),
	'!': rawMask(4, 0, []byte{
		0, 1, 1, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
	}),
	'"': rawMask(5, 0, []byte{
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		0, 1, 0, 0, 1,
		1, 0, 0, 1, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'#': rawMask(12, 0, []byte{
		0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0,
	}),
	'$': rawMask(8, 0, []byte{
		0, 0, 0, 1, 1, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 1, 1, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 1, 1, 0, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 1, 1, 0, 0, 0,
	}),
	'%': rawMask(14, 0, []byte{
		0, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0,
		// 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0,
		// 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0,
		// 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0,
		// 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0,
		// 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1,
		// 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1,
		// 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0,
	}),
	'&': rawMask(8, 0, []byte{
		0, 0, 1, 1, 1, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0, 0,
		0, 1, 1, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0, 0,
		0, 1, 1, 0, 1, 1, 0, 1,
		1, 1, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 0, 1,
	}),
	'\'': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
	}),
	'(': rawMask(3, 0, []byte{
		0, 1, 1,
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
		0, 1, 1,
	}),
	')': rawMask(3, 0, []byte{
		1, 1, 0,
		1, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		1, 1, 1,
		1, 1, 0,
	}),
	'*': rawMask(6, 0, []byte{
		0, 0, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 1, 1,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		// 1, 0, 1, 0, 1,
		// 0, 1, 1, 1, 0,
		// 0, 1, 1, 1, 0,
		// 1, 0, 1, 0, 1,
		// 0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0,
		// 0, 0, 0, 0, 0,
	}),
	'+': rawMask(6, 0, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1,
		0, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
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
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
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
		0, 0, 1, 1,
		0, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
	}),
	'0': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'1': rawMask(5, 0, []byte{
		0, 0, 1, 1, 0,
		0, 1, 1, 1, 0,
		1, 1, 1, 1, 0,
		1, 1, 1, 1, 0,
		0, 1, 1, 1, 0,
		0, 1, 1, 1, 0,
		0, 1, 1, 1, 0,
		0, 1, 1, 1, 0,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
	}),
	'2': rawMask(7, 0, []byte{
		0, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 1, 1,
		0, 0, 0, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 0,
	}),
	'3': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 1, 1, 1, 1, 0,
		0, 0, 0, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'4': rawMask(8, 0, []byte{
		0, 1, 1, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 0,
	}),
	'5': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'6': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'7': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 1, 1, 1, 1,
		0, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0, 0,
	}),
	'8': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'9': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 0,
	}),
	':': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		1, 1,
		1, 1,
	}),
	';': rawMask(2, 2, []byte{
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		1, 1,
		1, 1,
		0, 1,
		1, 0,
	}),
	'<': rawMask(5, 0, []byte{
		0, 0, 0, 1, 1,
		0, 0, 1, 1, 1,
		0, 1, 1, 1, 0,
		1, 1, 1, 0, 0,
		1, 1, 1, 0, 0,
		0, 1, 1, 1, 0,
		0, 0, 1, 1, 1,
		0, 0, 0, 1, 1,
		0, 0, 0, 0, 0,
	}),
	'=': rawMask(6, 0, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'>': rawMask(5, 0, []byte{
		1, 1, 0, 0, 0,
		1, 1, 1, 0, 0,
		0, 1, 1, 1, 0,
		0, 0, 1, 1, 1,
		0, 0, 1, 1, 1,
		0, 1, 1, 1, 0,
		1, 1, 1, 0, 0,
		1, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'?': rawMask(7, 0, []byte{
		0, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 1, 1, 1,
		0, 0, 0, 1, 1, 1, 1,
		0, 0, 0, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0,
	}),
	'@': rawMask(10, 0, []byte{
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 1, 1, 0, 0, 1, 1,
		1, 1, 0, 1, 0, 0, 1, 0, 1, 1,
		1, 1, 0, 1, 0, 0, 1, 0, 1, 1,
		1, 1, 0, 0, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
	}),
	'A': rawMask(9, 0, []byte{
		0, 0, 1, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 0, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 1, 1, 0,
	}),
	'B': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'C': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'D': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 1, 1, 0, 0,
	}),
	'E': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 0, 0, 0,
		1, 1, 1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 0, 0, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'F': rawMask(7, 0, []byte{
		0, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 0,
	}),
	'G': rawMask(9, 0, []byte{
		0, 0, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 0, 0, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'H': rawMask(9, 0, []byte{
		0, 1, 1, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 1, 1, 0,
	}),
	'I': rawMask(4, 0, []byte{
		0, 1, 1, 0,
		1, 1, 1, 1,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 0,
	}),
	'J': rawMask(7, 0, []byte{
		0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 0,
	}),
	'K': rawMask(8, 0, []byte{
		0, 1, 1, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, 1,
		0, 1, 1, 0, 0, 1, 1, 0,
	}),
	'L': rawMask(6, 0, []byte{
		0, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0,
		1, 1, 1, 1, 0, 0,
		1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1,
	}),
	'M': rawMask(10, 0, []byte{
		0, 1, 1, 0, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 1, 0, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 0, 1, 1, 0,
	}),
	'N': rawMask(9, 0, []byte{
		0, 1, 1, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 1, 1, 0,
	}),
	'O': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'P': rawMask(7, 0, []byte{
		0, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 0,
	}),
	'Q': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 0, 0, 1,
		1, 1, 1, 0, 0, 1, 0, 1,
		1, 1, 1, 0, 0, 1, 1, 0,
		0, 1, 1, 1, 1, 0, 1, 1,
	}),
	'R': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		0, 1, 1, 0, 0, 1, 1, 0,
	}),
	'S': rawMask(8, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 0, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'T': rawMask(9, 0, []byte{
		0, 1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 1, 1, 1, 0, 1, 1,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
	}),
	'U': rawMask(8, 0, []byte{
		0, 1, 1, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
	}),
	'V': rawMask(9, 0, []byte{
		0, 1, 1, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		// 0, 1, 0, 0, 0, 0, 0, 1, 0,
		// 1, 1, 0, 0, 0, 0, 0, 1, 1,
		// 1, 1, 0, 0, 0, 0, 0, 1, 1,
		// 1, 1, 0, 0, 0, 0, 0, 1, 1,
		// 1, 1, 0, 0, 0, 0, 0, 1, 1,
		// 1, 1, 1, 0, 0, 0, 1, 1, 1,
		// 1, 1, 1, 0, 0, 0, 1, 1, 1,
		// 0, 1, 1, 1, 0, 1, 1, 1, 0,
		// 0, 1, 1, 1, 1, 1, 1, 1, 0,
		// 0, 0, 1, 1, 1, 1, 1, 0, 0,
	}),
	'W': rawMask(10, 0, []byte{
		0, 1, 1, 0, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 0, 1, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 0, 1, 1, 0,
	}),
	'X': rawMask(9, 0, []byte{
		0, 1, 1, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 1, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 1, 1, 0,
	}),
	'Y': rawMask(8, 0, []byte{
		0, 1, 1, 0, 0, 1, 1, 0,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 0, 0, 0,
	}),
	'Z': rawMask(7, 0, []byte{
		0, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 1, 1, 1,
		0, 0, 0, 1, 1, 1, 1,
		0, 0, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 1, 0, 0,
		1, 1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 0,
	}),
	'[': rawMask(3, 0, []byte{
		1, 1, 1,
		1, 1, 1,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 0,
		1, 1, 1,
		1, 1, 1,
	}),
	'\\': rawMask(4, 0, []byte{
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
	}),
	']': rawMask(3, 0, []byte{
		1, 1, 1,
		1, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		0, 1, 1,
		1, 1, 1,
		1, 1, 1,
	}),
	'^': rawMask(6, 0, []byte{
		0, 0, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 0,
		1, 1, 0, 0, 1, 1,
		1, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'_': rawMask(7, 1, []byte{
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
	}),
	'`': rawMask(3, 0, []byte{
		1, 1, 0,
		0, 1, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'{': rawMask(4, 0, []byte{
		0, 0, 1, 1,
		0, 1, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 0, 0, 
		1, 1, 0, 0, 
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 1,
		0, 0, 1, 1,
	}),
	'|': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
	}),
	'}': rawMask(4, 0, []byte{
		1, 1, 0, 0,
		1, 1, 1, 0,
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 1, 1, 0,
		0, 1, 1, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,
	}),
	'~': rawMask(7, 0, []byte{
		0, 1, 1, 1, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),

	// Ñ and Ç
	'Ñ': rawMask(9, 0, []byte{
		0, 0, 1, 1, 1, 0, 0, 1, 0,
		0, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 0, 1, 1, 0,
		1, 1, 1, 1, 0, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 1, 1, 1, 1,
		0, 1, 1, 0, 0, 0, 1, 1, 0,
	}),
	'Ç': rawMask(8, 2, []byte{
		0, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0, 1, 1,
		1, 1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0,
	}),

	// extra punctuation
	'¡': rawMask(4, 2, []byte{
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 0,
	}),
	'¿': rawMask(7, 0, []byte{
		0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 0, 0, 0,
		1, 1, 1, 1, 0, 0, 0,
		1, 1, 1, 0, 0, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 0,
	}),
	'´': rawMask(3, 0, []byte{
		0, 1, 1,
		1, 1, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'¨': rawMask(5, 0, []byte{
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'·': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
	}),
	'¦': rawMask(2, 0, []byte{
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
	}),
	'‘': rawMask(2, 0, []byte{
		0, 1,
		1, 0,
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		0, 0,
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
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'•': rawMask(4, 0, []byte{
		0, 1, 1, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'…': rawMask(8, 0, []byte{
		1, 1, 0, 1, 1, 0, 1, 1,
		1, 1, 0, 1, 1, 0, 1, 1,
	}),
	'\uED01': rawMask(4, 0, []byte{
		0, 1, 1, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 1, 1, 0,
	}),

	// dashes
	'–': rawMask(4, 0, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'‑': rawMask(4, 0, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'—': rawMask(7, 0, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
}