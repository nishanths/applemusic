package applemusic

import (
	"os"
	"testing"
)

func TestParseHTML(t *testing.T) {
	testcases := []struct {
		path   string
		expect Info
	}{
		{
			"testdata/album-rumours-deluxe-651871544-651871679.html",
			Info{
				Artwork: Artwork{
					HttpURL:  "https://is1-ssl.mzstatic.com/image/thumb/Music2/v4/69/ef/ca/69efca19-5e7a-67a1-a9c4-a748fa8b3db6/603497925759.jpg/1200x630bb.jpg",
					HttpsURL: "https://is1-ssl.mzstatic.com/image/thumb/Music2/v4/69/ef/ca/69efca19-5e7a-67a1-a9c4-a748fa8b3db6/603497925759.jpg/1200x630bb.jpg",
					Type:     "image/jpg",
					Width:    1200,
					Height:   630,
				},
				AlbumURL:  "https://itunes.apple.com/us/album/rumours-deluxe/651871544",
				ArtistURL: "https://itunes.apple.com/us/artist/fleetwood-mac/158038",
			},
		},
	}

	for i, tt := range testcases {
		r, err := os.Open(tt.path)
		if err != nil {
			t.Errorf("[%d]: failed to open %q", i, tt.path)
			continue
		}
		defer r.Close()
		got, err := ParseHTML(r)
		if err != nil {
			t.Errorf("[%d]: unexpected error %s", i, err)
			continue
		}
		if tt.expect != got {
			t.Errorf("[%d]: expected: %+v, got: %+v", i, tt.expect, got)
		}
	}
}
