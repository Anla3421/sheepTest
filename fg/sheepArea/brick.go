package sheepArea

// LvOneTiles is a map of Tiles generated from an Excel file.
var LvOneTiles = map[int]*Tile{
	1:  {ID: 1, Symbol: "L2", X: 0, Y: 4, Z: 3, IsBlocked: false, IsMatched: false}, //
	2:  {ID: 2, Symbol: "L3", X: -1, Y: 3, Z: 2, IsBlocked: true, IsMatched: false},
	3:  {ID: 3, Symbol: "L3", X: 1, Y: 3, Z: 2, IsBlocked: true, IsMatched: false},
	4:  {ID: 4, Symbol: "L1", X: 0, Y: 2, Z: 1, IsBlocked: true, IsMatched: false},
	5:  {ID: 5, Symbol: "L1", X: -3, Y: 1, Z: 2, IsBlocked: true, IsMatched: false},
	6:  {ID: 6, Symbol: "L2", X: -1, Y: 1, Z: 0, IsBlocked: true, IsMatched: false},
	7:  {ID: 7, Symbol: "L4", X: 1, Y: 1, Z: 0, IsBlocked: true, IsMatched: false},
	8:  {ID: 8, Symbol: "H3", X: 3, Y: 1, Z: 2, IsBlocked: true, IsMatched: false},
	9:  {ID: 9, Symbol: "H2", X: -4, Y: 0, Z: 3, IsBlocked: false, IsMatched: false},
	10: {ID: 10, Symbol: "H2", X: -2, Y: 0, Z: 1, IsBlocked: true, IsMatched: false},
	11: {ID: 11, Symbol: "H2", X: 2, Y: 0, Z: 1, IsBlocked: true, IsMatched: false},
	12: {ID: 12, Symbol: "H2", X: 4, Y: 0, Z: 3, IsBlocked: false, IsMatched: false},
	13: {ID: 13, Symbol: "L1", X: -3, Y: -1, Z: 2, IsBlocked: true, IsMatched: false},
	14: {ID: 14, Symbol: "L4", X: -1, Y: -1, Z: 0, IsBlocked: true, IsMatched: false},
	15: {ID: 15, Symbol: "L2", X: 1, Y: -1, Z: 0, IsBlocked: true, IsMatched: false},
	16: {ID: 16, Symbol: "H3", X: 3, Y: -1, Z: 2, IsBlocked: true, IsMatched: false},
	17: {ID: 17, Symbol: "L1", X: 0, Y: -2, Z: 1, IsBlocked: true, IsMatched: false},
	18: {ID: 18, Symbol: "H1", X: -1, Y: -3, Z: 2, IsBlocked: true, IsMatched: false},
	19: {ID: 19, Symbol: "H1", X: 1, Y: -3, Z: 2, IsBlocked: true, IsMatched: false},
	20: {ID: 20, Symbol: "L2", X: 0, Y: -4, Z: 3, IsBlocked: false, IsMatched: false},  //
	21: {ID: 21, Symbol: "L2", X: 10, Y: -4, Z: 3, IsBlocked: false, IsMatched: false}, //
}
