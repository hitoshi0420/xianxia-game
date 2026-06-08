package game

import "fmt"

// 商店物品列表
type ShopItem struct {
	Item  *Item
	Price int
}

// 获取商店物品（根据玩家境界动态生成）
func GetShopItems(player *Player) []ShopItem {
	items := []ShopItem{
		{Item: &Item{Name: "回春丹", Type: ItemPill, Desc: "恢复80点生命", HP: 80, Amount: 1}, Price: 15},
		{Item: &Item{Name: "大还丹", Type: ItemPill, Desc: "恢复200点生命", HP: 200, Amount: 1}, Price: 35},
		{Item: &Item{Name: "聚气散", Type: ItemPill, Desc: "增加150点修为", Cult: 150, Amount: 1}, Price: 25},
		{Item: &Item{Name: "凝神丹", Type: ItemPill, Desc: "增加300点修为", Cult: 300, Amount: 1}, Price: 55},
	}

	// 根据境界解锁装备
	if int(player.Realm) >= 2 {
		items = append(items, ShopItem{Item: &Item{Name: "青锋剑", Type: ItemWeapon, Desc: "利剑 +15攻", Atk: 15, Amount: 1}, Price: 120})
		items = append(items, ShopItem{Item: &Item{Name: "铁布衫", Type: ItemArmor, Desc: "护甲 +10防", Def: 10, Amount: 1}, Price: 100})
	}
	if int(player.Realm) >= 4 {
		items = append(items, ShopItem{Item: &Item{Name: "玄铁重剑", Type: ItemWeapon, Desc: "重剑 +30攻", Atk: 30, Amount: 1}, Price: 350})
		items = append(items, ShopItem{Item: &Item{Name: "金丝甲", Type: ItemArmor, Desc: "软甲 +25防", Def: 25, Amount: 1}, Price: 300})
	}
	if int(player.Realm) >= 6 {
		items = append(items, ShopItem{Item: &Item{Name: "飞虹剑", Type: ItemWeapon, Desc: "神剑 +60攻", Atk: 60, Amount: 1}, Price: 900})
		items = append(items, ShopItem{Item: &Item{Name: "天蚕衣", Type: ItemArmor, Desc: "宝衣 +50防", Def: 50, Amount: 1}, Price: 700})
	}

	return items
}

// 购买物品
func BuyItem(player *Player, shopItems []ShopItem, index int) (string, bool) {
	if index < 0 || index >= len(shopItems) {
		return "无效的选择", false
	}

	si := shopItems[index]
	if player.Stones < si.Price {
		return fmt.Sprintf("灵石不足！需要 %d 灵石，你只有 %d", si.Price, player.Stones), false
	}

	player.Stones -= si.Price
	player.AddItem(si.Item)
	player.addLog(fmt.Sprintf("购买了 %s，花费 %d 灵石", si.Item.Name, si.Price))
	return fmt.Sprintf("购买成功！获得【%s】，消耗 %d 灵石", si.Item.Name, si.Price), true
}
