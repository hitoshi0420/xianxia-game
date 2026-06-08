package game

import "fmt"

// 境界
type Realm int

const (
	RealmMortal          Realm = iota // 凡人
	RealmQiRefining                   // 炼气
	RealmFoundation                   // 筑基
	RealmGoldenCore                   // 金丹
	RealmNascentSoul                  // 元婴
	RealmSpiritTransform              // 化神
	RealmTribulation                  // 渡劫
	RealmMahayana                     // 大乘
	RealmAscension                    // 飞升（通关）
)

var realmNames = map[Realm]string{
	RealmMortal:           "凡人",
	RealmQiRefining:       "炼气期",
	RealmFoundation:       "筑基期",
	RealmGoldenCore:       "金丹期",
	RealmNascentSoul:      "元婴期",
	RealmSpiritTransform:  "化神期",
	RealmTribulation:      "渡劫期",
	RealmMahayana:         "大乘期",
	RealmAscension:        "飞升之境",
}

func (r Realm) String() string {
	if name, ok := realmNames[r]; ok {
		return name
	}
	return "未知"
}

// 物品类型
type ItemType int

const (
	ItemPill     ItemType = iota // 丹药
	ItemWeapon                   // 武器
	ItemArmor                    // 防具
	ItemMaterial                 // 材料
)

// 物品
type Item struct {
	Name   string
	Type   ItemType
	Desc   string
	Value  int // 售价（灵石）
	HP     int // 恢复生命
	Cult   int // 增加修为
	Atk    int // 攻击加成
	Def    int // 防御加成
	Amount int // 数量（可堆叠）
}

// 妖兽
type Monster struct {
	Name    string
	HP      int
	MaxHP   int
	Atk     int
	Def     int
	Reward  int // 灵石奖励
	DropItem *Item
}

// 战斗结果
type CombatResult int

const (
	CombatWin  CombatResult = iota
	CombatLose
	CombatFlee
	CombatOngoing
)

// 事件
type GameEvent struct {
	Title   string
	Desc    string
	Choices []EventChoice
}

type EventChoice struct {
	Text   string
	Action func(*Player) string // 返回结果描述
}

// 辅助函数
func FormatLargeNum(n int) string {
	if n >= 100000000 {
		return fmt.Sprintf("%.1f亿", float64(n)/100000000)
	}
	if n >= 10000 {
		return fmt.Sprintf("%.1f万", float64(n)/10000)
	}
	return fmt.Sprintf("%d", n)
}
