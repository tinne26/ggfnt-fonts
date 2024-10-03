# tinny
[![Go Reference](https://pkg.go.dev/badge/tinne26/ggfnt-fonts/tinny.svg)](https://pkg.go.dev/github.com/tinne26/ggfnt-fonts/tinny)

A very low resolution, rounded and friendly-looking font with a fairly decent latin character set. This font was the second one created for the ggfnt format, having two main goals:
- Create the first proper font for ggfnt, covering at least the basic ASCII range.
- Test kernings for ggfnt and ptxt.
There's also a [bold version](https://github.com/tinne26/ggfnt-fonts/tree/main/tinnybold) of the font that can be used alongside it.

The name was suggested by [Zyko](https://github.com/Zyko0) as a pun for tiny + tinne. He regretted the silly joke almost instantly, but it was already too late... and the name stuck.

The import path is:
```Golang
import "github.com/tinne26/ggfnt-fonts/tinny"
```

## Glyphs

Common glyphs:
- Notdef.
- Full ASCII range.
- Basic diacritics: `àáäâ` (for `aeiou`, `AEIOU` and the diacritics in isolation).
- Other common punctuation symbols `¡ ¿`, `– ‑ —`, `‘’ “” …`.
- Spanish letters `ÑñÇç`.
- Common currency symbols `€ £ ¢ ¥ ¤`.
- A few extra math symbols `− × ÷ ± º ′ ″ π`.
- A few extra symbols `♩ ♪ ♫ ♬ �`.

Named glyphs:
- `"notdef"`, also mapped to `'\uE000'` for accessibility.

## License

The code in this folder is MIT licensed, the font itself is [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
