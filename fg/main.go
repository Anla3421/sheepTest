package main

import (
	"fmt"
	"log"
	"math/rand"
	"sheep-test/fg/sheepArea"
	"sort"
)

const totalSpins = 100000

var limit = 7
var luckySlot = &[]string{}

func main() {
	fmt.Printf("init luckySlot: %+v\n", luckySlot)
	gameBoard := SheepAreaInit()
	PrintBoard(gameBoard, "init")

	for i := 0; i < totalSpins; i++ {
		// fmt.Printf("!!!新的一局開始了!!! init luckySlot: %+v\n", luckySlot)
		// luckySlot := []string{"一萬", "兩萬", "五萬"}
		// needToCheck := SpinResult()
		needToCheck := NewSpinResult()
		// fmt.Println("!!!這次的連線的符號!!!", needToCheck.Symbol)
		Compare(gameBoard, luckySlot, needToCheck.Symbol)
		// 現有盤面確認，如果盤面清空了(len()=0)，獲得JP (這邊先用迴圈中止處理)
		if JackpotCheck(gameBoard) {
			log.Println("!!!!!JP", i)
			PrintBoard(gameBoard, "finish")
			fmt.Printf("!!!!!JP finish luckySlot: %+v\n", luckySlot)
			break
		}
	}
	log.Println("!!!!!no JP!!!!")
	// PrintBoard(gameBoard, "finish")
	// fmt.Printf("finish luckySlot: %+v\n", luckySlot)
}

// Add 新增項目到容器的方法
func Add(luckySlot *[]string, item string) string {
	// 檢查是否已滿
	if len((*luckySlot)) >= limit {
		return fmt.Sprintf("空間已滿 (限制 %d)，無法新增項目\n", limit)
	}
	(*luckySlot) = append((*luckySlot), item)
	return "新增至 lucky shot 成功"
}

// SetLimit 用於動態調整大小限制
func SetLimit(luckySlot []string, newLimit int) error {
	if newLimit <= 0 {
		return fmt.Errorf("新的 limit 必須大於 0")
	}

	fmt.Printf("\n>>>>> 正在將限制從 %d 調整為 %d <<<<<", limit, newLimit)
	limit = newLimit

	// 核心邏輯：如果當前項目數量超過了新的限制，就截斷 slice
	if len(luckySlot) > newLimit {
		fmt.Printf("!! 當前數量 (%d) 超過新限制 (%d)，將刪除多餘項目。\n", len(luckySlot), newLimit)
		luckySlot = luckySlot[0:newLimit]
	}

	return nil
}

// func SheepAreaInit() []*Tile {
func SheepAreaInit() map[int]*Tile {
	tilesByID := map[int]*Tile{}
	for k, v := range sheepArea.LvOneTiles {
		tilesByID[k] = &Tile{
			ID:     v.ID,
			Symbol: v.Symbol,
			X:      v.X,
			Y:      v.Y,
			Z:      v.Z,
		}
	}

	tilesByZ := map[int][]*Tile{}
	for _, tile := range tilesByID {
		tilesByZ[tile.Z] = append(tilesByZ[tile.Z], tile)
	}

	for z, tileAtZ := range tilesByZ {
		if tilesAtZMinus1, ok := tilesByZ[z-1]; ok {
			for _, tileA := range tileAtZ {
				for _, tileB := range tilesAtZMinus1 {
					if doTheyOverlap(tileA, tileB) {
						// 建立關係：A 蓋住 B
						tileA.Covers = append(tileA.Covers, tileB)
						// 建立關係：B 被 A 蓋住
						tileB.CoveredBy = append(tileB.CoveredBy, tileA)
					}
				}
			}
		}
	}

	for _, tile := range tilesByID {
		if len(tile.CoveredBy) > 0 {
			tile.IsBlocked = true
		}
	}

	// return tiles
	return tilesByID
}

