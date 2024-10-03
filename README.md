# ggfnt-fonts

This project exposes free [ggfnt](https://github.com/tinne26/ggfnt) fonts for use with Golang programs, most commonly [Ebitengine](https://github.com/hajimehoshi/ebiten) games that use the [ptxt](https://github.com/tinne26/ptxt) font rendering library.

## Font samples

`github.com/tinne26/ggfnt-fonts/tinny`
![sample_tinny](https://github.com/user-attachments/assets/1e3326d6-6831-4c1b-8116-0d0a858ab561)

`github.com/tinne26/ggfnt-fonts/tinnybold`
![sample_tinny_bold](https://github.com/user-attachments/assets/636ab70b-98f5-4b5d-ab6c-d035ec428026)

`github.com/tinne26/ggfnt-fonts/graybit`
![sample_graybit](https://github.com/user-attachments/assets/a9bd799b-9ded-43f2-8081-0d788940339f)

`github.com/tinne26/ggfnt-fonts/omen`
![sample_omen](https://github.com/user-attachments/assets/19e1739d-398d-41f4-9bea-4650decb1399)

`github.com/tinne26/ggfnt-fonts/candy`
![sample_candy](https://github.com/user-attachments/assets/25fbd09e-0030-40cf-ad38-02c53056224f)

`github.com/tinne26/ggfnt-fonts/strut`
![sample_strut](https://github.com/user-attachments/assets/8d195df0-09fb-4cd8-a5a2-1859b8745bd4)

`github.com/tinne26/ggfnt-fonts/starship`
![sample_starship](https://github.com/user-attachments/assets/46b62db4-2b61-40b7-abf7-02a1cc5107f0)

`github.com/tinne26/ggfnt-fonts/flick`
![sample_flick](https://github.com/user-attachments/assets/73b5ec93-4e73-47cc-be10-9a552af02050)

`github.com/tinne26/ggfnt-fonts/minitile`
![sample_minitile](https://github.com/user-attachments/assets/6fca2315-0efe-4d8d-9140-a5e13eab6a4e)

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
