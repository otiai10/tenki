// main
package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	"os"

	// _ "image/png"
	// _ "image/gif"

	"github.com/otiai10/gat/render"
	"github.com/otiai10/tenki/cli"
	"github.com/otiai10/tenki/tenki"
)

const version = "v1.2.3"

var (
	geo, mask bool
	usepix    bool
	scale     float64

	// debug
	// verbose bool

	// 以下、タイムラプスでのみ有効
	lapse   bool
	minutes int
	delay   int
	loop    bool

	// リストモード
	list bool
)

func setup() {
	flag.BoolVar(&lapse, "a", false, "タイムラプス表示")
	flag.IntVar(&minutes, "m", 40, "タイムラプスの取得直近時間（分）")
	flag.IntVar(&delay, "d", 200, "タイムラプスアニメーションのfps（msec）")
	flag.BoolVar(&loop, "l", false, "タイムラプスアニメーションをループ表示")
	flag.BoolVar(&geo, "g", true, "地形を描画")
	flag.BoolVar(&mask, "b", true, "県境を描画")
	flag.BoolVar(&usepix, "p", false, "iTermであってもピクセル画で表示")
	flag.Float64Var(&scale, "s", 1.2, "表示拡大倍率")
	flag.BoolVar(&list, "list", false, "サポートしている地域を一覧")
	// flag.BoolVar(&verbose, "v", false, "デバッグログ表示")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "tenki.jpをCLIに表示するコマンドです。(%v)\n利用可能なオプション:\n", version)
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {

	setup()

	if list {
		onerror(cli.List())
		return
	}

	renderer := render.GetDefaultRenderer()
	if usepix {
		renderer = &render.CellGrid{}
	}
	renderer.SetScale(scale)
	areaname := flag.Arg(0)
	if areaname == "" {
		areaname = "japan"
	}
	area, err := tenki.GetArea(areaname)
	onerror(err)

	if lapse {
		onerror(cli.Timelapse(renderer, area, minutes, delay, loop))
	} else {
		onerror(cli.Tenki(renderer, area))
	}
}

func onerror(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
