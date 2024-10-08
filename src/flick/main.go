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
	err := fontBuilder.SetName("flick")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("flick")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("The first font created to experiment with glyph animation on ggfnt. It's a very silly font otherwise, uppercase-only and with many glyphs being quite similar or identical to 'tinny'.")
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 2)
	creationDate := ggfnt.Date{ Day: 1, Month: 7, Year: 2024 }
	err = fontBuilder.SetFirstVerDate(creationDate)
	if err != nil { panic(err) }
	err = fontBuilder.SetMajorVerDate(creationDate)
	if err != nil { panic(err) }

	// set metrics
	fmt.Print("...setting metrics\n")
	fontBuilder.SetAscent(6)
	fontBuilder.SetExtraAscent(0)
	fontBuilder.SetUppercaseAscent(6)
	fontBuilder.SetMidlineAscent(0)
	fontBuilder.SetDescent(1)
	fontBuilder.SetExtraDescent(0)
	fontBuilder.SetHorzInterspacing(1)
	fontBuilder.SetLineGap(2)
	err = fontBuilder.GetMetricsStatus()
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
	// addRunes(fontBuilder, runeToUID,
	// 	'À', 'Á', 'Â', 'Ä', 'à', 'á', 'â', 'ä',
	// 	'È', 'É', 'Ê', 'Ë', 'è', 'é', 'ê', 'ë',
	// 	'Ì', 'Í', 'Î', 'Ï', 'ì', 'í', 'î', 'ï',
	// 	'Ò', 'Ó', 'Ô', 'Ö', 'ò', 'ó', 'ô', 'ö',
	// 	'Ù', 'Ú', 'Û', 'Ü', 'ù', 'ú', 'û', 'ü',
	// ) // accents and diacritics
	// addRunes(fontBuilder, runeToUID, '¡', '¿', '´', '¨', '·', '¦') // additional punctuation
	// addRunes(fontBuilder, runeToUID, '−', '×', '÷', '±', 'º', '¬') // ++maths
	// addRunes(fontBuilder, runeToUID, '€', '£', '¢', '¥', '¤') // currency symbols
	// addRunes(fontBuilder, runeToUID, 'Ñ', 'ñ', 'Ç', 'ç') // ++spanish letters
	// addRunes(fontBuilder, runeToUID, '–', '‑', '—', '\uE001') // ++dashes
	// addRunes(fontBuilder, runeToUID, '♩', '♪', '♫') // notes
	// addRunes(fontBuilder, runeToUID, ' ', ' ') // thin space and hair space for padding
	// addRunes(fontBuilder, runeToUID, '◀', '▶', '❤', '💔') // special symbols

	// show size
	font, err := fontBuilder.Build()
	if err != nil { panic(err) }
	err = font.Validate(ggfnt.FmtDefault)
	if err != nil { panic(err) }
	fmt.Printf("...raw size of %d bytes\n", font.RawSize())

	// export
	const FileName = "flick-6d0-v0p2.ggfnt"
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

		altBitmap, altFound := altBitmaps[codePoint]
		if altFound {
			altUID, err := fontBuilder.AddGlyph(altBitmap)
			if err != nil { panic(err) }
			err = fontBuilder.MapGroup(codePoint, ggfnt.AnimFlagLoopable | ggfnt.AnimFlagSequential, uid, altUID)
			if err != nil { panic(err) }
		} else {
			err = fontBuilder.Map(codePoint, uid)
			if err != nil { panic(err) }
		}
		
		codePointsMap[codePoint] = uid
	}
}

