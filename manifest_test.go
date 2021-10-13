package go_bagit

import (
	"path/filepath"
	"testing"
)

func TestManifests(t *testing.T) {
	t.Run("Test Parsing Manifest", func(t *testing.T) {
		manifestLoc := filepath.Join("test", "valid", "manifest-sha256.txt")
		manifestMap, err := ReadManifest(manifestLoc)
		if err != nil {
			t.Error(err)
		}

		want := 1
		got := len(manifestMap)
		if want != got {
			t.Errorf("Wanted %d Got %d", want, got)
		}
	})

	t.Run("Test Validate a manifest file", func(t *testing.T) {
		manifestLoc := filepath.Join("test", "valid", "manifest-sha256.txt")
		if err := ValidateManifest(manifestLoc, false); len(err) > 0 {
			t.Error(err[0])
		}
	})

	t.Run("Test Completeness Only Validation of a manifest file", func(t *testing.T) {
		manifestLoc := filepath.Join("test", "valid", "manifest-sha256.txt")
		if err := ValidateManifest(manifestLoc, true); len(err) > 0 {
			t.Error(err[0])
		}
	})
}
