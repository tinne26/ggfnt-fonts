package main

import "os"
import "fmt"
import "image"

import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt/builder"

// TODO:
// - zero-disambiguation-mark: on, off
// - numeric-style: clear, neutral, compact
// - maybe animated cursor? either as vertical bar or underscore
// switches:
// - zero-switch: zero-disambiguation-mark + numeric-style
// - numeric style: that

// TODO: add named gamepad keys and so on? hmmm...

// globals
var SwitchZeroKey, SwitchNumStyleKey uint8

func main() {
	// create font builder
	fmt.Print("creating new font builder\n")
	fontBuilder := builder.New()

	// add metadata
	fmt.Print("...adding metadata\n")
	err := fontBuilder.SetName("jammy")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("jammy")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("This font was born from tinne's entries for Ebitengine game jams. Originally, a few pixel art letters were made for small parts of the UI in Bindless (2022). The next year, many more ASCII characters were added for Transition. When developing ggfnt, this was the first font to be created and exported for the format. Right before public release, it was expanded to include lowercase characters and complete the full ASCII range.")
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 3)
	creationDate := ggfnt.Date{ Month: 6, Year: 2022 }
	err = fontBuilder.SetFirstVerDate(creationDate)
	if err != nil { panic(err) }
	err = fontBuilder.SetMajorVerDate(creationDate)
	if err != nil { panic(err) }

	// set metrics
	fmt.Print("...setting metrics\n")
	fontBuilder.SetAscent(6)
	fontBuilder.SetExtraAscent(1) // only for accents
	fontBuilder.SetUppercaseAscent(5)
	fontBuilder.SetMidlineAscent(3)
	fontBuilder.SetDescent(2)
	fontBuilder.SetHorzInterspacing(1)
	fontBuilder.SetLineGap(0)
	err = fontBuilder.GetMetricsStatus()
	if err != nil { panic(err) }

	// TODO:
	// add two settings. then two switch types. then we can map all that without issue.
	settingZeroDisKey , err := fontBuilder.AddSetting("zero-disambiguation-mark", "on", "off")
	if err != nil { panic(err) }
	settingNumStyleKey, err := fontBuilder.AddSetting("numeric-style", "clear", "neutral", "compact")
	if err != nil { panic(err) }
	SwitchZeroKey, err = fontBuilder.AddSwitch(settingZeroDisKey, settingNumStyleKey)
	if err != nil { panic(err) }
	SwitchNumStyleKey, err = fontBuilder.AddSwitch(settingNumStyleKey)
	if err != nil { panic(err) }

	// add notdef as the first glyph
	fmt.Printf("...registering glyphs\n")
	notdefUID, err := fontBuilder.AddGlyph(notdef)
	if err != nil { panic(err) }
	err = fontBuilder.SetGlyphName(notdefUID, "notdef")
	if err != nil { panic(err) }
	err = fontBuilder.Map('\uE000', notdefUID)
	if err != nil { panic(err) }

	// add all other glyphs
	runeToUID := make(map[rune]uint64, 128)
	addRuneRange(fontBuilder, runeToUID, ' ', '~') // ASCII
	addRunes(fontBuilder, runeToUID,
		'À', 'Á', 'Â', 'Ä', 'à', 'á', 'â', 'ä',
		'È', 'É', 'Ê', 'Ë', 'è', 'é', 'ê', 'ë',
		'Ì', 'Í', 'Î', 'Ï', 'ì', 'í', 'î', 'ï',
		'Ò', 'Ó', 'Ô', 'Ö', 'ò', 'ó', 'ô', 'ö',
		'Ù', 'Ú', 'Û', 'Ü', 'ù', 'ú', 'û', 'ü',
	) // accents and diacritics
	addRunes(fontBuilder, runeToUID, '¡', '¿', '´', '¨', '·', '¦') // additional punctuation
	addRunes(fontBuilder, runeToUID, '−', '×', '÷', '±', 'º', '¬') // ++maths
	addRunes(fontBuilder, runeToUID, '€', '£', '¢', '¥', '¤') // currency symbols
	addRunes(fontBuilder, runeToUID, 'Ñ', 'ñ', 'Ç', 'ç') // ++spanish letters
	addRunes(fontBuilder, runeToUID, '–', '‑', '—', '\uE001') // ++dashes
	addRunes(fontBuilder, runeToUID, '♩', '♪', '♫') // notes
	addRunes(fontBuilder, runeToUID, ' ', ' ') // thin space and hair space for padding
	addRunes(fontBuilder, runeToUID, '◀', '▶', '❤', '💔') // special symbols

	// set kerning pairs
	fmt.Printf("...configuring kerning pairs\n")
	for _, codePoint := range ".,;:!?" { // slightly reduce space after punctuation
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID[' '], -1)
	}
	fontBuilder.SetKerningPair(runeToUID['s'], runeToUID['s'], -1)
	fontBuilder.SetKerningPair(runeToUID['z'], runeToUID['z'], -1)
	// the following are logical but I don't like them in practice
	// fontBuilder.SetKerningPair(runeToUID['z'], runeToUID['\''], -1)
	// fontBuilder.SetKerningPair(runeToUID['z'], runeToUID['"'], -1)
	// fontBuilder.SetKerningPair(runeToUID['\''], runeToUID['s'], -1)
	// fontBuilder.SetKerningPair(runeToUID['"'], runeToUID['s'], -1)
	for _, codePoint := range "-–‑—\uE001~)]}\\&·'\"aàáâäbcçdeèéêëfghiìíïîjklmnñoòóôöpqtuùúûüvwxy" { // make 'a' attach to other letters
		fontBuilder.SetKerningPair(runeToUID['a'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['à'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['á'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['â'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['ä'], runeToUID[codePoint], -1)
	}
	for _, codePoint := range "'\"-–‑—\uE001_·~&\\/AÀÁÄÂBCÇDEÈÉËÊFGHIÌÍÏÎJKLMNÑOÒÓÖÔPQRSTUÙÚÜÛVWXYZbcçdeèéêëfhiìíîïklmnñoòóôöprstuùúûüvwxy" { // make 'j' closer to other letters
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['j'], -1)
	}
	for _, codePoint := range "aàáâä" { // make 'j' closer to other letters
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['j'], -2)
	}
	for _, codePoint := range "-–‑—\uE001_·~\\Láâäbcdeèghjkmnoòpqrstuùvwxyz" { // make ìîïÌÎÏ closer to other letters
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ì'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Ì'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Î'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Ï'], -1)
	}
	for _, codePoint := range "aà" { // make 'ìîïÌÎÏ' closer to other letters
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ì'], -2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['ï'], -2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['î'], -2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Ì'], -2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Î'], -2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Ï'], -2)
	}
	for _, codePoint := range "ñéëêóöôiïîúüû" { // extra coverage for ÌÎÏ
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Ì'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Î'], -1)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Ï'], -1)
	}
	for _, codePoint := range "aàáäâ" { // extra coverage for ÌÎÏ
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Ì'], -2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Î'], -2)
		fontBuilder.SetKerningPair(runeToUID[codePoint], runeToUID['Ï'], -2)
	}

	for _, codePoint := range ".,;:.-–‑\uE001—_~·aácçdeéfgmnoópqrsuúvwxyz" { // adjust right kerning for íïîÍÎÏ
		fontBuilder.SetKerningPair(runeToUID['í'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['î'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['Í'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['Ï'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['Î'], runeToUID[codePoint], -1)
	}
	fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID['f'], 0)
	fontBuilder.SetKerningPair(runeToUID['î'], runeToUID['f'], 0)
	for _, codePoint := range ")]}" { // further adjust right kerning for íïî
		fontBuilder.SetKerningPair(runeToUID['í'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['ï'], runeToUID[codePoint], -1)
		fontBuilder.SetKerningPair(runeToUID['î'], runeToUID[codePoint], -1)
	}
	for _, codePoint := range "AÁBCÇDEÉFGHIJKLMNOÓPQRSTUÚVWXYZâäêëóôöûübfhiíjklñt" { // further adjust right kerning for Í
		fontBuilder.SetKerningPair(runeToUID['Í'], runeToUID[codePoint], -1)
	}
	fontBuilder.SetKerningPair(runeToUID['Í'], runeToUID['j'], -2)
	//fontBuilder.SetKerningPair(runeToUID['r'], runeToUID['s'], -1) // this makes sense but doesn't look good in practice

	// add rewrite rules
	fontBuilder.AddSimpleUtf8RewriteRule('❤', '<', '3')
	fontBuilder.AddSimpleUtf8RewriteRule('💔', '<', '/', '3')
	
	setUID, err := fontBuilder.CreateGlyphSet()
	if err != nil { panic(err) }
	fontBuilder.AddGlyphSetRange(setUID, runeToUID['a'], runeToUID['z'])
	for _, codePoint := range "àáäâèéêëìíïîòóöôùúûü" {
		err := fontBuilder.AddGlyphSetListGlyph(setUID, runeToUID[codePoint])
		if err != nil { panic(err) }
	}
	fontBuilder.AddGlyphRewriteRule(1, 1, 1, []uint64{setUID, runeToUID['-'], setUID}, runeToUID['\uE001'])

	// show size
	fmt.Printf("...building font\n")
	font, err := fontBuilder.Build()
	if err != nil { panic(err) }
	err = font.Validate(ggfnt.FmtDefault)
	if err != nil { panic(err) }
	fmt.Printf("...raw size of %d bytes\n", font.RawSize())

	// export
	const FileName = "jammy-5d2-v0p3.ggfnt"
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
		if codePoint == '0' {
			err := mapZero(fontBuilder, uid, codePoint)
			if err != nil { panic(err) }
		} else if codePoint >= '1' && codePoint <= '9' {
			err := mapNum(fontBuilder, uid, codePoint)
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
		bitmap, found := pkgBitmaps[rune(codePoint)]
		if !found { panic("missing bitmap for '" + string(codePoint) + "'") }
		uid, err := fontBuilder.AddGlyph(bitmap)
		if err != nil { panic(err) }
		if codePoint == '0' {
			err := mapZero(fontBuilder, uid, codePoint)
			if err != nil { panic(err) }
		} else if codePoint >= '1' && codePoint <= '9' {
			err := mapNum(fontBuilder, uid, codePoint)
			if err != nil { panic(err) }
		} else {
			err = fontBuilder.Map(codePoint, uid)
			if err != nil { panic(err) }
		}
		codePointsMap[codePoint] = uid
	}
}

func mapZero(fontBuilder *builder.Font, clearMarkedUID uint64, codePoint rune) error {
	clearZeroUID, err := fontBuilder.AddGlyph(altZeros[0])
	if err != nil { return err }
	compactZeroUID, err := fontBuilder.AddGlyph(altZeros[1])
	if err != nil { return err }
	return fontBuilder.MapWithSwitchSingles(codePoint, SwitchZeroKey,
		clearMarkedUID, clearMarkedUID, compactZeroUID, // cases with disambiguation mark
		clearZeroUID  , clearZeroUID  , compactZeroUID, // cases without disambiguation mark
	)
}

func mapNum(fontBuilder *builder.Font, uid uint64, num rune) error {
	numIndex := uint8(num - '1')
	neutral, err := fontBuilder.AddGlyph(altNums[numIndex*2 + 0])
	if err != nil { panic(err) }
	compact, err := fontBuilder.AddGlyph(altNums[numIndex*2 + 1])
	if err != nil { panic(err) }
	return fontBuilder.MapWithSwitchSingles(num, SwitchNumStyleKey, uid, neutral, compact)	
}

// helper for mask creation
func rawAlphaMaskToWhiteMask(width int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0, -height + 2, width, 2))
	for i := 0; i < len(mask); i++ {
		img.Pix[i] = 255*mask[i]
	}
	return img
}

const KeyO   = '\x01'
const KeyI   = '\x02'
const KeyJ   = '\x03'
const KeyP   = '\x04'
const KeyL   = '\x05'
const KeyA   = '\x06'
const KeyD   = '\x07'
const KeyTAB = '\x09'

const KeyMsgI = '\x10'

const GpBtBottom  = '\uE026'
const GpBtTop     = '\uE027'
const GpBtRight   = '\uE028'
const GpBtLeft    = '\uE029'
const GpShoulderL = '\uE02A'
const GpShoulderR = '\uE02B'
const GpShoulders = '\uE02C'
const GpTriggL    = '\uE02D'
const GpTriggR    = '\uE02E'
const GpTriggers  = '\uE02F'

const Padder = '\u200A' // aka hair space
const HalfSpace = '\u2009' // aka thin space

var notdef = rawAlphaMaskToWhiteMask(3, []byte{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		1, 1, 1,
		0, 0, 0,
})

var altZeros = []*image.Alpha{
	// clearZero (can be confused with O)
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1, 
		1, 0, 0, 1, 
		1, 0, 0, 1, 
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compactZero (can be confused with O)
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline 
		0, 0, 0,
		0, 0, 0,
	}),
}

