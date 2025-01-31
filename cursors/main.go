// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"github.com/egonelbre/expgio/f32color"
)

const cursorCount = pointer.CursorNorthWestSouthEastResize + 1

func main() {
	ui := &UI{Theme: material.NewTheme()}
	go func() {
		w := app.NewWindow(app.Title("Image Viewer"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	Theme *material.Theme
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			key.InputOp{Tag: w, Keys: key.NameEscape}.Add(gtx.Ops)
			for _, ev := range gtx.Queue.Events(w) {
				if e, ok := ev.(key.Event); ok && e.Name == key.NameEscape {
					return nil
				}
			}

			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	n := float64(cursorCount)

	ratio := float64(gtx.Constraints.Max.X) / float64(gtx.Constraints.Max.Y)
	columnsBest := math.Sqrt(n * ratio)
	rowsBest := columnsBest / ratio

	cols := int(math.Ceil(columnsBest))
	rows := int(math.Ceil(rowsBest))

	squareSize := gtx.Constraints.Max.X / cols
	square := image.Point{X: squareSize, Y: squareSize}

	i := pointer.Cursor(0)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if i >= cursorCount {
				break
			}
			func(cursor pointer.Cursor) {
				p := image.Point{X: col * squareSize, Y: row * squareSize}
				defer op.Offset(p).Push(gtx.Ops).Pop()
				defer clip.Rect{Max: square}.Push(gtx.Ops).Pop()

				col := f32color.HSL(float32(i)*math.Phi, 0.6, 0.6)
				paint.ColorOp{Color: col}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)

				pointer.InputOp{Tag: i}.Add(gtx.Ops)
				cursor.Add(gtx.Ops)

				gtx := gtx
				gtx.Constraints = layout.Exact(square)
				layout.Center.Layout(gtx,
					material.Body1(ui.Theme, cursor.String()).Layout)
			}(i)
			i++
		}
	}

	return layout.Dimensions{
		Size: image.Pt(squareSize*cols, squareSize*rows),
	}
}
