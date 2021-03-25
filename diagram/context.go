package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Context struct {
	Transform
	layout.Context
	Theme   *material.Theme
	Diagram *Diagram
}

func NewContext(gtx layout.Context, th *material.Theme, zoom *Zoom, diagram *Diagram) *Context {
	return &Context{
		Transform: NewTransform(gtx, zoom),
		Context:   gtx,
		Theme:     th,
		Diagram:   diagram,
	}
}

type Transform struct {
	Dp        int
	PxPerUnit int
}

func NewTransform(gtx layout.Context, zoom *Zoom) Transform {
	px := gtx.Px(unit.Dp(30))
	px = (px / 24) * 24 // make it divisible by 2,3,4,6,12
	px = int(float32(px) * zoom.Multiplier())
	return Transform{
		Dp:        gtx.Px(unit.Dp(1)),
		PxPerUnit: px,
	}
}

func (tr *Transform) Px(v Unit) int {
	return tr.PxPerUnit * int(v)
}

func (tr *Transform) Pt(v Vector) image.Point {
	return image.Point{
		X: int(v.X) * tr.PxPerUnit,
		Y: int(v.Y) * tr.PxPerUnit,
	}
}

func (tr *Transform) FPt(v Vector) f32.Point {
	return f32.Point{
		X: float32(int(v.X) * tr.PxPerUnit),
		Y: float32(int(v.Y) * tr.PxPerUnit),
	}
}

func (tr *Transform) Inv(p image.Point) Vector {
	return Vector{
		X: Unit(p.X / tr.PxPerUnit),
		Y: Unit(p.Y / tr.PxPerUnit),
	}
}

func (tr *Transform) FInv(p f32.Point) Vector {
	return Vector{
		X: Unit(int(p.X) / tr.PxPerUnit),
		Y: Unit(int(p.Y) / tr.PxPerUnit),
	}
}

func (tr *Transform) Bounds(box Box) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(box.Pos.X) * tr.PxPerUnit,
			Y: int(box.Pos.Y) * tr.PxPerUnit,
		},
		Max: image.Point{
			X: int(box.Pos.X+box.Size.X) * tr.PxPerUnit,
			Y: int(box.Pos.Y+box.Size.Y) * tr.PxPerUnit,
		},
	}
}
