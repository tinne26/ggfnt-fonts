# candy
[![Go Reference](https://pkg.go.dev/badge/tinne26/ggfnt-fonts/candy.svg)](https://pkg.go.dev/github.com/tinne26/ggfnt-fonts/candy)

A comfy, rounded and happy display font. Only uppercase.

The import path is:
```Golang
import "github.com/tinne26/ggfnt-fonts/candy"
```

## Glyphs

Common glyphs:
- Notdef.
- Full ASCII range except for lowercase letters `a-z`.
- Other common punctuation symbols `¡ ¿`, `– ‑ —`, `‘’ “” • …`.
- Spanish letters `ÑÇ`.
- Hair and thin spaces (` `, ` `).

Named glyphs:
- `"notdef"`, also mapped to `'\uE000'` for accessibility.
- `"candy"`, mapped to `'\U0001F36C'` (🍬).
- `"dye-candy"`, like `"candy"` but using the monochrome main dye color instead of a palette. Also mapped to `'\uEDCA'` for accessibility.
- `"fat-dot"`, also mapped to `'\uED01'`. This is 4x4 instead of the default 2x2. It does look better, but it's not consistent with other punctuation symbols (`¿¡:;!?`).

> [!TIP]
> Named glyphs are also accessible as package constants (`candy.CandyRune`, `candy.Candy`, `candy.FatDotRune`, `candy.FatDot`, etc).

## License

The code in this folder is MIT licensed, the font itself is [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
