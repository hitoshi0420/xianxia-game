package game

import (
	"fmt"
	"math/rand"
)

// 玩家
type Player struct {
	Name      string
	Realm     Realm  // 当前境界
	Layer     int    // 当前层 (1-9)
	HP        int
	MaxHP     int
	Atk       int
	Def       int
	Cult      int // 当前修为值
	MaxCult   int // 突破所需修为
	Stones    int // 灵石
	Inventory []*Item
	Log       []string // 最近的游戏日志
	TotalKills int
}

// 创建新角色
func NewPlayer(name string) *Player {
	p := &Player{
		Name:    name,
		Realm:   RealmMortal,
		Layer:   1,
		HP:      100,
		MaxHP:   100,
		Atk:     10,
		Def:     5,
		Cult:    0,
		MaxCult: 100,
		Stones:  50,
	}

	// 初始物品
	p.Inventory = append(p.Inventory, &Item{
		Name: "回春丹", Type: ItemPill, Desc: "恢复50点生命", Value: 10, HP: 50, Amount: 3,
	})
	p.Inventory = append(p.Inventory, &Item{
		Name: "聚气散", Type: ItemPill, Desc: "增加100点修为", Value: 20, Cult: 100, Amount: 2,
	})

	return p
}

// 修炼
func (p *Player) Cultivate() string {
	gain := p.cultivationGain()
	p.Cult += gain

	msg := fmt.Sprintf("修炼中……获得 %d 点修为", gain)
	p.addLog(msg)

	if p.Cult >= p.MaxCult {
		p.Cult = p.MaxCult
		msg += "\n修为已满！可以尝试突破了！"
	}

	return msg
}

// 每次修炼获得的修为
func (p *Player) cultivationGain() int {
	base := 10 + int(p.Realm)*5 + p.Layer*2
	// 随机浮动 ±20%
	variance := base / 5
	return base + rand.Intn(variance*2+1) - variance
}

// 尝试突破
func (p *Player) Breakthrough() string {
	if p.Cult < p.MaxCult {
		return "修为不足，无法突破"
	}

	if p.Realm == RealmAscension {
		return "你已飞升，无需再突破"
	}

	// 突破成功率
	successRate := p.breakthroughRate()
	roll := rand.Intn(100) + 1
	p.Cult = 0

	if roll <= successRate {
		return p.advanceLayer()
	}
	return p.failBreakthrough()
}

// 突破成功率
func (p *Player) breakthroughRate() int {
	base := 90
	if p.Layer == 9 {
		base = 60 - int(p.Realm)*5 // 大境界突破更危险
	}
	if p.Realm >= RealmTribulation {
		base -= 20 // 渡劫期后更凶险
	}
	if base < 20 {
		base = 20
	}
	return base
}

// 突破成功：提升层数
func (p *Player) advanceLayer() string {
	p.Layer++
	p.updateStats()

	msg := ""
	if p.Layer > 9 {
		// 晋升大境界
		p.Realm++
		p.Layer = 1
		msg = fmt.Sprintf("突破成功！晋升至【%s】！天地变色，灵气汇聚！", p.Realm)

		if p.Realm == RealmAscension {
			msg += "\n\n🎉 恭喜飞升！你已踏入仙界，游戏通关！"
		}
	} else {
		msg = fmt.Sprintf("突破成功！晋升至【%s %d层】！修为大增！", p.Realm, p.Layer)
	}

	p.MaxCult = p.calcMaxCult()
	p.updateStats()
	p.addLog(msg)
	return msg
}

// 突破失败
func (p *Player) failBreakthrough() string {
	damage := p.MaxHP / 4
	p.HP -= damage
	p.addLog(fmt.Sprintf("突破失败！受到 %d 点内伤", damage))

	msg := fmt.Sprintf("突破失败！心魔反噬，损失 %d 点生命", damage)

	if p.HP <= 0 {
		p.HP = 0
		msg += "\n\n💀 突破失败，身死道消……"
	}

	return msg
}

// 更新属性
func (p *Player) updateStats() {
	p.MaxHP = 100 + int(p.Realm)*50 + p.Layer*20
	p.Atk = 10 + int(p.Realm)*8 + p.Layer*3
	p.Def = 5 + int(p.Realm)*5 + p.Layer*2

	// 装备加成
	for _, item := range p.Inventory {
		if item.Amount > 0 {
			p.Atk += item.Atk
			p.Def += item.Def
		}
	}
}

// 计算突破所需修为
func (p *Player) calcMaxCult() int {
	return 100 + int(p.Realm)*200 + p.Layer*50
}

// 使用物品
func (p *Player) UseItem(index int) string {
	if index < 0 || index >= len(p.Inventory) {
		return "无效的物品"
	}

	item := p.Inventory[index]
	if item.Amount <= 0 {
		return "物品已用完"
	}

	switch item.Type {
	case ItemPill:
		p.HP += item.HP
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
		p.Cult += item.Cult
		if p.Cult > p.MaxCult {
			p.Cult = p.MaxCult
		}
		item.Amount--
		return fmt.Sprintf("使用%s，%s", item.Name, item.Desc)

	case ItemWeapon, ItemArmor:
		return fmt.Sprintf("%s 已装备，提供 +%d攻 +%d防", item.Name, item.Atk, item.Def)

	default:
		return "此物品无法使用"
	}
}

// 添加物品
func (p *Player) AddItem(item *Item) {
	for _, inv := range p.Inventory {
		if inv.Name == item.Name && inv.Type == ItemPill {
			inv.Amount += item.Amount
			return
		}
	}
	p.Inventory = append(p.Inventory, item)
}

// 添加日志
func (p *Player) addLog(msg string) {
	p.Log = append(p.Log, msg)
	if len(p.Log) > 50 {
		p.Log = p.Log[len(p.Log)-50:]
	}
}

// 获取最近日志
func (p *Player) RecentLogs(n int) []string {
	if len(p.Log) < n {
		return p.Log
	}
	return p.Log[len(p.Log)-n:]
}

// 玩家状态摘要
func (p *Player) StatusLine() string {
	if p.Realm == RealmAscension {
		return fmt.Sprintf("%s · 飞升之境", p.Name)
	}
	return fmt.Sprintf("%s · %s %d层 · 修为 %d/%d",
		p.Name, p.Realm, p.Layer, p.Cult, p.MaxCult)
}

// 完整状态
func (p *Player) FullStatus() string {
	alive := "存活"
	if p.IsDead() {
		alive = "💀 已陨落"
	}

	return fmt.Sprintf(`姓名: %s  状态: %s
境界: %s %d层
生命: %d/%d  攻击: %d  防御: %d
修为: %d/%d  灵石: %d
击杀: %d`,
		p.Name, alive,
		p.Realm, p.Layer,
		p.HP, p.MaxHP, p.Atk, p.Def,
		p.Cult, p.MaxCult, p.Stones,
		p.TotalKills,
	)
}

func (p *Player) IsDead() bool {
	return p.HP <= 0
}

func (p *Player) IsAscended() bool {
	return p.Realm == RealmAscension
}
