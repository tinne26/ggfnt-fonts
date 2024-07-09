package omen

import "testing"
import "github.com/tinne26/ggfnt"

func TestFont(t *testing.T) {
	// test initial state and parsing
	if cachedFont != nil { t.Fatal("cachedFont != nil") }
	font := Font() // if this panics it will come out on the test
	if cachedFont == nil { t.Fatal("cachedFont == nil") }

	// ensure font name is what we expected
	const ExpectedFontName = "omen"
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

	// check setting names and keys
	var oStyleSettingFound bool
	font.Settings().Each(func(key ggfnt.SettingKey, name string) {
		switch name {
		case OStyleSettingName:
			if oStyleSettingFound { panic("broken font") }
			if key != OStyleSettingKey {
				t.Fatalf("expected OStyleSettingKey to be %d, found %d", OStyleSettingKey, key)
			}
			oStyleSettingFound = true
		}
	})
	if !oStyleSettingFound {
		t.Fatalf("expected to find OStyle setting as '%s'", OStyleSettingName)
	}

	// test release
	Release()
	if cachedFont != nil {
		t.Fatal("font release failed")
	}
}
