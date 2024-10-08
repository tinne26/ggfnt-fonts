package starship

import "testing"

import "github.com/tinne26/ggfnt"

func TestFont(t *testing.T) {
	// test initial state and parsing
	if cachedFont != nil { t.Fatal("cachedFont != nil") }
	font := Font() // if this panics it will come out on the test
	if cachedFont == nil { t.Fatal("cachedFont == nil") }

	// ensure font name is what we expected
	const ExpectedFontName = "starship"
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
	var settingFound bool
	font.Settings().Each(func(key ggfnt.SettingKey, name string) {
		switch name {
		case CutsSettingName:
			if settingFound { panic("broken font") }
			if key != CutsSettingKey {
				t.Fatalf("expected CutsSettingKey to be %d, found %d", CutsSettingKey, key)
			}
			settingFound = true
		}
	})
	if !settingFound {
		t.Fatalf("expected to find Cuts setting as '%s'", CutsSettingName)
	}

	// test release
	Release()
	if cachedFont != nil {
		t.Fatal("font release failed")
	}
}
