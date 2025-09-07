package main

import (
	"fmt"
	"log"
	"math/rand"
	sheepfunc "sheep-test/fg/sheepFunc"
)

const (
	rows       = 3
	columns    = 3
	totalSpins = 1000000 // 1,000,000
	betPerSpin = 5
)

var (
	// Define paylines (row, column coordinates)
	paylines = [][]struct{ r, c int }{
		// Horizontal lines
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},
		// Diagonal lines
		{{0, 0}, {1, 1}, {2, 2}},
		{{2, 0}, {1, 1}, {0, 2}},
	}
	// PayTable，WW 不自付
	payTable = map[string]map[int]int{
		"H1": {3: 100},
		"H2": {3: 50},
		"H3": {3: 20},
		"L1": {3: 10},
		"L2": {3: 5},
		"L3": {3: 3},
		"L4": {3: 2},
	}
	// 輪帶資料
	reels = [][]string{
		{"H1", "H1", "H1", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L4", "L4", "L4", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3"},
		{"H1", "H1", "H1", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L4", "L4", "L4", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3"},
		{"H1", "H1", "H1", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L4", "L4", "L4", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3", "L1", "L1", "L1", "H3", "H3", "H3", "H2", "H2", "H2", "L2", "L2", "L2", "L3", "L3", "L3"},
	}
	winSymbolCount = map[string]int{
		"H1": 0,
		"H2": 0,
		"H3": 0,
		"L1": 0,
		"L2": 0,
		"L3": 0,
		"L4": 0,
	}
	luckySlot           = &[]string{}
	totalSheepWinForRtp = 0
)

func main() {
	// Initialize random seed (important for rand.Intn)
	// For Go 1.20+, consider using rand.New(rand.NewSource(time.Now().UnixNano()))
	// For older versions, uncomment: rand.Seed(time.Now().UnixNano())
	sheepGameBoard := sheepfunc.SheepAreaInit()
	totalWin := 0
	for i := 0; i < totalSpins; i++ {
		currentGrid := spin()
		// currentGrid = [][]string{
		// 	{"L1", "L1", "L1"},
		// 	{"H1", "H2", "H1"},
		// 	{"L1", "L1", "L1"},
		// }
		isGolden := goldenSymbolPick()
		bgWinAmount, bgWinSymbols, changeToWild := calculateWinBG(currentGrid, isGolden, payTable)
		// 對每個 BG 獲勝符號執行進羊區邏輯
		for _, winSymbol := range bgWinSymbols {
			sheepWin, newGameBoard := sheepfunc.SheepTrigger(winSymbol, sheepGameBoard, luckySlot)
			totalSheepWinForRtp += sheepWin
			sheepGameBoard = newGameBoard
		}

		// log.Printf("BG Win: %d\n", winAmount)
		if len(changeToWild) > 0 {
			currentGrid := spin()
			fgWinAmount, fgWinSymbols := calculateWinFG(currentGrid, changeToWild, payTable)
			totalWin += fgWinAmount
			// log.Printf("FG Win: %d\n", winAmount)
			// fg 獲勝近羊區邏輯
			for _, winSymbol := range fgWinSymbols {
				sheepWin, newGameBoard := sheepfunc.SheepTrigger(winSymbol, sheepGameBoard, luckySlot)
				totalSheepWinForRtp += sheepWin
				sheepGameBoard = newGameBoard
			}
			// sheepWin := sheepfunc.SheepTrigger(Symbol, sheepGameBoard, luckySlot)
		}

		// log.Printf("Round Win: %d\n", winAmount)
		totalWin += bgWinAmount
	}
	totalWin += totalSheepWinForRtp
	// sheepfunc.PrintBoard(sheepGameBoard, "main final")
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
func calculateWinBG(grid [][]string, isGolden []bool, payoutsBySymbol map[string]map[int]int) (int, []string, map[int]bool) {
	totalWin := 0
	winSymbol := []string{}
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
					winSymbolCount[targetSymbol]++ // rtp 計分用
				}
			}
		}

		if isWin && isGolden[line[1].r] {
			// 確認剛剛中獎的是有黃金版本符號則記錄起來
			// 記錄進FG要變百搭鎖住的位置
			changeToWild[line[1].r] = true
			winSymbol = append(winSymbol, targetSymbol)
		}
	}
	return totalWin, winSymbol, changeToWild
}

// calculateWinFG 計算贏分並記錄各項贏分
// 當黃金版本的符號有連線得獎時，那個格子會在下一輪轉變成百搭，並執行一次free game (本次執行時不會有黃金版本的符號即這局權重是0)
// 百搭符號只會維持一局
func calculateWinFG(grid [][]string, changeToWild map[int]bool, payoutsBySymbol map[string]map[int]int) (int, []string) {
	totalWin := 0
	winSymbol := []string{}
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
					winSymbolCount[targetSymbol]++ // rtp 計分用
					winSymbol = append(winSymbol, targetSymbol)
				}
			}
		}
	}

	return totalWin, winSymbol
}

func resultTable(totalWin int) {
	totalBet := float64(betPerSpin * totalSpins)
	rtp := float64(totalWin) / totalBet
	checkSumAccuracyRtp := 0.0
	fmt.Println(`---模擬結果---`)
	log.Printf("總模擬次數: %v\n", totalSpins)
	fmt.Printf("%-5s %-10s %-10s\n", "總下注金額", "總贏取金額", "總體RTP")
	fmt.Printf("%-5.0f %-15v %-10.8f\n", totalBet, totalWin, rtp)
	fmt.Printf("%-5s %-15s %-10s\n", "SYM", "TotalPayout", "RTP")
	for k, v := range winSymbolCount {
		for k2, v2 := range payTable {
			if k == k2 {
				partialRtp := float64(v*v2[3]) / totalBet
				fmt.Printf("%-5s %-15v %-10.8f\n", k, v*v2[3], float64(v*v2[3])/totalBet)
				checkSumAccuracyRtp += partialRtp
			}
		}
	}
	fmt.Printf("%-5s %-15v %-10.8f\n", "sheepJP", totalSheepWinForRtp, float64(totalSheepWinForRtp)/totalBet)
	partialRtp := float64(totalSheepWinForRtp) / totalBet
	checkSumAccuracyRtp += partialRtp
	log.Printf("checkSumAccuracyRtp: %-10.8f\n", checkSumAccuracyRtp)
}

// func PrintBoard(gameBoard map[int]*Tile, remark string) {
// 	keys := make([]int, 0, len(gameBoard))
// 	for k := range gameBoard {
// 		keys = append(keys, k)
// 	}
// 	// 接著對 key 進行排序
// 	sort.Ints(keys)

// 	// 最後按照排好的順序來迭代 map
// 	for _, k := range keys {
// 		v := gameBoard[k]
// 		fmt.Printf("%v gameBoard (key: %d): %+v\n", remark, k, v)
// 	}
// }
