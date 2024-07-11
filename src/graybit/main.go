package main

import "os"
import "fmt"
import "image"

import "github.com/tinne26/ggfnt"
import "github.com/tinne26/ggfnt/builder"

// TODO: rename. graybit, graylite, edge, graphite, replica, naas, grayness. grayed, ash, carbon

// TODO:
// - maybe animated cursor? either as vertical bar or underscore
// switches:
// - add named gamepad keys and so on? hmmm...

// globals
var SwitchZeroKey, SwitchNumStyleKey uint8

func main() {
	// create font builder
	fmt.Print("creating new font builder\n")
	fontBuilder := builder.New()

	// add metadata
	fmt.Print("...adding metadata\n")
	err := fontBuilder.SetName("graybit")
	if err != nil { panic(err) }
	err = fontBuilder.SetFamily("graybit")
	if err != nil { panic(err) }
	err = fontBuilder.SetAuthor("tinne")
	if err != nil { panic(err) }
	err = fontBuilder.SetAbout("Born from tinne's entries for Ebitengine game jams and expanded through the years, this was the first font to be ever encoded in the ggfnt format.")
	if err != nil { panic(err) }
	fontBuilder.SetVersion(0, 5)
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
	fontBuilder.SetLineGap(1)
	err = fontBuilder.GetMetricsStatus()
	if err != nil { panic(err) }

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
	fontBuilder.SetKerningPair(runeToUID['.'], runeToUID['?'], -1) // improve "..?"
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
	//fontBuilder.SetKerningPair(runeToUID['L'], runeToUID['Y'], -1) // this is not bad, but in case of ambivalence best avoid touching
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
	const FileName = "graybit-5d2-v0p5.ggfnt"
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
func rawMask(width int, mask []byte) *image.Alpha {
	height := len(mask)/width
	img := image.NewAlpha(image.Rect(0, -height + 2, width, 2))
	copy(img.Pix, mask)
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

var notdef = rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	rawMask(4, []byte{
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
	rawMask(3, []byte{
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
	GpBtBottom: rawMask(9, []byte{
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		1, 0, 1, 0, 0, 0, 1, 0, 1,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0, // baseline
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpBtRight: rawMask(9, []byte{
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		1, 0, 1, 0, 0, 0, 1, 1, 1,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0, // baseline
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpBtLeft: rawMask(9, []byte{
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		1, 1, 1, 0, 0, 0, 1, 0, 1,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0, // baseline
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpBtTop: rawMask(9, []byte{
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		1, 0, 1, 0, 0, 0, 1, 0, 1,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 1, 0, 0, 0, // baseline
		0, 0, 0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpShoulderL: rawMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1, // baseline
		0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpShoulderR: rawMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, // baseline
		1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpShoulders: rawMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, // baseline
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpTriggL: rawMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpTriggR: rawMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	GpTriggers: rawMask(13, []byte{
		0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	KeyTAB: rawMask(13, []byte{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1,
		1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1,
		1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, // baseline
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	KeyO: rawMask(6, []byte{
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 1, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 1, 0, 0, 1, 1, // baseline
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
	}),
	KeyI: rawMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		1, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	KeyP: rawMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 1, 1,
		1, 0, 1, 1, 1, // baseline
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	KeyL: rawMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 0, 1, 1, 1,
		1, 0, 0, 0, 1, // baseline
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	KeyJ: rawMask(6, []byte{
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 1,
		1, 1, 1, 0, 1, 1,
		1, 1, 1, 0, 1, 1,
		1, 0, 1, 0, 1, 1,
		1, 1, 0, 1, 1, 1, // baseline
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
	}),
	KeyD: rawMask(6, []byte{
		1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 1, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 0, 1,
		1, 0, 0, 0, 1, 1, // baseline
		1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0,
	}),
	KeyA: rawMask(5, []byte{
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1, // baseline
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 0,
	}),
	KeyMsgI: rawMask(3, []byte{
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
	'A': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'B': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		1, 0, 1,
		1, 1, 0,
		1, 0, 1,
		1, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'C': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'D': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'E': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'F': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'G': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 0, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'H': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'I': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
		0,
	}),
	'J': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// 'J': rawMask(3, []byte{
	// 	0, 0, 0,
	// 	0, 0, 1,
	// 	0, 0, 1,
	// 	0, 0, 1,
	// 	1, 0, 1,
	// 	0, 1, 0, // baseline
	// 	0, 0, 0,
	// 	0, 0, 0,
	// }),
	'K': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 0,
		1, 1, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// 'K': rawMask(4, []byte{
	// 	0, 0, 0, 0,
	// 	1, 0, 0, 1,
	// 	1, 0, 1, 0,
	// 	1, 1, 0, 0,
	// 	1, 0, 1, 0,
	// 	1, 0, 0, 1, // baseline
	// 	0, 0, 0, 0,
	// 	0, 0, 0, 0,
	// }),
	'L': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'M': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 1, 0, 1, 1,
		1, 1, 0, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		1, 0, 0, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'N': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 1, 0, 1,
		1, 1, 1, 1,
		1, 0, 1, 1,
		1, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'O': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'P': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Q': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 0, 1, 1,
		1, 1, 1, 1, // baseline
		0, 0, 1, 0,
		0, 0, 0, 0,
	}),
	'R': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 0,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'S': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'T': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'U': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'V': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// 'V': rawMask(5, []byte{
	// 	0, 0, 0, 0, 0,
	// 	1, 0, 0, 0, 1,
	// 	1, 0, 0, 0, 1,
	// 	0, 1, 0, 1, 0,
	// 	0, 1, 0, 1, 0,
	// 	0, 0, 1, 0, 0, // baseline
	// 	0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0,
	// }),
	'W': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 0, 1, 0,
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'X': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	// 'X': rawMask(5, []byte{
	// 	0, 0, 0, 0, 0,
	// 	1, 0, 0, 0, 1,
	// 	0, 1, 0, 1, 0,
	// 	0, 0, 1, 0, 0,
	// 	0, 1, 0, 1, 0,
	// 	1, 0, 0, 0, 1, // baseline
	// 	0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0,
	// }),
	'Y': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Z': rawMask(3, []byte{
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
	'0': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 1, 1, 
		1, 0, 0, 1, 
		1, 1, 0, 1, 
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'1': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'2': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 1,
		0, 1, 1, 0,
		1, 0, 0, 0,
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'3': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 1, 1, 1,
		0, 0, 0, 1,
		1, 1, 1, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'4': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline 
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'5': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'6': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'7': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1,
		0, 0, 0, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'8': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'9': rawMask(4, []byte{
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
	'a': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'b': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'c': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'd': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'e': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'f': rawMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		1, 1,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
	}),
	'g': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 1,
		1, 1, 1,
	}),
	'h': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'i': rawMask(1, []byte{
		0,
		1,
		0,
		1,
		1,
		1, // baseline
		0,
		0,
	}),
	'j': rawMask(2, []byte{
		0, 0,
		0, 1,
		0, 0,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 1,
		1, 0,
	}),
	'k': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 1,
		1, 1, 0,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'l': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
		0,
	}),
	'm': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'n': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'o': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'p': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		1, 0, 0,
		1, 0, 0,
	}),
	'q': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 1,
		0, 0, 1,
	}),
	'r': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
	}),
	's': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 1, 1,
		0, 1, 0,
		1, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	't': rawMask(2, []byte{
		0, 0,
		1, 0,
		1, 1,
		1, 0,
		1, 0,
		1, 1, // baseline
		0, 0,
		0, 0,
	}),
	'u': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'v': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'w': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 0, 1, 0, 1,
		1, 0, 1, 0, 1,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'x': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'y': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 1,
		1, 1, 1,
	}),
	'z': rawMask(3, []byte{
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
	' ': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'.': rawMask(1, []byte{
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		0,
		0,
	}),
	',': rawMask(1, []byte{
		0,
		0,
		0,
		0,
		0,
		1, // baseline
		1,
		0,
	}),
	':': rawMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		1, // baseline
		0,
		0,
	}),
	';': rawMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		1, // baseline
		1,
		0,
	}),
	'!': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		0,
		1, // baseline
		0,
		0,
	}),
	'?': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 0,
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'\'': rawMask(1, []byte{
		0,
		1,
		1,
		0,
		0,
		0, // baseline
		0,
		0,
	}),
	'(': rawMask(2, []byte{
		0, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		0, 1,
		0, 0,
	}),
	')': rawMask(2, []byte{
		1, 0,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		1, 0,
		0, 0,
	}),
	'[': rawMask(2, []byte{
		1, 1,
		1, 0,
		1, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		1, 1,
		0, 0,
	}),
	']': rawMask(2, []byte{
		1, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1,
		0, 1, // baseline
		1, 1,
		0, 0,
	}),
	'{': rawMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 1,
		0, 0, 0,
	}),
	'}': rawMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0, // baseline
		1, 0, 0,
		0, 0, 0,
	}),
	'"': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'_': rawMask(3, []byte{ // NOTE: could go one pixel lower
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'-': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'\uE001': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'+': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'/': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 1, 0,
		1, 0, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'|': rawMask(1, []byte{
		0,
		1,
		1,
		1,
		1,
		1, // baseline
		0,
		0,
	}),
	'#': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1, 
		0, 1, 0, 1, 0, 
		1, 1, 1, 1, 1, 
		0, 1, 0, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'~': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		1, 0, 1, 0, 1,
		0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, // baseline
		0, 0, 0, 0, 0,
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
		0, 0, 0, 0, 0,
	}),
	'%': rawMask(9, []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0, 0, 0,
		1, 0, 1, 0, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 0, 1, 0, 1,
		0, 0, 0, 1, 0, 0, 0, 1, 0, // baseline
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}),
	'&': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 1, 0, 0,
		1, 0, 1, 0,
		0, 1, 1, 0,
		1, 0, 0, 1,
		1, 1, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'*': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'<': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 1,
		1, 0,
		0, 1,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'>': rawMask(2, []byte{
		0, 0,
		0, 0,
		1, 0,
		0, 1,
		1, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'=': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'@': rawMask(6, []byte{
		0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0,
		0, 1, 0, 0, 1, 0,
		1, 0, 1, 1, 0, 1,
		1, 0, 1, 1, 1, 1,
		0, 1, 0, 0, 0, 0, // baseline
		0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0,
	}),
	'\\': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'^': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'´': rawMask(2, []byte{
		0, 0,
		0, 1,
		1, 0,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'`': rawMask(2, []byte{
		0, 0,
		1, 0,
		0, 1,
		0, 0,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'¨': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'¦': rawMask(1, []byte{
		0,
		1,
		1,
		0,
		1,
		1, // baseline
		0,
		0,
	}),
	'·': rawMask(1, []byte{
		0,
		0,
		0,
		1,
		0,
		0, // baseline
		0,
		0,
	}),
	'¡': rawMask(1, []byte{
		0,
		0,
		1,
		0,
		1,
		1, // baseline
		1,
		0,
	}),
	'¿': rawMask(3, []byte{
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
	'º': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 0, 1,
		0, 1, 0,
		0, 0, 0,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'−': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'×': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		0, 1, 0,
		1, 0, 1,
		0, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'÷': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'±': rawMask(3, []byte{
		0, 0, 0,
		0, 1, 0,
		1, 1, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'¬': rawMask(3, []byte{
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
	'À': rawMask(3, []byte{
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
	'Á': rawMask(3, []byte{
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
	'Â': rawMask(3, []byte{
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
	'Ä': rawMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'È': rawMask(3, []byte{
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
	'É': rawMask(3, []byte{
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
	'Ê': rawMask(3, []byte{
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
	'Ë': rawMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ì': rawMask(2, []byte{
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
	'Í': rawMask(2, []byte{
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
	'Î': rawMask(3, []byte{
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
	'Ï': rawMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ò': rawMask(3, []byte{
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
	'Ó': rawMask(3, []byte{
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
	'Ô': rawMask(3, []byte{
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
	'Ö': rawMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ù': rawMask(3, []byte{
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
	'Ú': rawMask(3, []byte{
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
	'Û': rawMask(3, []byte{
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
	'Ü': rawMask(3, []byte{
		1, 0, 1,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),

	'à': rawMask(4, []byte{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'á': rawMask(4, []byte{
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'â': rawMask(4, []byte{
		0, 1, 0, 0,
		1, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'ä': rawMask(4, []byte{
		0, 0, 0, 0,
		1, 0, 1, 0,
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'è': rawMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'é': rawMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ê': rawMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ë': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 1, 0,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ì': rawMask(2, []byte{
		1, 0,
		0, 1,
		0, 0,
		0, 1,
		0, 1,
		0, 1, // baseline
		0, 0,
		0, 0,
	}),
	'í': rawMask(2, []byte{
		0, 1,
		1, 0,
		0, 0,
		1, 0,
		1, 0,
		1, 0, // baseline
		0, 0,
		0, 0,
	}),
	'î': rawMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ï': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ò': rawMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ó': rawMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ô': rawMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ö': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ù': rawMask(3, []byte{
		1, 0, 0,
		0, 1, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ú': rawMask(3, []byte{
		0, 0, 1,
		0, 1, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'û': rawMask(3, []byte{
		0, 1, 0,
		1, 0, 1,
		0, 0, 0,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'ü': rawMask(3, []byte{
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
	'Ñ': rawMask(4, []byte{
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
	'ñ': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'Ç': rawMask(3, []byte{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1, // baseline
		0, 1, 0,
		1, 0, 0,
	}),
	'ç': rawMask(3, []byte{
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
	'€': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 1,
		0, 1, 0, 0,
		1, 1, 1, 0,
		0, 1, 0, 0,
		0, 0, 1, 1, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'£': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 1, 0,
		0, 1, 0, 0, 1,
		1, 1, 1, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 1, 1, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'¢': rawMask(4, []byte{
		0, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 1, 1,
		1, 0, 1, 0,
		0, 1, 1, 1,
		0, 0, 1, 0, // baseline
		0, 0, 0, 0,
		0, 0, 0, 0,
	}),
	'¥': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 1,
		0, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'¤': rawMask(5, []byte{
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
	'♩': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		1, 1, 1,
		1, 1, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'♪': rawMask(5, []byte{
		0, 0, 0, 0, 0,
		0, 0, 1, 1, 0,
		0, 0, 1, 0, 1,
		0, 0, 1, 0, 0,
		1, 1, 1, 0, 0,
		1, 1, 1, 0, 0, // baseline
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}),
	'♫': rawMask(7, []byte{
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
	'–': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'‑': rawMask(2, []byte{
		0, 0,
		0, 0,
		0, 0,
		1, 1,
		0, 0,
		0, 0, // baseline
		0, 0,
		0, 0,
	}),
	'—': rawMask(3, []byte{
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
	'◀': rawMask(3, []byte{
		0, 0, 0,
		0, 0, 1,
		0, 1, 1,
		1, 1, 1,
		0, 1, 1,
		0, 0, 1, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'▶': rawMask(3, []byte{
		0, 0, 0,
		1, 0, 0,
		1, 1, 0,
		1, 1, 1,
		1, 1, 0,
		1, 0, 0, // baseline
		0, 0, 0,
		0, 0, 0,
	}),
	'❤': rawMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 0, 0, // baseline
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	'💔': rawMask(7, []byte{
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		1, 1, 1, 0, 1, 1, 1,
		1, 1, 1, 0, 1, 1, 1,
		0, 1, 0, 1, 1, 1, 0,
		0, 0, 0, 1, 1, 0, 0, // baseline
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}),
	' ': rawMask(2, []byte{ // aka thin space \u2009
		0, 0, 
		0, 0, 
		0, 0, 
		0, 0, 
		0, 0, 
		0, 0, // baseline
		0, 0, 
		0, 0, 
	}),
	' ': rawMask(1, []byte{ // aka hair space \u200A
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
