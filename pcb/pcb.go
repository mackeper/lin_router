package pcb

import (
	"math"
)

type Position struct {
	X float64
	Y float64
}

func (p Position) Distance(other Position) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type Board struct {
	Pads     []Pad
	Segments []Segment
	Vias     []Via
}

func NewBoard() *Board {
	return &Board{
		Pads:     []Pad{},
		Segments: []Segment{},
		Vias:     []Via{},
	}
}

func (b *Board) AddPad(pad Pad) {
	b.Pads = append(b.Pads, pad)
}

func (b *Board) AddSegment(seg Segment) {
	b.Segments = append(b.Segments, seg)
}

func (b *Board) AddVia(via Via) {
	b.Vias = append(b.Vias, via)
}

func (b *Board) GetPadsByNet(netNum int) []Pad {
	var pads []Pad
	for _, pad := range b.Pads {
		if pad.Net.Number == netNum {
			pads = append(pads, pad)
		}
	}
	return pads
}

func (b *Board) GetViasByNet(netNum int) []Via {
	var vias []Via
	for _, via := range b.Vias {
		if via.Net == netNum {
			vias = append(vias, via)
		}
	}
	return vias
}
