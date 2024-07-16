# minitile
[![Go Reference](https://pkg.go.dev/badge/tinne26/ggfnt-fonts/minitile.svg)](https://pkg.go.dev/github.com/tinne26/ggfnt-fonts/minitile)

An uppercase-only negative font where the glyphs are contained within slightly rounded squares of 7x7. This font was the first one encoded as a ggfnt to test vertical metrics and layouts.

One could use it to implement crossword puzzles and similar games, but it's unclear if that would be much better than using bitmaps directly, given a grid context.

This is the ugliest font I've ever created. It's so ugly that I initially named it "monotile", but later changed it in case I ever wanted to make a better "monotile" font more deserving of the name.

The import path is:
```Golang
import "github.com/tinne26/ggfnt-fonts/minitile"
```

## Glyphs

Common glyphs:
- Notdef.
- Full ASCII range except for lowercase letters `a-z`.
- Other common punctuation symbols `¡ ¿ — …`.
- Spanish letters `ÑÇ`.

Named glyphs:
- `"notdef"`, also mapped to `'\uE000'` for accessibility.

## License

The code in this folder is MIT licensed, the font itself is [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
