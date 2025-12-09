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

func (p Pad) GetSharedCopperLayers(other Pad) []string {
	shared := []string{}
	for _, layer1 := range p.Layers {
		if layer1 == "F.Cu" || layer1 == "B.Cu" {
			for _, layer2 := range other.Layers {
				if layer1 == layer2 {
					shared = append(shared, layer1)
				}
			}
		}
	}
	return shared
}