func addRunes(fontBuilder *builder.Font, codePointsMap map[rune]uint64, runes ...rune) {
	for _, codePoint := range runes {
		bitmap, found := pkgBitmaps[codePoint]
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
	img := image.NewAlpha(image.Rect(0, -6, width, -6 + height))
	copy(img.Pix, mask)
	return img
}

var notdef = rawMask(3, []byte{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
})

var pkgBitmaps = map[rune]*image.Alpha{
	' ': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'A': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
	}),
	'B': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 1,
	}),
	'C': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 0,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'D': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0,
	}),
	'E': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1,
	}),
	'F': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
	}),
	'G': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'H': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
	}),
	'I': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	}),
	'J': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'K': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 1, 0,
		1, 1, 0, 0,
		1, 0, 1, 0,
		1, 0, 0, 1,
	}),
	'L': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
	}),
	'M': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
	}),
	'N': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 0, 1, 1,
		1, 0, 0, 1,
	}),
	'O': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'P': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 0,
	}),
	'Q': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		0, 1, 1, 1,
	}),
	'R': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
	}),
	'S': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 1,
		1, 0, 0, 0,
		0, 1, 1, 0,
		0, 0, 0, 1,
		1, 1, 1, 0,
	}),
	'T': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	}),
	'U': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'V': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0,
		0, 0, 1, 0, 0,
	}),
	'W': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 0, 1, 0,
	}),
	'X': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		1, 0, 1,
	}),
	'Y': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	}),
	'Z': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 0, 1,
		0, 1, 0,
		1, 0, 0,
		1, 1, 1,
	}),

	'.': rawMask(1, []byte{
		0,
		0,
		0,
		0,
		0,
		1,
	}),
	',': rawMask(1, []byte{
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		1,
	}),
	':': rawMask(1, []byte{
		0,
		0,
		1,
		0,
		0,
		1, // baseline
		0,
	}),
	';': rawMask(1, []byte{
		0,
		0,
		1,
		0,
		0,
		1, // baseline
		1,
	}),
	'!': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		0,
		1,
	}),
	'¡': rawMask(1, []byte{
		0,
		0,
		1,
		0,
		1,
		1, // baseline
		1,
	}),
	'?': rawMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
	}),
	'¿': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		1, 0, 0,
		1, 0, 1, // baseline
		0, 1, 0,
	}),

	// remaining ASCII except lowercase
	'"': rawMask(3, []byte{
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'#': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 1, 0, 1, 0,
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
		0, 1, 0, 0, 0, 0, 0,
		1, 0, 1, 0, 1, 0, 0,
		0, 1, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 1, 0,
		0, 0, 1, 0, 1, 0, 1, // baseline
		0, 0, 0, 0, 0, 1, 0,
	}),
	'&': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		1, 0, 1, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
	}),
	'\'': rawMask(1, []byte{
		1,
		1,
		0,
		0,
		0,
		0,
	}),
	'(': rawMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		1, 0,
		1, 0,
		0, 1,
	}),
	')': rawMask(2, []byte{
		0, 0,
		1, 0,
		0, 1,
		0, 1,
		0, 1,
		1, 0,
	}),
	'*': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 1, 0, 
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
	}),
	'+': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 1, 0, 
		1, 1, 1,
		0, 1, 0,
		0, 0, 0,
	}),
	'-': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0,
	}),
	'/': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 1, 0,
		1, 0, 0,
		1, 0, 0,
	}),
	'0': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'1': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	}),
	'2': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		1, 1, 1, 1,
	}),
	'3': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 1,
		0, 1, 1, 0,
		0, 0, 0, 1,
		1, 1, 1, 0,
	}),
	'4': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 1, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		0, 0, 0, 1,
	}),
	'5': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 1,
		1, 1, 1, 0,
	}),
	'6': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'7': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 1, 0, 0,
	}),
	'8': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'9': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		0, 1, 1, 1,
		0, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'<': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 1,
		1, 0,
		0, 1,
		0, 0,
	}),
	'=': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
	}),
	'>': rawMask(2, []byte{
		0, 0,
		0, 0,
		1, 0,
		0, 1,
		1, 0,
		0, 0,
	}),
	'@': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0,
		0, 1, 0, 0, 1, 0,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 1, 1,
		0, 1, 0, 0, 0, 0, // baseline
		0, 0, 1, 1, 0, 0,
	}),
	'[': rawMask(2, []byte{
		0, 0,
		1, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 1,
	}),
	'\\': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 1,
	}),
	']': rawMask(2, []byte{
		0, 0,
		1, 1,
		0, 1,
		0, 1,
		0, 1,
		1, 1,
	}),
	'^': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
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
		0, 0, 0, // baseline
		1, 1, 1,
	}),
	'`': rawMask(2, []byte{
		0, 0,
		1, 0,
		0, 1,
		0, 0,
		0, 0,
		0, 0,
	}),
	'´': rawMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0,
	}),
	'¨': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}),
	'{': rawMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 1,
	}),
	'|': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		1,
	}),
	'¦': rawMask(1, []byte{
		0,
		1,
		1,
		0,
		1,
		1,
	}),
	'·': rawMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		0,
	}),
	'}': rawMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0, // baseline
		1, 0, 0,
	}),
	'~': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		1, 0, 1, 0, 1,
		0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
	}),
}

// --- alter forms ---
var altBitmaps = map[rune]*image.Alpha{
	'A': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
	}),
	'B': rawMask(4, []byte{
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
	}),
	'C': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'D': rawMask(4, []byte{
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0,
	}),
	'E': rawMask(3, []byte{
		1, 1, 1,
		1, 0, 0,
		1, 1, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
	}),
	'F': rawMask(3, []byte{
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
	}),
	'G': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'H': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
	}),
	'I': rawMask(3, []byte{
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	}),
	'J': rawMask(4, []byte{
		0, 0, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'K': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 1, 0,
		1, 1, 0, 0,
		1, 0, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
	}),
	'L': rawMask(3, []byte{
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
	}),
	'M': rawMask(5, []byte{
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
	}),
	'N': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 1, 0, 1,
		1, 0, 1, 1,
		1, 0, 1, 1,
		1, 0, 0, 1,
	}),
	'O': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'P': rawMask(4, []byte{
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
	}),
	'Q': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		0, 1, 1, 1,
	}),
	'R': rawMask(4, []byte{
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
	}),
	'S': rawMask(4, []byte{
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 0, 0, 0,
		0, 1, 1, 0,
		0, 0, 0, 1,
		1, 1, 1, 0,
	}),
	'T': rawMask(3, []byte{
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	}),
	'U': rawMask(4, []byte{
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}),
	'V': rawMask(5, []byte{
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0,
		0, 0, 1, 0, 0,
	}),
	'W': rawMask(5, []byte{
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 0, 1, 0,
	}),
	'X': rawMask(3, []byte{
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		1, 0, 1,
	}),
	'Y': rawMask(3, []byte{
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	}),
	'Z': rawMask(3, []byte{
		1, 1, 1,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		1, 1, 1,
	}),
}
