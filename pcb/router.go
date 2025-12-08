package pcb

const MaxRoutingDistance = 3.0
const DefaultTraceWidth = 0.2

func RouteBoard(board *Board) {
	netMap := make(map[int]bool)
	for _, pad := range board.Pads {
		netMap[pad.Net.Number] = true
	}

	for netNum := range netMap {
		pads := board.GetPadsByNet(netNum)

		for i := 0; i < len(pads); i++ {
			for j := i + 1; j < len(pads); j++ {
				if pads[i].Distance(pads[j]) <= MaxRoutingDistance &&
					pads[i].Layer == pads[j].Layer {
					seg := Segment{
						Start: pads[i].Position,
						End:   pads[j].Position,
						Width: DefaultTraceWidth,
						Layer: pads[i].Layer,
						Net:   netNum,
					}
					board.AddSegment(seg)
				}
			}
		}
	}
}
