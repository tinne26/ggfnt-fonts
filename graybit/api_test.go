package graybit

import "testing"

import "github.com/tinne26/ggfnt"

func TestFont(t *testing.T) {
	// test initial state and parsing
	if cachedFont != nil { t.Fatal("cachedFont != nil") }
	font := Font() // if this panics it will come out on the test
	if cachedFont == nil { t.Fatal("cachedFont == nil") }

	// ensure font name is what we expected
	const ExpectedFontName = "graybit"
	name := font.Header().Name()
	if name != ExpectedFontName {
		t.Fatalf("expected font name to be \"" + ExpectedFontName + "\", got \"%s\" instead", name)
	}

	// check hardcoded glyph ids
	mappingGroup, found := font.Mapping().Utf8(NotdefRune, nil)
	if !found {
		t.Fatalf("mapping for NotdefRune not found")
	}
	if mappingGroup.Size() != 1 {
		t.Fatalf("expected mapping for NotdefRune to have size 1, found %d", mappingGroup.Size())
	}
	notdefGlyphIndex := mappingGroup.Select(0)
	if notdefGlyphIndex != Notdef {
		t.Fatalf("expected Notdef to be mapped to %d, found %d", Notdef, notdefGlyphIndex)
	}
	
	mappingGroup, found = font.Mapping().Utf8(LowHyphenRune, nil)
	if !found {
		t.Fatalf("mapping for LowHyphenRune not found")
	}
	if mappingGroup.Size() != 1 {
		t.Fatalf("expected mapping for LowHyphenRune to have size 1, found %d", mappingGroup.Size())
	}
	lowHyphenGlyphIndex := mappingGroup.Select(0)
	if lowHyphenGlyphIndex != LowHyphen {
		t.Fatalf("expected LowHyphen to be mapped to %d, found %d", LowHyphen, lowHyphenGlyphIndex)
	}

	// check setting names and keys
	var zeroDisMarkFound, numStyleFound bool
	font.Settings().Each(func(key ggfnt.SettingKey, name string) {
		switch name {
		case ZeroDisambiguationMarkSettingName:
			if zeroDisMarkFound { panic("broken font") }
			if key != ZeroDisambiguationMarkSettingKey {
				t.Fatalf("expected ZeroDisambiguationMarkSettingKey to be %d, found %d", ZeroDisambiguationMarkSettingKey, key)
			}
			zeroDisMarkFound = true
		case NumericStyleSettingName:
			if numStyleFound { panic("broken font") }
			if key != NumericStyleSettingKey {
				t.Fatalf("expected NumericStyleSettingKey to be %d, found %d", NumericStyleSettingKey, key)
			}
			numStyleFound = true
		}
	})
	if !zeroDisMarkFound {
		t.Fatalf("expected to find ZeroDisambiguationMark setting as '%s'", ZeroDisambiguationMarkSettingName)
	}
	if !numStyleFound {
		t.Fatalf("expected to find NumericStyle setting as '%s'", NumericStyleSettingName)
	}

	// test release
	Release()
	if cachedFont != nil {
		t.Fatal("font release failed")
	}
}