var altNums = []*image.Alpha{
	// neutral one
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact one
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// neutral two
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact two
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// neutral three
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 1, 1, 1,
		0, 0, 0, 1,
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact three
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1, // baseline 
		0, 0, 0,
		0, 0, 0,
	}),
	// neutral four
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact four
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1, // baseline 
		0, 0, 0,
		0, 0, 0,
	}),
	// neutral five (can be confused with S)
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact five (can be confused with S)
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// neutral six
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact six
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// neutral seven
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact seven
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// neutral eight
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact eight
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// neutral nine
	rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	// compact nine
	rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
}

var pkgBitmaps = map[rune]*image.Alpha{
	// --- special hacks ----
	GpBtBottom: rawAlphaMaskToWhiteMask(9, []byte{
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		1, 0, 1, 0, 0, 0, 1, 0, 1,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0, // baseline
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpBtRight: rawAlphaMaskToWhiteMask(9, []byte{
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		1, 0, 1, 0, 0, 0, 1, 1, 1,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0, // baseline
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpBtLeft: rawAlphaMaskToWhiteMask(9, []byte{
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		1, 1, 1, 0, 0, 0, 1, 0, 1,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0, // baseline
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpBtTop: rawAlphaMaskToWhiteMask(9, []byte{
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		1, 0, 1, 0, 0, 0, 1, 0, 1,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0, // baseline
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpShoulderL: rawAlphaMaskToWhiteMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1, // baseline
		0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpShoulderR: rawAlphaMaskToWhiteMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, // baseline
		1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpShoulders: rawAlphaMaskToWhiteMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, // baseline
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpTriggL: rawAlphaMaskToWhiteMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpTriggR: rawAlphaMaskToWhiteMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpTriggers: rawAlphaMaskToWhiteMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	KeyTAB: rawAlphaMaskToWhiteMask(13, []byte{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1,
		1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1,
		1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, // baseline
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	KeyO: rawAlphaMaskToWhiteMask(6, []byte{
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 1, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 1, 0, 0, 1, 1, // baseline
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
	}),
	KeyI: rawAlphaMaskToWhiteMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		1, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	KeyP: rawAlphaMaskToWhiteMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 1, 1,
		1, 0, 1, 1, 1, // baseline
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	KeyL: rawAlphaMaskToWhiteMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	KeyJ: rawAlphaMaskToWhiteMask(6, []byte{
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 1,
		1, 1, 1, 0, 1, 1,
		1, 1, 1, 0, 1, 1,
		1, 0, 1, 0, 1, 1,
		1, 1, 0, 1, 1, 1, // baseline
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
	}),
	KeyD: rawAlphaMaskToWhiteMask(6, []byte{
		1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 1, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 0, 0, 1, 1, // baseline
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
	}),
	KeyA: rawAlphaMaskToWhiteMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1, // baseline
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	KeyMsgI: rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),

	// --- main alphabet ---
	'A': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'B': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		1, 0, 1,
		1, 1, 0,
		1, 0, 1,
		1, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'C': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'D': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'E': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'F': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'G': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'H': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'I': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
		0,
	}),
	'J': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// 'J': rawAlphaMaskToWhiteMask(3, []byte{
	// 	0, 0, 0,
	// 	0, 0, 1,
	// 	0, 0, 1,
	// 	0, 0, 1,
	// 	1, 0, 1,
	// 	0, 1, 0, // baseline
	// 	0, 0, 0,
	// 	0, 0, 0,
	// }),
	'K': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 0,
		1, 1, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// 'K': rawAlphaMaskToWhiteMask(4, []byte{
	// 	0, 0, 0, 0,
	// 	1, 0, 0, 1,
	// 	1, 0, 1, 0,
	// 	1, 1, 0, 0,
	// 	1, 0, 1, 0,
	// 	1, 0, 0, 1, // baseline
	// 	0, 0, 0, 0,
	// 	0, 0, 0, 0,
	// }),
	'L': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'M': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'N': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 0, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'O': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'P': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Q': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 1, 0,
		0, 0, 0, 0,
	}),
	'R': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 0,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'S': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'T': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'U': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'V': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// 'V': rawAlphaMaskToWhiteMask(5, []byte{
	// 	0, 0, 0, 0, 0,
	// 	1, 0, 0, 0, 1,
	// 	1, 0, 0, 0, 1,
	// 	0, 1, 0, 1, 0,
	// 	0, 1, 0, 1, 0,
	// 	0, 0, 1, 0, 0, // baseline
	// 	0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0,
	// }),
	'W': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'X': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// 'X': rawAlphaMaskToWhiteMask(5, []byte{
	// 	0, 0, 0, 0, 0,
	// 	1, 0, 0, 0, 1,
	// 	0, 1, 0, 1, 0,
	// 	0, 0, 1, 0, 0,
	// 	0, 1, 0, 1, 0,
	// 	1, 0, 0, 0, 1, // baseline
	// 	0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0,
	// }),
	'Y': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Z': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 0, 1,
		0, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),

	// ---- numbers ----
	'0': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 1, 1, 
		1, 0, 0, 1, 
		1, 1, 0, 1, 
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'1': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'2': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 1,
		0, 1, 1, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'3': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 1, 1, 1,
		0, 0, 0, 1,
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'4': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'5': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'6': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'7': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'8': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'9': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),

	// --- lowercase ---
	'a': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'b': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'c': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'd': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'e': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'f': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		1, 1,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
	}),
	'g': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 1,
		1, 1, 1,
	}),
	'h': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'i': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		1,
		0,
		1,
		1,
		1, // baseline
		0,
		0,
	}),
	'j': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 1,
		0, 0,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 1,
		1, 0,
	}),
	'k': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 1,
		1, 1, 0,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'l': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
		0,
	}),
	'm': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'n': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'o': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'p': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		1, 0, 0,
		1, 0, 0,
	}),
	'q': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 1,
		0, 0, 1,
	}),
	'r': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
	}),
	's': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 1, 1,
		0, 1, 0,
		1, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	't': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		1, 0,
		1, 1,
		1, 0,
		1, 0,
		1, 1, // baseline
		0, 0,
		0, 0,
	}),
	'u': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'v': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'w': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'x': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'y': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 1,
		1, 1, 1,
	}),
	'z': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 0,
		0, 1, 0,
		0, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),

	// ---- symbols and punctuation ----
	// Note: space is special and only shifts the
	//       position 3 pixels forwards.
	' ': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'.': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		0,
		0,
	}),
	',': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		1,
		0,
	}),
	':': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		1, // baseline
		0,
		0,
	}),
	';': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		1, // baseline
		1,
		0,
	}),
	'!': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		1,
		1,
		1,
		0,
		1, // baseline
		0,
		0,
	}),
	'?': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'\'': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		1,
		1,
		0,
		0,
		0, // baseline
		0,
		0,
	}),
	'(': rawAlphaMaskToWhiteMask(2, []byte{
		0, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		0, 1,
		0, 0,
	}),
	')': rawAlphaMaskToWhiteMask(2, []byte{
		1, 0,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		1, 0,
		0, 0,
	}),
	'[': rawAlphaMaskToWhiteMask(2, []byte{
		1, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		1, 1,
		0, 0,
	}),
	']': rawAlphaMaskToWhiteMask(2, []byte{
		1, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		1, 1,
		0, 0,
	}),
	'{': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 1,
		0, 0, 0,
	}),
	'}': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0, // baseline
		1, 0, 0,
		0, 0, 0,
	}),
	'"': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'_': rawAlphaMaskToWhiteMask(3, []byte{ // NOTE: could go one pixel lower
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'-': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'\uE001': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'+': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'/': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'|': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
		0,
	}),
	'#': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1, 
		0, 1, 0, 1, 0, 
		1, 1, 1, 1, 1, 
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'~': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		1, 0, 1, 0, 1,
		0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'$': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 1, 0, 0,
		1, 1, 1, 1, 1,
		1, 0, 1, 0, 0,
		1, 1, 1, 1, 1,
		0, 0, 1, 0, 1,
		1, 1, 1, 1, 1, // baseline
		0, 0, 1, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'%': rawAlphaMaskToWhiteMask(9, []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0, 0, 0,
		1, 0, 1, 0, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 0, 1, 0, 1,
		0, 0, 0, 1, 0, 0, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	'&': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		1, 0, 1, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'*': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'<': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 0,
		0, 1,
		1, 0,
		0, 1,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'>': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 0,
		1, 0,
		0, 1,
		1, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'=': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'@': rawAlphaMaskToWhiteMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0,
		0, 1, 0, 0, 1, 0,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 1, 1,
		0, 1, 0, 0, 0, 0, // baseline
		0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'\\': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'^': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'´': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'`': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		1, 0,
		0, 1,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'¨': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'¦': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		1,
		1,
		0,
		1,
		1, // baseline
		0,
		0,
	}),
	'·': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		0, // baseline
		0,
		0,
	}),
	'¡': rawAlphaMaskToWhiteMask(1, []byte{
		0,
		0,
		1,
		0,
		1,
		1, // baseline
		1,
		0,
	}),
	'¿': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0, // baseline
		1, 0, 0,
		0, 1, 1,
		0, 0, 0,
	}),

	// --- maths ---
	'º': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'−': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'×': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'÷': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'±': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'¬': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 1,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),

	// --- accents and diacritics ---
	'À': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Á': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Â': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ä': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'È': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'É': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ê': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ë': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ì': rawAlphaMaskToWhiteMask(2, []byte{
		1, 0,
		0, 1,
		0, 0,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 0,
		0, 0,
	}),
	'Í': rawAlphaMaskToWhiteMask(2, []byte{
		0, 1,
		1, 0,
		0, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
	}),
	'Î': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ï': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ò': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ó': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ô': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ö': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ù': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ú': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Û': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ü': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),

	'à': rawAlphaMaskToWhiteMask(4, []byte{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'á': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'â': rawAlphaMaskToWhiteMask(4, []byte{
		0, 1, 0, 0,
		1, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ä': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'è': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'é': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ê': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ë': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ì': rawAlphaMaskToWhiteMask(2, []byte{
		1, 0,
		0, 1,
		0, 0,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 0,
		0, 0,
	}),
	'í': rawAlphaMaskToWhiteMask(2, []byte{
		0, 1,
		1, 0,
		0, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
	}),
	'î': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ï': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ò': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ó': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ô': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ö': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ù': rawAlphaMaskToWhiteMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ú': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'û': rawAlphaMaskToWhiteMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ü': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),

	// --- additional letters for spanish ---
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
	}),
	'ñ': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ç': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 1, 0,
		1, 0, 0,
	}),
	'ç': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 1, 0,
		1, 0, 0,
	}),

	// --- currencies ---
	'€': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 1, 0, 0,
		1, 1, 1, 0,
		0, 1, 0, 0,
		0, 0, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'£': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 1, 0,
		0, 1, 0, 0, 1,
		1, 1, 1, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'¢': rawAlphaMaskToWhiteMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 1, 1,
		1, 0, 1, 0,
		0, 1, 1, 1,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'¥': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'¤': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
		0, 1, 0, 1, 0,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),

	// --- notes ---
	'♩': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		1, 1, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'♪': rawAlphaMaskToWhiteMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 1, 0,
		0, 0, 1, 0, 1,
		0, 0, 1, 0, 0,
		1, 1, 1, 0, 0,
		1, 1, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'♫': rawAlphaMaskToWhiteMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1,
		0, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 0, 0, 0, 1,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1, // baseline
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),

	// --- dashes ---
	'–': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'‑': rawAlphaMaskToWhiteMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'—': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),

	// --- additional symbols ---
	'◀': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 1, 1,
		1, 1, 1,
		0, 1, 1,
		0, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'▶': rawAlphaMaskToWhiteMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 1, 0,
		1, 1, 1,
		1, 1, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'❤': rawAlphaMaskToWhiteMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 0, 0, // baseline
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'💔': rawAlphaMaskToWhiteMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1,
		0, 1, 0, 1, 1, 1, 0,
		0, 0, 0, 1, 1, 0, 0, // baseline
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	' ': rawAlphaMaskToWhiteMask(2, []byte{ // aka thin space \u2009
		0, 0, 
		0, 0, 
		0, 0, 
		0, 0, 
		0, 0, 
		0, 0, // baseline
		0, 0, 
		0, 0, 
	}),
	' ': rawAlphaMaskToWhiteMask(1, []byte{ // aka hair space \u200A
		0,
		0,
		0,
		0,
		0,
		0, // baseline
		0,
		0,
	}),
}
