package main

import (
	"image"
	"image/color"
	"time"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Hoverable struct {
	hovered bool
}

func (h *Hoverable) Layout(gtx layout.Context) layout.Dimensions {
	h.update(gtx)

	stack := op.Push(gtx.Ops)
	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   h,
		Types: pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)
	stack.Pop()

	return layout.Dimensions{
		Size: gtx.Constraints.Min,
	}
}

func (h *Hoverable) Active() bool { return h.hovered }

func (h *Hoverable) update(gtx layout.Context) {
	for _, e := range gtx.Events(h) {
		ev, ok := e.(pointer.Event)
		if !ok {
			continue
		}

		switch ev.Type {
		case pointer.Enter:
			h.hovered = true
		case pointer.Leave:
			h.hovered = false
		}
	}
}

type AnimationTimer struct {
	Duration time.Duration

	progress time.Duration
	last     time.Time
}

func (anim *AnimationTimer) Progress() float32 {
	return float32(anim.progress) / float32(anim.Duration)
}

func (anim *AnimationTimer) Update(gtx layout.Context, active bool) float32 {
	delta := gtx.Now.Sub(anim.last)
	anim.last = gtx.Now

	if delta > 15*time.Millisecond {
		delta = 15 * time.Millisecond
	}

	if active {
		if anim.progress < anim.Duration {
			anim.progress += delta
			op.InvalidateOp{}.Add(gtx.Ops)
			if anim.progress > anim.Duration {
				anim.progress = anim.Duration
			}
		}
	} else {
		if anim.progress > 0 {
			anim.progress -= delta
			op.InvalidateOp{}.Add(gtx.Ops)
			if anim.progress < 0 {
				anim.progress = 0
			}
		}
	}

	return anim.Progress()
}

// BorderSmooth lays out a widget and draws a border inside it, with non-pixel-perfect border.
type BorderSmooth struct {
	Color        color.NRGBA
	CornerRadius unit.Value
	Width        float32
}

func (b BorderSmooth) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := w(gtx)
	sz := dims.Size
	rr := float32(gtx.Px(b.CornerRadius))
	st := op.Push(gtx.Ops)
	clip.Border{
		Rect: f32.Rectangle{
			Max: layout.FPt(sz),
		},
		NE: rr, NW: rr, SE: rr, SW: rr,
		Width: b.Width,
	}.Add(gtx.Ops)
	paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return dims
}
