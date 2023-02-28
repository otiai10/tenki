package tenki

import (
	"fmt"
	"image"
	"net/http"
	"time"
)

const (
	// TenkiStaticURL ...
	TenkiStaticURL = "https://static.tenki.jp/static-images"
	// TenkiDynamicTimestampPath ...
	TenkiDynamicTimestampPath = "/radar/2006/01/02/15/04/00"
)

func TruncateTime(t time.Time, unit time.Duration) time.Time {
	return t.Add(-2 * unit).Round(1 * unit)
}

type Area struct {
	AreaURLPath  string
	ReferenceURL string
}

type Entry struct {
	Time time.Time
	URL  string
}

func GetArea(name string) (Area, error) {
	p, ok := area[name]
	if !ok {
		return Area{}, fmt.Errorf("地域名 %v はサポートされていません.", name)
	}
	return Area{
		AreaURLPath:  p,
		ReferenceURL: "https://tenki.jp/",
	}, nil
}

func (area Area) GetEntry(t time.Time) Entry {
	return Entry{
		URL:  TenkiStaticURL + t.Format(TenkiDynamicTimestampPath) + area.AreaURLPath,
		Time: t,
	}
}

func (entry Entry) Image(client ...*http.Client) (image.Image, error) {
	if len(client) == 0 {
		client = append(client, http.DefaultClient)
	}
	res, err := client[0].Get(entry.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	return img, err
}
