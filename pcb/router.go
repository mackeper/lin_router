package pcb

import "github.com/mackeper/lin_router/utils"

const MaxRoutingDistance = 3.0
const DefaultTraceWidth = 0.2

func AddTrivialSegments(board *Board) {
	netMap := make(map[int]bool)
	for _, pad := range board.Pads {
		netMap[pad.Net.Number] = true
	}

	for netNum := range netMap {
		pads := board.GetPadsByNet(netNum)

		for i := range pads {
			for j := i + 1; j < len(pads); j++ {
				if pads[i].Distance(pads[j]) <= MaxRoutingDistance {
					sharedLayers := pads[i].GetSharedCopperLayers(pads[j])
					for _, layer := range sharedLayers {
						seg := Segment{
							Start: pads[i].Position,
							End:   pads[j].Position,
							Width: DefaultTraceWidth,
							Layer: layer,
							Net:   netNum,
							UUID:  utils.GenerateUUID(),
						}
						board.AddSegment(seg)
					}
				}
			}
		}
	}
}
