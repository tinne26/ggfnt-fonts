package candy

import "testing"

func TestFont(t *testing.T) {
	// test initial state and parsing
	if cachedFont != nil { t.Fatal("cachedFont != nil") }
	font := Font() // if this panics it will come out on the test
	if cachedFont == nil { t.Fatal("cachedFont == nil") }

	// ensure font name is what we expected
	const ExpectedFontName = "candy"
	name := font.Header().Name()
	if name != ExpectedFontName {
		t.Fatalf("expected font name to be \"" + ExpectedFontName + "\", got \"%s\" instead", name)
	}

	// check hardcoded glyph ids
	mappingGroup, found := font.Mapping().Utf8(NotdefRune, nil)
	if !found { t.Fatalf("mapping for NotdefRune not found") }
	if mappingGroup.Size() != 1 {
		t.Fatalf("expected mapping for NotdefRune to have size 1, found %d", mappingGroup.Size())
	}
	notdefGlyphIndex := mappingGroup.Select(0)
	if notdefGlyphIndex != Notdef {
		t.Fatalf("expected Notdef to be mapped to %d, found %d", Notdef, notdefGlyphIndex)
	}

	mappingGroup, found = font.Mapping().Utf8(FatDotRune, nil)
	if !found { t.Fatalf("mapping for FatDotRune not found") }
	if mappingGroup.Size() != 1 {
		t.Fatalf("expected mapping for FatDotRune to have size 1, found %d", mappingGroup.Size())
	}
	fatDotGlyphIndex := mappingGroup.Select(0)
	if fatDotGlyphIndex != FatDot {
		t.Fatalf("expected FatDot to be mapped to %d, found %d", FatDot, fatDotGlyphIndex)
	}

	mappingGroup, found = font.Mapping().Utf8(CandyRune, nil)
	if !found { t.Fatalf("mapping for CandyRune not found") }
	if mappingGroup.Size() != 1 {
		t.Fatalf("expected mapping for CandyRune to have size 1, found %d", mappingGroup.Size())
	}
	candyGlyphIndex := mappingGroup.Select(0)
	if candyGlyphIndex != Candy {
		t.Fatalf("expected Candy to be mapped to %d, found %d", Candy, candyGlyphIndex)
	}

	mappingGroup, found = font.Mapping().Utf8(DyeCandyRune, nil)
	if !found { t.Fatalf("mapping for DyeCandyRune not found") }
	if mappingGroup.Size() != 1 {
		t.Fatalf("expected mapping for DyeCandyRune to have size 1, found %d", mappingGroup.Size())
	}
	dyeCandyGlyphIndex := mappingGroup.Select(0)
	if dyeCandyGlyphIndex != DyeCandy {
		t.Fatalf("expected DyeCandy to be mapped to %d, found %d", DyeCandy, dyeCandyGlyphIndex)
	}

	// test release
	Release()
	if cachedFont != nil {
		t.Fatal("font release failed")
	}
}
