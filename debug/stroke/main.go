package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("Drawing Test"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	ops := new(op.Ops)

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			var p clip.Path
			p.Begin(gtx.Ops)
			p.MoveTo(f32.Pt(70, 20))
			p.CubeTo(f32.Pt(70, 20), f32.Pt(70, 110), f32.Pt(120, 120))
			p.LineTo(f32.Pt(20, 120))
			p.Close()
			ps := p.End()

			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, clip.Stroke{Path: ps, Width: 20}.Op())

			e.Frame(gtx.Ops)

		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
