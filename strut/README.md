# unnamed
[![Go Reference](https://pkg.go.dev/badge/tinne26/ggfnt-fonts/strut.svg)](https://pkg.go.dev/github.com/tinne26/ggfnt-fonts/strut)

A swaggy font with some afrobeat / funky / 80s energy to it, but still quite angular and rigid despite some of the more playful features. The font covers the full ASCII range and a few more common glyphs.

> [!NOTE]
> If you want to use lowercase + accents and diacritics on uppercase letters, you will most likely want to increase the line gap by 2 pixels to avoid vertical overlaps (e.g. q + Á). For example, on `ptxt` you can adjust this with the line interspacing shift.

The import path is:
```Golang
import "github.com/tinne26/ggfnt-fonts/strut"
```

## Glyphs

Common glyphs:
- Notdef.
- Complete ASCII range.
- Basic diacritics: `ÀÁÄÂ` (for both `AEIOU` and `aeiou`).
- Other common punctuation symbols `¡ ¿`, `– ‑ —`, `‘’ “” …`.
- Spanish letters `ÑñÇç`.
- Some currency symbols `€ ¢ ¥ ¤`.

Named glyphs:
- `"notdef"`, also mapped to `'\uE000'` for accessibility.

## License

The code in this folder is MIT licensed, the font itself is [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
