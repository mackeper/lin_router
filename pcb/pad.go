package pcb

type Pad struct {
	Position Position
	Net      Net
	Number   string
	Layers   []string
}

func (p Pad) Distance(other Pad) float64 {
	return p.Position.Distance(other.Position)
}
