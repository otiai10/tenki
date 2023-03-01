package cli

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/otiai10/gat/render"
	"github.com/otiai10/tenki/tenki"
)

// Tenki
func Tenki(r render.Renderer, area tenki.Area) error {
	t, err := tenki.GetNow()
	if err != nil {
		return err
	}
	t = tenki.TruncateTime(t, tenki.Unit)
	entry := area.GetEntry(t)
	if err != nil {
		return err
	}
	img, err := entry.Image()
	if err != nil {
		return err
	}
	if err := r.Render(os.Stdout, img); err != nil {
		return err
	}
	fmt.Printf("Powered by %s\n", area.ReferenceURL)
	return nil
}

// Timelapse タイムラプス表示
func Timelapse(r render.Renderer, area tenki.Area, minutes, delay int, loop bool) error {

	fmt.Printf("直近%d分間の降雨画像を取得中", minutes)

	now, err := tenki.GetNow()
	if err != nil {
		return err
	}
	now = tenki.TruncateTime(now, tenki.Unit)

	t := now.Add(time.Duration(-1*minutes) * time.Minute)
	end := now

	var entries []tenki.Entry
	entries = append(entries, area.GetEntry(t))
	for t := t.Add(tenki.Unit); t.Before(end); t = t.Add(tenki.Unit) {
		entries = append(entries, area.GetEntry(t))
	}

	progress := func(i int) { fmt.Print(".") }

	// images := make([]*image.RGBA, len(entries), len(entries))
	images := make([]image.Image, len(entries))
	for i, entry := range entries {
		img, err := entry.Image()
		if err != nil {
			return err
		}
		images[i] = img
		progress(i)
	}

	// まずクリアする
	fmt.Printf("\033c")

	var moveCursorToTop = func() {
		fmt.Print("\033[s\033[H\033[1;32m")
	}

	for i, img := range images {
		moveCursorToTop()
		r.Render(os.Stdout, img)
		fmt.Fprintln(os.Stdout, entries[i].Time.String())
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	fmt.Print("\033[0m")

	return nil
}
