package tenki

import (
	"encoding/json"
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

var (
	Unit time.Duration = 5 * time.Minute
	area               = map[string]string{}
)

const defaultLocation = "Asia/Tokyo"

// ゆうて日本国内なら同じtimezoneなのであんまり問題無いと思うけど、
// 今後世界都市対応するなら、これは問題になりますね。
func GetNow() (time.Time, error) {
	location, err := time.LoadLocation(defaultLocation)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(location), nil
}

func ListAreas() (names []string) {
	for k := range area {
		names = append(names, k)
	}
	return names
}

func Load(jsonlike []byte) error {
	return json.Unmarshal(jsonlike, &area)
}

func TruncateTime(t time.Time, unit time.Duration) time.Time {
	return t.Add(-1 * unit).Round(1 * unit)
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
		return Area{}, fmt.Errorf("地域名 %v はサポートされていません", name)
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
	if res.StatusCode >= 400 {
		fmt.Println(res.Status, res.Request.URL.String())
		return nil, fmt.Errorf(res.Status)
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	return img, err
}
