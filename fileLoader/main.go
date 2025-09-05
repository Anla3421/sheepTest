package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 從 Excel 讀取關卡佈局
	// 假設工作表(Sheet)的名稱為 "Sheet1"，如果不是請修改
	gameBoard, err := LoadTilesFromExcel("C:\\lt_game\\testZone\\sheep-test\\test1.xlsx", "pic")
	if err != nil {
		fmt.Printf("從Excel讀取關卡失敗: %v\n", err)
		return
	}
	fmt.Printf("成功從 Excel 載入 %d 個磚塊。\n", len(gameBoard))

	InitializeTileDependencies(gameBoard)
	for id, v := range gameBoard {
		fmt.Printf("gameBoard[%d]: %+v\n", id, v)
	}

	// Generate the brick.go file
	err = GenerateBrickFile(gameBoard, "C:\\lt_game\\testZone\\sheep-test\\fileLoader\\brick.go")
	if err != nil {
		fmt.Printf("Failed to generate brick.go: %v\n", err)
	} else {
		fmt.Println("\nSuccessfully generated brick.go")
	}
}

// LoadTilesFromExcel reads an Excel file and converts it into a map of Tile structs.
func LoadTilesFromExcel(filePath string, sheetName string) (map[int]*Tile, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet '%s': %w", sheetName, err)
	}

	var originRow, originCol int = -1, -1

	// First pass: find the "ZERO" cell to establish the origin
	for i, row := range rows {
		for j, cell := range row {
			if strings.ToUpper(cell) == "ZERO" {
				originRow = i
				originCol = j
				break
			}
		}
		if originRow != -1 {
			break
		}
	}

	if originRow == -1 {
		return nil, fmt.Errorf("could not find 'ZERO' cell in sheet '%s'", sheetName)
	}

	tiles := make(map[int]*Tile)
	idCounter := 1

	// Second pass: create tiles based on relative positions and cell content
	for i, row := range rows {
		for j, cell := range row {
			if cell == "" || (i == originRow && j == originCol) {
				continue // Skip empty cells and the origin cell itself
			}

			// Parse content like '2,"L1"'
			parts := strings.Split(cell, ",")
			if len(parts) != 2 {
				// Skip cells that don't match the Z,Symbol format
				continue
			}

			// Trim spaces from the Z part and convert to integer
			z, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				// Skip cells where the Z value is not a valid integer
				continue
			}

			// Trim quotes and spaces from the symbol part
			symbol := strings.Trim(parts[1], ` "`)

			tile := &Tile{
				Symbol:    symbol,
				X:         j - originCol,
				Y:         originRow - i, // Y increases as we go up in the sheet
				Z:         z,
				IsBlocked: false, // This will be calculated later
				IsMatched: false,
			}
			tiles[idCounter] = tile
			idCounter++
		}
	}

	return tiles, nil
}

// InitializeTileDependencies 計算所有磚塊的覆蓋與被覆蓋關係
func InitializeTileDependencies(tiles map[int]*Tile) {
	for _, tileA := range tiles {
		for _, tileB := range tiles {
			if tileA == tileB {
				continue
			}

			// 檢查 tileA 是否在 tileB 的正上一層
			if tileA.Z == tileB.Z+1 {
				// 檢查它們在 XY 平面上是否重疊
				if doTheyOverlap(tileA, tileB) {
					// 建立關係：A 蓋住 B
					tileA.Covers = append(tileA.Covers, tileB)
					// 建立關係：B 被 A 蓋住
					tileB.CoveredBy = append(tileB.CoveredBy, tileA)
				}
			}
		}
	}

	// 建立完圖後，根據 CoveredBy 列表初始化 IsBlocked 狀態
	for _, tile := range tiles {
		if len(tile.CoveredBy) > 0 {
			tile.IsBlocked = true
		}
	}
}

// doTheyOverlap 是一個輔助函數，判斷兩塊牌在 XY 平面是否重疊
// 這裡是一個簡化的範例，假設牌的大小是 2x1 個單位
func doTheyOverlap(a, b *Tile) bool {
	// 假設牌的寬度是 2，高度是 1
	a_minX, a_maxX := a.X-2, a.X+2
	a_minY, a_maxY := a.Y-1, a.Y+1

	b_minX, b_maxX := b.X-2, b.X+2
	b_minY, b_maxY := b.Y-1, b.Y+1

	// 檢查 A 和 B 的矩形區域是否重疊
	if a_maxX <= b_minX || b_maxX <= a_minX {
		return false
	}
	if a_maxY <= b_minY || b_maxY <= a_minY {
		return false
	}
	return true
}

// GenerateBrickFile takes the tile data and generates a brick.go file.
func GenerateBrickFile(tiles map[int]*Tile, outputPath string) error {
	var builder strings.Builder
	builder.WriteString("package main\n\n")
	builder.WriteString("// LvOneTiles is a map of Tiles generated from an Excel file.\n")
	builder.WriteString("var LvOneTiles = map[int]*Tile{\n")

	for id, t := range tiles {
		// Note: Covers and CoveredBy are intentionally left empty.
		// They must be initialized at runtime by InitializeTileDependencies.
		builder.WriteString(fmt.Sprintf(
			"\t%d: {ID: %d, Symbol: \"%s\", X: %d, Y: %d, Z: %d, IsBlocked: %t, IsMatched: %t},\n",
			id, id, t.Symbol, t.X, t.Y, t.Z, t.IsBlocked, t.IsMatched,
		))
	}

	builder.WriteString("}\n")

	// Write the generated content to the file
	err := os.WriteFile(outputPath, []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %%w", outputPath, err)
	}
	return nil
}

// Tile 代表一個麻將牌 (更新版)

type Tile struct {
	ID        int
	Symbol    string
	X, Y, Z   int
	IsBlocked bool
	IsMatched bool

	// Covers: 這塊牌直接蓋住了哪些牌 (它的孩子們)
	Covers []*Tile
	// CoveredBy: 這塊牌被哪些牌直接蓋住 (它的父母們)
	CoveredBy []*Tile
}