// JPCheck
func JackpotCheck(gameBoard map[int]*Tile) bool {
	isBoardEmpty := false
	tileMatched := 0
	for _, v := range gameBoard {
		if v.IsMatched {
			tileMatched++
		}
	}
	if len(gameBoard) == tileMatched {
		isBoardEmpty = true
	}
	return isBoardEmpty
}

// 盤面任何有變動後，都要檢查一次現在的盤面及 lucky shot 是不是會因為打開下層能有配對消除
// 應該是用 map[string]int 來計數，實作
// 前端演繹?
func Intermission(gameBoard map[int]*Tile, luckySlot *[]string) {
	isEreased := false
	symbolCount := map[string]int{}
	markToErease := []*Tile{}
	for _, tile := range gameBoard {
		// 牌沒被擋，也存在
		if !tile.IsBlocked && !tile.IsMatched {
			symbolCount[tile.Symbol]++
			markToErease = append(markToErease, tile)
		}
	}
	for symbol, count := range symbolCount {
		ereaseNow := []*Tile{}
		indexNeedToCut := []int{}
		switch count {
		case 2: // luckyshot 1 + 盤面 2
			// log.Println("luckyshot 1 + 盤面 2")
			for i, s := range *luckySlot {
				if s == symbol {
					// 記錄起來留到最後處理
					indexNeedToCut = append(indexNeedToCut, i)
					// 這邊只要1個，找到 1 個符合的符號就跳出
					break
				}
			}
			// 如果 luckyshot 數量不足 1 就不執行
			if len(indexNeedToCut) < 1 {
				break
			}
			for _, t := range markToErease {
				// 取出符合 t.Symbol 的麻將 2 個，將他們消除
				i := 0
				if t.Symbol == symbol && i < 2 {
					ereaseNow = append(ereaseNow, t)
					i++
				}
				for _, t := range ereaseNow {
					t.IsMatched = true
					fmt.Printf("Set IsMatched: %+v\n", t)
					// 當一個牌被消除 (IsMatched = true)，我們需要檢查它所覆蓋的牌
					for _, coveredTile := range t.Covers {
						// 檢查這張被覆蓋的牌是否還有其他未消除的牌壓著它
						isStillBlocked := false
						for _, parentTile := range coveredTile.CoveredBy {
							if !parentTile.IsMatched {
								// 只要還有任何一個壓著它的牌是未消除的，它就仍然被阻擋
								isStillBlocked = true
								break
							}
						}
						// 如果所有壓著它的牌都已經被消除了，就更新它的狀態
						if !isStillBlocked {
							coveredTile.IsBlocked = false
							fmt.Printf("Unblocked Tile: %s (ID: %d)\n", coveredTile.Symbol, coveredTile.ID)
						}
					}
				}
			}
			// log.Println("luckyshot1 + 盤面2 isEreased", isEreased)
			isEreased = true
		case 1: // luckyshot2 + 盤面1
			// log.Println("luckyshot2 + 盤面1")
			for i, s := range *luckySlot {
				if s == symbol {
					// 記錄起來留到最後處理
					indexNeedToCut = append(indexNeedToCut, i)
					// 這邊只要 2 個，找到 2 個符合的符號就跳出
					if len(indexNeedToCut) == 2 {
						break
					}
				}
			}
			// 如果 luckyshot 數量不足 1 就不執行
			if len(indexNeedToCut) < 2 {
				break
			}
			for _, t := range markToErease {
				// 取出符合 t.Symbol 的麻將 2 個，將他們消除
				i := 0
				if t.Symbol == symbol && i < 1 {
					ereaseNow = append(ereaseNow, t)
					i++
				}
				for _, t := range ereaseNow {
					t.IsMatched = true
					fmt.Printf("Set IsMatched: %+v\n", t)
					// 當一個牌被消除 (IsMatched = true)，我們需要檢查它所覆蓋的牌
					for _, coveredTile := range t.Covers {
						// 檢查這張被覆蓋的牌是否還有其他未消除的牌壓著它
						isStillBlocked := false
						for _, parentTile := range coveredTile.CoveredBy {
							if !parentTile.IsMatched {
								// 只要還有任何一個壓著它的牌是未消除的，它就仍然被阻擋
								isStillBlocked = true
								break
							}
						}
						// 如果所有壓著它的牌都已經被消除了，就更新它的狀態
						if !isStillBlocked {
							coveredTile.IsBlocked = false
							fmt.Printf("Unblocked Tile: %s (ID: %d)\n", coveredTile.Symbol, coveredTile.ID)
						}
					}
				}
			}
			// log.Println("luckyshot2 + 盤面1 isEreased", isEreased)
			isEreased = true
		case 0: // 盤面0
		default: // 盤面 >=3
			// log.Println("盤面3")
			for _, t := range markToErease {
				// 取出符合 t.Symbol 的麻將 3 個，將他們消除
				i := 0
				if t.Symbol == symbol && i < 3 {
					ereaseNow = append(ereaseNow, t)
					i++
				}
			}
			for _, v := range ereaseNow {
				fmt.Printf("ereaseNow: %+v\n", v)
			}
			for _, t := range ereaseNow {
				t.IsMatched = true
				fmt.Printf("Set IsMatched: %+v\n", t)

				// 當一個牌被消除 (IsMatched = true)，我們需要檢查它所覆蓋的牌
				for _, coveredTile := range t.Covers {
					// 檢查這張被覆蓋的牌是否還有其他未消除的牌壓著它
					isStillBlocked := false
					for _, parentTile := range coveredTile.CoveredBy {
						if !parentTile.IsMatched {
							// 只要還有任何一個壓著它的牌是未消除的，它就仍然被阻擋
							isStillBlocked = true
							break
						}
					}

					// 如果所有壓著它的牌都已經被消除了，就更新它的狀態
					if !isStillBlocked {
						coveredTile.IsBlocked = false
						fmt.Printf("Unblocked Tile: %s (ID: %d)\n", coveredTile.Symbol, coveredTile.ID)
					}
				}
			}
			// log.Println("盤面3", isEreased)
			isEreased = true
		}

		// 執行移除 luckySlot 需要砍掉的符號
		cutMap := map[int]int{}
		for _, index := range indexNeedToCut {
			cutMap[index]++
		}
		n := 0
		tempLuckySlot := *luckySlot
		for i, v := range tempLuckySlot {
			if _, exist := cutMap[i]; exist {
				tempLuckySlot[n] = v
				n++
			}
		}
		*luckySlot = tempLuckySlot[:n]
	}

	// 如果有三消成功則進遞迴
	if isEreased {
		Intermission(gameBoard, luckySlot)
	}
}

