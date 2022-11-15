package atlas

type TileType string

const (
	Base    TileType = "Base"
	Overlay TileType = "Overlay"
)

func (t TileType) toString() string {
	if t == Base {
		return "Base"
	}
	if t == Overlay {
		return "Overlay"
	}
	panic("unknown TileType")
}
