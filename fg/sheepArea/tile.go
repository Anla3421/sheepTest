package sheepArea

type Tile struct {
	ID        int
	Symbol    string
	X, Y, Z   int
	IsBlocked bool
	IsMatched bool

	// --- 新增的依賴關係 ---
	// Covers: 這塊牌直接蓋住了哪些牌 (它的孩子們)
	Covers []*Tile
	// CoveredBy: 這塊牌被哪些牌直接蓋住 (它的父母們)
	CoveredBy []*Tile
}
