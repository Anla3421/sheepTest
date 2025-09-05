package main

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	rows       = 3
	columns    = 3
	totalSpins = 1000000 // 1,000,000
	betPerSpin = 5
)

// Define paylines (row, column coordinates)
var paylines = [][]struct{ r, c int }{
	// Horizontal lines
	{{0, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},
	// Diagonal lines
	{{0, 0}, {1, 1}, {2, 2}},
	{{2, 0}, {1, 1}, {0, 2}},
}

// PayTable，WW 不自付
var payTable = map[string]map[int]int{
	"H1": {3: 100},
	"H2": {3: 50},
	"H3": {3: 20},
	"L1": {3: 10},
	"L2": {3: 5},
	"L3": {3: 3},
	"L4": {3: 2},
}

// 輪帶資料
var reels = [][]string{
	{"H1", "H1", "H1", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L4", "L4", "L4", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3"},
	{"H1", "H1", "H1", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L4", "L4", "L4", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3"},
	{"H1", "H1", "H1", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L4", "L4", "L4", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3"},
}

var winSymbolCount = map[string]int{
	"H1": 0,
	"H2": 0,
	"H3": 0,
	"L1": 0,
	"L2": 0,
	"L3": 0,
	"L4": 0,
}

func main() {
	// Initialize random seed (important for rand.Intn)
	// For Go 1.20+, consider using rand.New(rand.NewSource(time.Now().UnixNano()))
	// For older versions, uncomment: rand.Seed(time.Now().UnixNano())

	totalWin := 0
	for i := 0; i < totalSpins; i++ {
		currentGrid := spin()
		// currentGrid = [][]string{
		// 	{"L1", "L1", "L1"},
		// 	{"H1", "H2", "H1"},
		// 	{"L1", "L1", "L1"},
		// }
		isGolden := goldenSymbolPick()
		winAmount, _, changeToWild := calculateWinBG(currentGrid, isGolden, payTable)
		// log.Printf("BG Win: %d\n", winAmount)
		if len(changeToWild) > 0 {
			currentGrid := spin()
			winAmount += calculateWinFG(currentGrid, changeToWild, payTable)
			// log.Printf("FG Win: %d\n", winAmount)
		}
		// log.Printf("Round Win: %d\n", winAmount)
		totalWin += winAmount
	}
	resultTable(totalWin)
}

// spin 產生一個盤面
func spin() [][]string {
	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, columns)
	}

	for c := 0; c < columns; c++ {
		reel := reels[c]
		startPos := rand.Intn(len(reel))
		for r := 0; r < rows; r++ {
			pos := (startPos + r) % len(reel)
			grid[r][c] = reel[pos]
		}
	}
	return grid
}

func goldenSymbolPick() (isGolden []bool) {
	// 正中間的 column 有機會出現黃金版本的符號(機率看權重，預設10%)
	i := 0
	for i < 3 {
		g := false
		if rand.Intn(100) < 10 {
			g = true
		}
		isGolden = append(isGolden, g)
		i++
	}
	// log.Println("isGolden", isGolden)
	return isGolden
}

// calculateWinBG 計算贏分並記錄各項贏分
func calculateWinBG(grid [][]string, isGolden []bool, payoutsBySymbol map[string]map[int]int) (int, map[bool]int, map[int]bool) {
	totalWin := 0
	gotFreeSpin := map[bool]int{}
	changeToWild := map[int]bool{}
	for _, line := range paylines {
		s1 := grid[line[0].r][line[0].c]
		s2 := grid[line[1].r][line[1].c]
		s3 := grid[line[2].r][line[2].c]

		var targetSymbol string
		// Find the first non-wild symbol as the target for the win
		targetSymbol = s1
		// Check if all symbols on the line match the target or are wild
		isWin := (s2 == targetSymbol) && (s3 == targetSymbol)
		if isWin {
			if payoutsForSymbol, ok := payoutsBySymbol[targetSymbol]; ok {
				if payout, ok := payoutsForSymbol[3]; ok { // Always count 3 for a 3-symbol line win
					totalWin += payout
					winSymbolCount[targetSymbol]++
				}
			}
		}

		if isWin && isGolden[line[1].r] {
			// 確認剛剛中獎的是有黃金版本符號則記錄起來
			gotFreeSpin[true]++ // 應該不需要了? 待確認後刪除
			// 記錄進FG要變百搭鎖住的位置
			changeToWild[line[1].r] = true
		}
	}
	return totalWin, gotFreeSpin, changeToWild
}

// calculateWinFG 計算贏分並記錄各項贏分
// 當黃金版本的符號有連線得獎時，那個格子會在下一輪轉變成百搭，並執行一次free game (本次執行時不會有黃金版本的符號即這局權重是0)
// 百搭符號只會維持一局
func calculateWinFG(grid [][]string, changeToWild map[int]bool, payoutsBySymbol map[string]map[int]int) int {
	totalWin := 0
	for k := range changeToWild {
		grid[k][1] = "WW"
	}
	// log.Println("FG grid", grid)

	for _, line := range paylines {
		s1 := grid[line[0].r][line[0].c]
		s2 := grid[line[1].r][line[1].c]
		s3 := grid[line[2].r][line[2].c]

		var targetSymbol string
		// Find the first non-wild symbol as the target for the win
		if s1 != "WW" {
			targetSymbol = s1
		} else if s2 != "WW" {
			targetSymbol = s2
		} else if s3 != "WW" {
			targetSymbol = s3
		} else {
			// All symbols are WW, and WW doesn't pay by itself
			continue // No win for a line of all WWs
		}

		// Check if all symbols on the line match the target or are wild
		isWin := (s1 == targetSymbol || s1 == "WW") &&
			(s2 == targetSymbol || s2 == "WW") &&
			(s3 == targetSymbol || s3 == "WW")

		if isWin {
			if payoutsForSymbol, ok := payoutsBySymbol[targetSymbol]; ok {
				if payout, ok := payoutsForSymbol[3]; ok { // Always count 3 for a 3-symbol line win
					totalWin += payout
					winSymbolCount[targetSymbol]++
				}
			}
		}
	}

	return totalWin
}

func resultTable(totalWin int) {
	totalBet := float64(betPerSpin * totalSpins)
	rtp := float64(totalWin) / totalBet
	log.Printf("RTP: %v\n", rtp)
	fmt.Println(`---模擬結果---`)
	fmt.Printf("總模擬次數: %v\n", totalSpins)
	fmt.Printf("%-5s %-10s %-10s\n", "總下注金額", "總贏取金額", "總體RTP")
	fmt.Printf("%-5.0f %-15v %-10.8f\n", totalBet, totalWin, rtp)
	fmt.Printf("%-5s %-15s %-10s\n", "SYM", "TotalPayout", "RTP")
	for k, v := range winSymbolCount {
		for k2, v2 := range payTable {
			if k == k2 {
				fmt.Printf("%-5s %-15v %-10.8f\n", k, v*v2[3], float64(v*v2[3])/totalBet)
			}
		}
	}
}
