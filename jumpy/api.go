package jumpy

import "io"
import "time"
import "math"
import _ "embed"
import "github.com/tinne26/ggfnt"

//go:embed jumpy-6d0-v0p1.ggfnt
var bytes []byte

var cachedFont *ggfnt.Font

const NotdefRune = '\uE000'
const Notdef = ggfnt.GlyphIndex(0)

func Release() { cachedFont = nil }
func Font() *ggfnt.Font {
	if cachedFont == nil {
		font, err := ggfnt.Parse(&byteSliceReader{ data: bytes })
		if err != nil { panic(err) } // (go test .)
		cachedFont = font
	}
	return cachedFont
}

// --- glyph picker ---

// mimics ptxt/strand.GlyphPickerPass
type glyphPickerPass = uint8
const (
	measurePass glyphPickerPass = iota
	drawPass
	bufferPass
)

// Can be configured and used with ptxt:
//   var picker jumpy.GoldenPicker
//   strand.GlyphPickers().Add(&picker)
// This glyph picker is a basic implementation using an alternating
// sequence for glyphs based on time and the golden ratio. In general,
// using ticks is more appropriate, but this code doesn't have access
// to that information.
type GoldenPicker struct {
	initialized bool
	nextFlipTime time.Time
	stateDuration time.Duration
	flipped bool
	sequenceValue float64 // 1.61803399
}

// The default duration is 400*time.Millisecond.
func (self *GoldenPicker) SetStateDuration(duration time.Duration) {
	self.stateDuration = duration
}

// Implements ptxt/strand.GlyphPicker.
func (self *GoldenPicker) Pick(codePoint rune, groupSize uint8, flags ggfnt.AnimationFlags, numQueuedGlyphs int) uint8 {
	if groupSize == 1 { return 0 }
	self.sequenceValue += 1.61803399
	if self.sequenceValue >= 1.0 { self.sequenceValue -= 1.0 }
	if self.sequenceValue >= 1.0 { self.sequenceValue -= 1.0 }
	pick := uint8(self.sequenceValue*float64(groupSize))
	if self.flipped {
		pick = (pick + groupSize/2)
		if pick >= groupSize {
			pick -= groupSize
		}
	}
	return pick
}

// Implements ptxt/strand.GlyphPicker.
func (self *GoldenPicker) NotifyAddedGlyph(ggfnt.GlyphIndex, rune, uint8, ggfnt.AnimationFlags) {}

// Implements ptxt/strand.GlyphPicker.
func (self *GoldenPicker) NotifyPass(pass glyphPickerPass, start bool) {
	if !self.initialized { self.initialize() }
	if start {
		now := time.Now()
		if now.After(self.nextFlipTime) {
			self.nextFlipTime = now.Add(self.stateDuration)
			self.flipped = !self.flipped
		}

		self.sequenceValue = 0.0
	}
}

// Implements ptxt/strand.GlyphPicker.
func (self *GoldenPicker) initialize() {
	self.initialized = true
	if self.stateDuration == 0 {
		self.stateDuration = 400*time.Millisecond
	}
	self.nextFlipTime = time.Now().Add(self.stateDuration)
}

// Can be configured and used with ptxt:
//   var picker jumpy.PulsePicker
//   strand.GlyphPickers().Add(&picker)
// This glyph picker chooses glyphs based on a pulse function, which has
// a certain low, transition and high widths, and some speed. In general
// using ticks is more appropriate than using time, but this code doesn't
// have access to that information.
type PulsePicker struct {
	initialized bool
	lastStart time.Time
	pulseHighWidth float64
	pulseTransitionWidth float64
	pulseLowWidth float64
	cycleLen float64
	glyphsPerSecond float64
	startPoint float64
	operationPoint float64
}

// The defaults are 35, 1, 2. The widths are given in glyphs.
func (self *PulsePicker) SetPulseWidth(low, transition, high float64) {
	self.pulseLowWidth = low
	self.pulseTransitionWidth = transition
	self.pulseHighWidth = high
	self.cycleLen = low + transition + high + transition
}

// The default is 20.
func (self *PulsePicker) SetPulseSpeed(glyphsPerSecond float64) {
	self.glyphsPerSecond = glyphsPerSecond
}

// Implements ptxt/strand.GlyphPicker.
func (self *PulsePicker) Pick(codePoint rune, groupSize uint8, flags ggfnt.AnimationFlags, numQueuedGlyphs int) uint8 {
	// update operation angle
	self.operationPoint -= 1.0
	for self.operationPoint < 0 {
		self.operationPoint += self.cycleLen
	}

	// trivial case
	if groupSize == 1 { return 0 }
	
	// general case
	var lerp = func(a, b, t float64) float64 {
		return a + (b - a)*t
	}
	value := self.operationPoint
	if value < self.pulseLowWidth { return 0 }
	value -= self.pulseLowWidth
	if value < self.pulseTransitionWidth {
		t := value/self.pulseTransitionWidth
		return uint8(lerp(0.0, float64(groupSize) - 0.00001, t))
	}
	value -= self.pulseTransitionWidth
	if value < self.pulseHighWidth { return groupSize - 1 }
	value -= self.pulseHighWidth
	t := value/self.pulseTransitionWidth
	return uint8(lerp(float64(groupSize) - 0.00001, 0.0, t))
}

// Implements ptxt/strand.GlyphPicker.
func (self *PulsePicker) NotifyAddedGlyph(ggfnt.GlyphIndex, rune, uint8, ggfnt.AnimationFlags) {}

// Implements ptxt/strand.GlyphPicker.
func (self *PulsePicker) NotifyPass(pass glyphPickerPass, start bool) {
	if !self.initialized { self.initialize() }
	if start {
		now := time.Now()
		diff := now.Sub(self.lastStart)
		self.lastStart = now
		self.startPoint += diff.Seconds()*self.glyphsPerSecond
		self.startPoint  = math.Mod(self.startPoint, self.cycleLen)
		self.operationPoint = self.startPoint
	}
}

func (self *PulsePicker) initialize() {
	self.initialized = true
	if self.cycleLen == 0 { self.SetPulseWidth(35, 1, 2) }
	if self.glyphsPerSecond == 0 { self.glyphsPerSecond = 20.0 }
	self.startPoint = self.pulseLowWidth
	self.lastStart = time.Now()
}

// --- helpers ---

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
