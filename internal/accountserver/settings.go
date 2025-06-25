package accountserver

type AccountServerSettings struct {
	ClassWiseStarterInfo ClassWiseStarterInfo `json:"class_wise_starter_info"`
}

type StarterInfo struct {
	Experience     uint32          `json:"experience"`
	Woonz          uint32          `json:"woonz"`
	Level          byte            `json:"level"`
	WearItems      []Item          `json:"wear_items"`
	InventoryItems []InventoryItem `json:"inventory_items"`
}

type Item struct {
	ItemCode   uint32 `json:"item_code"`
	ItemOption uint32 `json:"item_option"`
}

type InventoryItem struct {
	ItemCode   uint32 `json:"item_code"`
	ItemOption uint32 `json:"item_option"`
	Slot       byte   `json:"slot"`
}

type ClassWiseStarterInfo struct {
	Warrior StarterInfo `json:"warrior"`
	HK      StarterInfo `json:"hk"`
	Mage    StarterInfo `json:"mage"`
	Archer  StarterInfo `json:"archer"`
}