// 比對：先看盤面跟 symbol 有沒有一樣的(盤面至少有一個)
// yes -> 湊的了3個?
//
//	yes -> 消除
//	no -> 檢查 luckySlot 有沒有，有就消沒有加入(如果沒滿)
//
// no -> 加入 luckySlot (如果沒滿)
func Compare(gameBoard map[int]*Tile, luckySlot *[]string, symbol string) {
	count := 0
	markToErease := []*Tile{}
	isEreased := false
	for _, tile := range gameBoard {
		// 符號有符合，牌沒被擋，也存在
		if tile.Symbol == symbol && !tile.IsBlocked && !tile.IsMatched {
			count++
			markToErease = append(markToErease, tile)
		}
		// 輪盤 1 + 盤面上有 2 個可消除的符號
		if count > 1 {
			for _, t := range markToErease {
				t.IsMatched = true
				fmt.Printf("Set IsMatched: %+v\n", t)

				// 當一個牌被消除 (IsMatched = true)，我們需要檢查它所覆蓋的牌
				for _, coveredTile := range t.Covers {
					// 檢查這張被覆蓋的牌是否還有其他未消除的牌壓著它
					isStillBlocked := false
					for _, parentTile := range coveredTile.CoveredBy {
						if !parentTile.IsMatched {
							// 只要還有任何一個壓著它的牌是未消除的，它就仍然被阻擋
							isStillBlocked = true
							break
						}
					}

					// 如果所有壓著它的牌都已經被消除了，就更新它的狀態
					if !isStillBlocked {
						coveredTile.IsBlocked = false
						fmt.Printf("Unblocked Tile: %s (ID: %d)\n", coveredTile.Symbol, coveredTile.ID)
					}
				}
			}
			isEreased = true
			break
		}
	}

	for _, v := range markToErease {
		fmt.Printf("markToErease: %+v\n", v)
	}

	// 盤面沒有消除東西，檢查luckySlot裡面
	// 是否有項目可以跟上面湊成3各
	// 把東西加到裡面(有空位的話)
	if !isEreased {
		// newLuckyshot := *[]string{}
		indexNeedToCut := []int{}
		for i, s := range *luckySlot {
			if s == symbol {
				count++
			}
			// 輪盤 1 + 盤面上有 0~1 個可消除的符號 + lucky shot 1~2 -> 待討論，lucky shot 可以湊 3 個消掉嗎？
			if count > 1 {
				for _, t := range markToErease {
					t.IsMatched = true
					// 當一個牌被消除 (IsMatched = true)，我們需要檢查它所覆蓋的牌
					for _, coveredTile := range t.Covers {
						// 檢查這張被覆蓋的牌是否還有其他未消除的牌壓著它
						isStillBlocked := false
						for _, parentTile := range coveredTile.CoveredBy {
							if !parentTile.IsMatched {
								// 只要還有任何一個壓著它的牌是未消除的，它就仍然被阻擋
								isStillBlocked = true
								break
							}
						}

						// 如果所有壓著它的牌都已經被消除了，就更新它的狀態
						if !isStillBlocked {
							coveredTile.IsBlocked = false
							fmt.Printf("Unblocked Tile: %s (ID: %d)\n", coveredTile.Symbol, coveredTile.ID)
						}
					}
				}
				// 紀錄 luckySlot 中所有需要砍掉的符號的 index
				indexNeedToCut = append(indexNeedToCut, i)
				isEreased = true
				break
			}
		}

		// 執行移除 luckySlot 需要砍掉的符號
		for i := range *luckySlot {
			for _, v := range indexNeedToCut {
				if i == v {
					(*luckySlot) = append((*luckySlot)[:i], (*luckySlot)[i+1:]...)
				}
			}
		}
	}
	fmt.Printf("compre finish, luckySlot: %+v\n", luckySlot)
	fmt.Printf("compre finish, symbol: %+v\n", symbol)
	// 無法消除任何項目，把符號加入 lucky shot
	if !isEreased {
		// fmt.Printf("before add, luckySlot: %+v\n", luckySlot)
		msg := Add(luckySlot, symbol)
		fmt.Println("msg", msg)
		// fmt.Printf("alter luckySlot: %+v\n", luckySlot)
	}

	Intermission(gameBoard, luckySlot)
}

func NewSpinResult() *Tile {
	tiles := []*Tile{}
	for _, v := range sheepArea.LvOneTiles {
		tempTiles := &Tile{
			// ID:     v.ID,
			Symbol: v.Symbol,
			X:      v.X,
			Y:      v.Y,
			Z:      v.Z,
		}
		tiles = append(tiles, tempTiles)
	}
	// index 0, {ID: 1, Symbol: "L2", X: 0, Y: 4, Z: 3, IsBlocked: false, IsMatched: false},
	// index 11, {ID: 12, Symbol: "H2", X: 4, Y: 0, Z: 3, IsBlocked: false, IsMatched: false},
	result := rand.Intn(20) // 0 ~ 19
	// result = 5
	return tiles[result]
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

// Tile 代表一個麻將牌 (更新版)
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

func PrintBoard(gameBoard map[int]*Tile, remark string) {
	keys := make([]int, 0, len(gameBoard))
	for k := range gameBoard {
		keys = append(keys, k)
	}
	// 接著對 key 進行排序
	sort.Ints(keys)

	// 最後按照排好的順序來迭代 map
	for _, k := range keys {
		v := gameBoard[k]
		fmt.Printf("%v gameBoard (key: %d): %+v\n", remark, k, v)
	}
}
