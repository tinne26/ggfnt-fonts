# ggfnt-fonts

This project exposes free [ggfnt](https://github.com/tinne26/ggfnt) fonts for use with Golang programs —most commonly [Ebitengine](https://github.com/hajimehoshi/ebiten/v2) games that use the [ptxt](https://github.com/tinne26/ptxt) font rendering library—.

## Licenses

Most fonts in this project are licensed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/), though some might be adaptations of well known fonts designed in other formats and with their own licenses. See each font subfolder for additional details.

The Golang code used to provide access to the fonts is extremely trivial and MIT licensed. The code for building fonts is not licensed, as it contains representations of the font glyphs themselves (but feel free to use the structure and basic methods as reference).

## Usage

All subpackages expose two methods:
- `Font()`, which parses the font if it wasn't parsed yet, caches it and returns it as a `*ggfnt.Font`.
- `Release()`, which frees the cached font if it was loaded at any point by `Font()`.

Some subpackages contain additional constants for specific glyphs or code point mappings. The most common is probably `NotdefRune` (`rune`) and `Notdef` (`ggfnt.GlyphIndex`). Glyph pickers might also be provided for some fancy fonts in order to make life easier when working with ptxt.

Example program:
```Golang
package main

import "fmt"
import "github.com/tinne26/ggfnt-fonts/jammy"

func main() {
	font := jammy.Font()
	verMajor, verMinor := font.Header().VersionMajor(), font.Header().VersionMinor()
	fmt.Printf("Font: %s (v%d.%02d)\n", font.Header().Name(), verMajor, verMinor)
	fmt.Printf("Num. glyphs: %d\n", font.Glyphs().Count())
	fmt.Printf("Author: %s\n", font.Header().Author())
	fmt.Printf("About:\n> %s\n", font.Header().About())
}
```
The names of the subpackages and their paths match what you can find in this repository.

## Subpackage cheatsheet

- **jammy 5d2** | `github.com/tinne26/fonts/jammy`
- **tinny 6d3** | `github.com/tinne26/fonts/tinny`
