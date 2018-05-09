package applemusic

import (
	"errors"
	"io"
	"strconv"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Info struct {
	Artwork   Artwork
	AlbumURL  string
	ArtistURL string
}

type Artwork struct {
	HttpURL       string
	HttpsURL      string
	Type          string
	Width, Height int
}

// https://itunes.apple.com/us/album/651871544?i=651871679
// <meta> tags attributes
// ----------------------
//
// property og:image
// content https://is1-ssl.mzstatic.com/image/thumb/Music2/v4/69/ef/ca/69efca19-5e7a-67a1-a9c4-a748fa8b3db6/603497925759.jpg/1200x630bb.jpg
//
// property og:image:secure_url
// content https://is1-ssl.mzstatic.com/image/thumb/Music2/v4/69/ef/ca/69efca19-5e7a-67a1-a9c4-a748fa8b3db6/603497925759.jpg/1200x630bb.jpg
//
// property og:image:type
// content image/jpg
//
// property og:image:width
// content 1200
//
// property og:image:height
// content 630
//
// property music:musician
// content https://itunes.apple.com/us/artist/fleetwood-mac/158038

// ParseHTML parses info from HTML pages such as
// https://itunes.apple.com/us/album/651871544?i=651871679.
func ParseHTML(r io.Reader) (Info, error) {
	z := html.NewTokenizer(r)
	var i Info

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// should never reach io.EOF in practice since
			// we exit earlier at </head>, so don't filter io.EOF
			// as an error. If we do reach io.EOF it's better to fail
			// so we can figure out what's going on.
			return Info{}, z.Err()

		case html.EndTagToken:
			// reached </head>, there shouldn't be <meta> tags after this
			// point.
			if z.Token().DataAtom == atom.Head {
				return i, nil
			}

		case html.StartTagToken:
			// Note that <meta> elements don't have end tags.
			// TODO: why is z.Token() empty?
			// Hence we're doing manual accumulation of attributes below.
			tn, _ := z.TagName()

			if atom.String(tn) == "meta" {
				var as []html.Attribute
				for {
					k, v, more := z.TagAttr()
					if !more {
						break
					}
					as = append(as, html.Attribute{Key: string(k), Val: string(v)}) // ignore Namespace
				}

				pv := attrVal(as, "property")
				nv := attrVal(as, "name")
				cv := attrVal(as, "content")

				if (pv == "" && nv == "") || cv == "" {
					continue
				}

				switch pv {
				case "og:image":
					i.Artwork.HttpURL = cv
				case "og:image:secure_url":
					i.Artwork.HttpsURL = cv
				case "og:image:type":
					i.Artwork.Type = cv
				case "og:image:width":
					if w, err := strconv.Atoi(cv); err == nil {
						i.Artwork.Width = w
					}
				case "og:image:height":
					if h, err := strconv.Atoi(cv); err == nil {
						i.Artwork.Height = h
					}
				case "music:musician":
					i.ArtistURL = cv
				}

				switch nv {
				case "music:album":
					i.AlbumURL = cv
				} // ... possibly more cases ...
			}
		}
	}

	// to make the compiler happy
	return Info{}, errors.New("should never be reached")
}

func attrVal(as []html.Attribute, key string) string {
	for _, a := range as {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
