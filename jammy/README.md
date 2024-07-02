# jammy
[![Go Reference](https://pkg.go.dev/badge/tinne26/ggfnt-fonts/jammy.svg)](https://pkg.go.dev/github.com/tinne26/ggfnt-fonts/jammy)

A really small font with a fairly neutral character, angular and slightly modern. This font was the first bitmap font I created for some parts of the UI in [bindless](https://github.com/tinne26/bindless), my entry for the first Ebitengine game jam in 2022. The font was greatly expanded during the next year's jam, and it eventually became the first font to be encoded in the ggfnt format (before that I was using it directly as images).

Right before making all the ptxt and ggfnt projects public, I expanded this font even more to cover the full ASCII range and a few more characters.

The import path is:
```Golang
import "github.com/tinne26/ggfnt-fonts/jammy"
```

## Glyphs

Common glyphs:
- Notdef.
- Full ASCII range (notice: this font was originally designed for uppercase; lowercase letters are quite cramped).
- Basic diacritics: `àáäâ` (for `aeiou`, `AEIOU` and the diacritics in isolation).
- Other common punctuation symbols `¡ ¿`, `– ‑ —`.
- Spanish letters `ÑñÇç`.
- Common currency symbols `€ £ ¢ ¥ ¤`.
- A few extra math symbols `− × ÷ ± º`.
- A few extra symbols `♩ ♪ ♫ ◀ ▶ ❤ 💔`.
- Thin space ` ` and hair space ` `.

Named glyphs:
- `"notdef"`, also mapped to `\uE000` for accessibility.
- TODO: expose all the special characters for buttons? Probably yes.

Rewrite rules:
- `<3` will translate to ❤, and `</3` to 💔.
- Using hyphens (-) between lowercase characters will use a hyphen that's one pixel lower than usual. This is something weird that I did while I was testing rewrite rules involving glyph sets.

Other private use characters:
- `\uE001`: low hyphen.

## License

The code in this folder is MIT licensed, the font itself is [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
