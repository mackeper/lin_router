package lexer

type ExprType int

const (
	ExprUnknown ExprType = iota
	ExprKicadPcb
	ExprFootprint
	ExprPad
	ExprSegment
	ExprVia
	ExprNet
	ExprLayer
	ExprGrLine
	ExprGrArc
	ExprAt
)

func (et ExprType) String() string {
	switch et {
	case ExprKicadPcb:
		return "kicad_pcb"
	case ExprFootprint:
		return "footprint"
	case ExprPad:
		return "pad"
	case ExprSegment:
		return "segment"
	case ExprVia:
		return "via"
	case ExprNet:
		return "net"
	case ExprLayer:
		return "layer"
	case ExprGrLine:
		return "gr_line"
	case ExprGrArc:
		return "gr_arc"
	case ExprAt:
		return "at"
	default:
		return "unknown"
	}
}

func IdentifierToExprType(identifier string) ExprType {
	switch identifier {
	case "kicad_pcb":
		return ExprKicadPcb
	case "footprint":
		return ExprFootprint
	case "pad":
		return ExprPad
	case "segment":
		return ExprSegment
	case "via":
		return ExprVia
	case "net":
		return ExprNet
	case "layer":
		return ExprLayer
	case "gr_line":
		return ExprGrLine
	case "gr_arc":
		return ExprGrArc
	case "at":
		return ExprAt
	default:
		return ExprUnknown
	}
}
