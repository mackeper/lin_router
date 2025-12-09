package pcb

import "github.com/mackeper/lin_router/utils"

const MaxRoutingDistance = 3.0
const DefaultTraceWidth = 0.2

func AddTrivialSegments(board *Board) {
	netMap := make(map[int]bool)
	for _, pad := range board.Pads {
		netMap[pad.Net.Number] = true
	}
	for _, via := range board.Vias {
		netMap[via.Net] = true
	}

	for netNum := range netMap {
		pads := board.GetPadsByNet(netNum)
		vias := board.GetViasByNet(netNum)

		// Pad to Pad
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

		// Pad to Via
		for _, pad := range pads {
			for _, via := range vias {
				if pad.Position.Distance(via.Position) <= MaxRoutingDistance {
					sharedLayers := getSharedLayers(pad.Layers, via.Layers)
					for _, layer := range sharedLayers {
						seg := Segment{
							Start: pad.Position,
							End:   via.Position,
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

		// Via to Via
		for i := range vias {
			for j := i + 1; j < len(vias); j++ {
				if vias[i].Distance(vias[j]) <= MaxRoutingDistance {
					sharedLayers := getSharedLayers(vias[i].Layers, vias[j].Layers)
					for _, layer := range sharedLayers {
						seg := Segment{
							Start: vias[i].Position,
							End:   vias[j].Position,
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

func getSharedLayers(layers1, layers2 []string) []string {
	layerMap := make(map[string]bool)
	for _, layer := range layers1 {
		if isCopperLayer(layer) {
			layerMap[layer] = true
		}
	}

	var shared []string
	for _, layer := range layers2 {
		if isCopperLayer(layer) && layerMap[layer] {
			shared = append(shared, layer)
		}
	}
	return shared
}

func isCopperLayer(layer string) bool {
	return layer == "F.Cu" || layer == "B.Cu"
}
