package main

import (
	"fmt"
)

// DataObject 假設這是你說的 "object"
type DataObject struct {
	ID   int
	Name string
}

// SlotContainer 建立一個容器 Struct 來管理 slots 和限制
type SlotContainer struct {
	slots []any // 使用 any (interface{}) 來存放不同類型的資料
	limit int   // 動態限制
}

// NewSlotContainer "建構函式"，用來建立一個新的容器
func NewSlotContainer(limit int) (*SlotContainer, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit 必須大於 0")
	}
	return &SlotContainer{
		slots: make([]any, 0, limit), // 預先分配容量以提高效能
		limit: limit,
	}, nil
}

// Add 新增項目到容器的方法
func (c *SlotContainer) Add(item any) error {
	// 檢查是否已滿
	if len(c.slots) >= c.limit {
		return fmt.Errorf("空間已滿 (限制 %d)，無法新增項目", c.limit)
	}
	c.slots = append(c.slots, item)
	return nil
}

// SetLimit 用於動態調整大小限制
func (c *SlotContainer) SetLimit(newLimit int) error {
	if newLimit <= 0 {
		return fmt.Errorf("新的 limit 必須大於 0")
	}

	fmt.Printf("\n>>>>> 正在將限制從 %d 調整為 %d <<<<<", c.limit, newLimit)
	c.limit = newLimit

	// 核心邏輯：如果當前項目數量超過了新的限制，就截斷 slice
	if len(c.slots) > newLimit {
		fmt.Printf("!! 當前數量 (%d) 超過新限制 (%d)，將刪除多餘項目。\n", len(c.slots), newLimit)
		c.slots = c.slots[0:newLimit]
	}

	return nil
}

// Display 顯示容器內的項目 (使用 Type Switch)
func (c *SlotContainer) Display() {
	fmt.Printf("--- 容器狀態 (限制: %d, 當前: %d) ---", c.limit, len(c.slots))
	for i, item := range c.slots {
		switch v := item.(type) {
		case string:
			fmt.Printf("Slot %d: [字串] %s\n", i+1, v)
		case DataObject:
			fmt.Printf("Slot %d: [物件] ID=%d, Name=%s\n", i+1, v.ID, v.Name)
		default:
			fmt.Printf("Slot %d: [未知類型]\n", i+1)
		}
	}
	fmt.Println("------------------------------------")
}

func main() {
	// 1. 建立一個限制為 5 的容器
	container, _ := NewSlotContainer(5)
	fmt.Println("初始化容器，限制為 5。")

	// 2. 塞滿它
	container.Add("一萬")
	container.Add("兩萬")
	container.Add("三萬")
	container.Add("四萬")
	container.Add(DataObject{ID: 99, Name: "五萬物件"})
	container.Display()

	// 3. 現在，將限制從 5 縮小到 3
	container.SetLimit(3)
	fmt.Println("調整限制完成。")
	container.Display() // 檢查結果，應該只剩下前 3 個項目

	// 4. 嘗試新增一個項目，應該會失敗，因為空間已滿 (3/3)
	fmt.Println("\n嘗試新增 '新項目' (預期會失敗)")
	err := container.Add("新項目")
	if err != nil {
		fmt.Println("新增失敗:", err)
	}

	// 5. 現在，將限制從 3 擴大到 4
	container.SetLimit(4)
	fmt.Println("調整限制完成。")
	container.Display()

	// 6. 再次嘗試新增，這次應該會成功
	fmt.Println("\n嘗試新增 '新項目' (預期會成功)")
	err = container.Add("新項目")
	if err != nil {
		fmt.Println("新增失敗:", err)
	}
	container.Display() // 檢查結果，新項目已加入
}
