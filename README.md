# ggfnt-fonts

This project exposes free [ggfnt](https://github.com/tinne26/ggfnt) fonts for use with Golang programs, most commonly [Ebitengine](https://github.com/hajimehoshi/ebiten) games that use the [ptxt](https://github.com/tinne26/ptxt) font rendering library.

## Font samples

`github.com/tinne26/ggfnt-fonts/tinny`
![sample_tinny](https://github.com/tinne26/ggfnt-fonts/assets/95440833/a87208ad-e963-470d-91d2-a8a8348ef10b)

`github.com/tinne26/ggfnt-fonts/graybit`
![sample_graybit](https://github.com/tinne26/ggfnt-fonts/assets/95440833/66f07930-bde0-4e23-8adf-82e2a60487f7)

`github.com/tinne26/ggfnt-fonts/omen`
![sample_omen](https://github.com/tinne26/ggfnt-fonts/assets/95440833/b08fdc74-cc46-4d4a-8a8c-4c9f44a999d1)

`github.com/tinne26/ggfnt-fonts/candy`
![sample_candy](https://github.com/user-attachments/assets/00f377b1-985e-4061-b367-32bd055e5a03)

`github.com/tinne26/ggfnt-fonts/strut`


`github.com/tinne26/ggfnt-fonts/starship`
![sample_starship](https://github.com/tinne26/ggfnt-fonts/assets/95440833/c4309f82-1a38-475a-932b-defdecaf396e)

`github.com/tinne26/ggfnt-fonts/flick`
![sample_flick](https://github.com/tinne26/ggfnt-fonts/assets/95440833/244ae966-8191-4544-933b-8fb6954598c0)

`github.com/tinne26/ggfnt-fonts/minitile`
![sample_minitile](https://github.com/user-attachments/assets/7c410318-c395-43c5-9bda-51c6e26e6f4f)


## Licenses

Most fonts in this project are licensed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/), though some might be adaptations of well known fonts designed in other formats and with their own licenses. See each font subfolder for additional details.

The Golang code used to provide access to the fonts is extremely trivial and MIT licensed. The code for building fonts is not licensed, as it contains representations of the font glyphs themselves (but feel free to use the structure and basic methods as reference).

## Usage

All subpackages expose at least two methods:
- `Font()`, which parses the font if it wasn't parsed yet, caches it and returns it as a `*ggfnt.Font`.
- `Release()`, which frees the cached font if it was loaded at any point by `Font()`.

Some subpackages contain additional constants for font settings, specific glyphs and code point mappings. The most common is probably `NotdefRune` (`rune`) and `Notdef` (`ggfnt.GlyphIndex`). Glyph pickers might also be provided for some fancy fonts in order to make life easier when working with ptxt.

Example program:
```Golang
package main

import "fmt"
import "github.com/tinne26/ggfnt-fonts/graybit"

func main() {
	font := graybit.Font()
	verMajor, verMinor := font.Header().VersionMajor(), font.Header().VersionMinor()
	fmt.Printf("Font: %s (v%d.%02d)\n", font.Header().Name(), verMajor, verMinor)
	fmt.Printf("Num. glyphs: %d\n", font.Glyphs().Count())
	fmt.Printf("Author: %s\n", font.Header().Author())
	fmt.Printf("About:\n> %s\n", font.Header().About())
}
```
The names of the subpackages and their paths match what you can find in this repository.

## Pangrams

I actually made some of the pangrams for the examples on my own, here they are:
- Twin axes ablaze, the grumpy viking reconquered the fjord.
- The zombie geeks acquired explosive jellyfish now!?
- Saxophonists frequently acknowledge my jazzy vibes.
- Catalyzer hijack verified, equip low oxygen bombs.
- Objectively speaking, frozen marshmallows are darn exquisite.
- Hendrix's jam —if unequivocably squawky— was hypnotizing.
- Josephine, buddy, the squeezy wolfkin are exclusively mine!

They are a tad long, but at least some of them have nice and understandable stories. A few more variations:
- Wanted zombie geeks aquire explosive jellyfish. *(shorter but less funny imo)*
- Saxophonists disqualify overembellished spacewalking jazz. *(only five words, but convoluted)*
- Hendrix's jam was hypnotizing, unequivocably funky.

Sadly, I don't know the authors of "amazingly few discotheques provide jukeboxes", "sphinx of black quartz, judge my vow" and other cool pangrams; I'd love to give them credit otherwise!
