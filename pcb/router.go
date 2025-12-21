package pcb

import (
	"log/slog"

	"github.com/mackeper/lin_router/utils"
)

const DefaultTraceWidth = 0.2

func AddTrivialSegments(board *Board, maxRoutingDistance float64) {
	netMap := make(map[int]bool)
	for _, pad := range board.Pads {
		netMap[pad.Net.Number] = true
	}
	for _, via := range board.Vias {
		netMap[via.Net] = true
	}

	slog.Debug("Router starting", "total_pads", len(board.Pads), "total_vias", len(board.Vias), "nets", len(netMap))

	for netNum := range netMap {
		if netNum == 0 {
			slog.Debug("Skipping net 0 (no net)")
			continue
		}
		pads := board.GetPadsByNet(netNum)
		vias := board.GetViasByNet(netNum)
		slog.Debug("Processing net", "net", netNum, "pads", len(pads), "vias", len(vias))

		// Pad to Pad
		for i := range pads {
			for j := i + 1; j < len(pads); j++ {
				dist := pads[i].Distance(pads[j])
				if dist <= maxRoutingDistance {
					sharedLayers := getSharedLayers(pads[i].Layers, pads[j].Layers)
					slog.Debug("Found pad pair within distance", "net", netNum, "dist", dist, "shared_layers", len(sharedLayers), "pad1_layers", pads[i].Layers, "pad2_layers", pads[j].Layers)
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
						slog.Debug("Added segment", "net", netNum, "layer", layer)
					}
				}
			}
		}

		// Pad to Via
		for _, pad := range pads {
			for _, via := range vias {
				if pad.Position.Distance(via.Position) <= maxRoutingDistance {
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
				if vias[i].Distance(vias[j]) <= maxRoutingDistance {
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
