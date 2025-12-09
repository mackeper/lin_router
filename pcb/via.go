package pcb

type Via struct {
	Position Position
	Layers   []string
	Net      int
	UUID     string
}

func (v Via) Distance(other Via) float64 {
	return v.Position.Distance(other.Position)
}
