package game

import (
	"fmt"
	"math/rand"
)

// 妖兽名字库
var monsterNames = [][2]string{
	{"噬金鼠", "一只浑身泛着金属光泽的巨大鼠妖"},
	{"赤练蛇", "通体赤红的毒蛇，喷吐毒雾"},
	{"幽冥狼", "来自幽冥之地的狼妖，双眼闪烁着绿光"},
	{"石甲熊", "皮肤如岩石般坚硬的熊妖"},
	{"血翼蝠", "双翼呈血红色的巨型蝙蝠"},
	{"碧水蛟", "潜伏在水中的蛟龙，鳞片闪烁着碧光"},
	{"烈焰虎", "周身环绕烈焰的猛虎"},
	{"冰霜蟒", "吐息可冻结万物的巨蟒"},
	{"金翅鹏", "双翅展开遮天蔽日的金色大鹏"},
	{"紫电貂", "速度极快、周身缠绕紫色闪电的貂妖"},
	{"玄龟", "背甲如山岳般厚重的千年玄龟"},
	{"九尾狐", "生有九尾的妖狐，魅惑众生"},
	{"铁甲蝎", "尾钩闪烁着寒光的巨蝎"},
	{"鬼面蛛", "腹部花纹似鬼面的巨型蜘蛛"},
	{"黑水玄蛇", "来自深渊的黑色巨蛇"},
	{"银月狼王", "月圆之夜力量倍增的狼族王者"},
	{"三眼蟾蜍", "额头生有第三只眼的妖蟾"},
	{"青鸾", "羽毛如青玉般的神鸟后裔"},
	{"地龙", "潜伏地下的龙族旁支"},
	{"雷鹰", "羽翼间有雷电闪烁的巨鹰"},
}

// 根据玩家境界生成妖兽
func GenerateMonster(player *Player) *Monster {
	idx := rand.Intn(len(monsterNames))
	name := monsterNames[idx][0]

	// 妖兽属性随玩家境界增长
	realmBonus := int(player.Realm)*20 + player.Layer*5
	level := rand.Intn(3) // 0=普通 1=精英 2=BOSS

	m := &Monster{
		Name:  name,
		MaxHP: 50 + realmBonus + level*50 + rand.Intn(50),
		Atk:   8 + realmBonus/2 + level*5 + rand.Intn(10),
		Def:   3 + realmBonus/3 + level*3 + rand.Intn(5),
	}
	m.HP = m.MaxHP

	// 灵石奖励
	m.Reward = 10 + int(player.Realm)*15 + level*30 + rand.Intn(30)

	// 掉落物品（概率）
	dropChance := 20 + int(player.Realm)*5
	if level == 2 {
		dropChance = 60
	}
	if rand.Intn(100) < dropChance {
		m.DropItem = randomDrop(player.Realm)
	}

	if level == 1 {
		m.Name = "【精英】" + m.Name
	}
	if level == 2 {
		m.Name = "【BOSS】" + m.Name
		m.MaxHP *= 2
		m.HP = m.MaxHP
		m.Atk = m.Atk * 3 / 2
		m.Reward *= 3
	}

	return m
}

// 随机掉落
func randomDrop(realm Realm) *Item {
	drops := []*Item{
		{Name: "回春丹", Type: ItemPill, Desc: "恢复80点生命", Value: 15, HP: 80, Amount: 1},
		{Name: "大还丹", Type: ItemPill, Desc: "恢复200点生命", Value: 30, HP: 200, Amount: 1},
		{Name: "聚气散", Type: ItemPill, Desc: "增加150点修为", Value: 25, Cult: 150, Amount: 1},
		{Name: "凝神丹", Type: ItemPill, Desc: "增加300点修为", Value: 50, Cult: 300, Amount: 1},
		{Name: "青锋剑", Type: ItemWeapon, Desc: "凡铁锻造的利剑", Value: 100, Atk: 15, Amount: 1},
		{Name: "玄铁重剑", Type: ItemWeapon, Desc: "以玄铁铸成", Value: 300, Atk: 30, Amount: 1},
		{Name: "飞虹剑", Type: ItemWeapon, Desc: "剑身如虹，削铁如泥", Value: 800, Atk: 60, Amount: 1},
		{Name: "铁布衫", Type: ItemArmor, Desc: "普通护甲", Value: 80, Def: 10, Amount: 1},
		{Name: "金丝甲", Type: ItemArmor, Desc: "以金丝编织的软甲", Value: 250, Def: 25, Amount: 1},
		{Name: "天蚕衣", Type: ItemArmor, Desc: "以天蚕丝织成的宝衣", Value: 600, Def: 50, Amount: 1},
	}

	idx := rand.Intn(len(drops))
	// 根据境界提高高级物品概率
	if int(realm) > 3 && rand.Intn(2) == 0 {
		idx = rand.Intn(len(drops))
	}

	return drops[idx]
}

// 玩家攻击
func PlayerAttack(p *Player, m *Monster) string {
	dmg := p.Atk - m.Def
	if dmg < 1 {
		dmg = 1
	}
	// 暴击 15%
	if rand.Intn(100) < 15 {
		dmg *= 2
		m.HP -= dmg
		return fmt.Sprintf("暴击！对%s造成 %d 点伤害！%s剩余HP: %d/%d",
			m.Name, dmg, m.Name, m.HP, m.MaxHP)
	}

	m.HP -= dmg
	return fmt.Sprintf("你对%s造成 %d 点伤害。%s剩余HP: %d/%d",
		m.Name, dmg, m.Name, m.HP, m.MaxHP)
}

// 妖兽攻击
func MonsterAttack(p *Player, m *Monster) string {
	dmg := m.Atk - p.Def
	if dmg < 1 {
		dmg = 1
	}
	// 暴击 10%
	if rand.Intn(100) < 10 {
		dmg *= 2
		p.HP -= dmg
		return fmt.Sprintf("%s暴击！对你造成 %d 点伤害！你的HP: %d/%d",
			m.Name, dmg, p.HP, p.MaxHP)
	}

	p.HP -= dmg
	return fmt.Sprintf("%s对你造成 %d 点伤害。你的HP: %d/%d",
		m.Name, dmg, p.HP, p.MaxHP)
}

// 尝试逃跑
func TryFlee(p *Player, m *Monster) bool {
	// 逃跑成功率与境界相关
	chance := 50 + int(p.Realm)*5
	if chance > 90 {
		chance = 90
	}
	return rand.Intn(100) < chance
}

// 战胜妖兽后的奖励描述
func CombatReward(m *Monster, p *Player) string {
	p.Stones += m.Reward
	p.TotalKills++

	msg := fmt.Sprintf("战胜%s！获得 %d 灵石", m.Name, m.Reward)

	if m.DropItem != nil {
		p.AddItem(m.DropItem)
		msg += fmt.Sprintf("，掉落【%s】", m.DropItem.Name)
	}

	p.addLog(msg)
	return msg
}
