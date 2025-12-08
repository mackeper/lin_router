package pcb

type Segment struct {
	Start Position
	End   Position
	Width float64
	Layer string
	Net   int
}

func (s Segment) Length() float64 {
	return s.Start.Distance(s.End)
}
