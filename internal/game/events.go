package game

import (
	"fmt"
	"math/rand"
)

// 获取一个随机事件（修炼时有概率触发）
func GetRandomEvent(p *Player) *GameEvent {
	roll := rand.Intn(100)
	category := rand.Intn(5)

	switch {
	case roll < 30: // 30% 概率触发事件
		switch category {
		case 0:
			return treasureEvent(p)
		case 1:
			return immortalEvent(p)
		case 2:
			return ambushEvent(p)
		case 3:
			return herbEvent(p)
		default:
			return meditationEvent(p)
		}
	}

	return nil
}

// 发现宝物
func treasureEvent(p *Player) *GameEvent {
	stones := 20 + rand.Intn(int(p.Realm)*30+50)

	return &GameEvent{
		Title: "意外发现",
		Desc:  fmt.Sprintf("修炼时，你感应到地底有灵力波动。仔细探查，发现了一处前人遗留的藏宝洞！\n里面藏有 %d 灵石。", stones),
		Choices: []EventChoice{
			{
				Text: "取走灵石",
				Action: func(p *Player) string {
					p.Stones += stones
					p.addLog(fmt.Sprintf("发现藏宝洞，获得 %d 灵石", stones))
					return fmt.Sprintf("你小心翼翼地取走灵石，获得 %d 灵石！", stones)
				},
			},
			{
				Text: "留给后人，不取分毫",
				Action: func(p *Player) string {
					cultGain := p.cultivationGain() * 3
					p.Cult += cultGain
					if p.Cult > p.MaxCult {
						p.Cult = p.MaxCult
					}
					p.addLog("你选择不取宝物，心境提升")
					return fmt.Sprintf("你的无私之举感动了洞府主人残存的神念，获得 %d 点修为领悟！", cultGain)
				},
			},
		},
	}
}

// 遇到仙人
func immortalEvent(p *Player) *GameEvent {
	return &GameEvent{
		Title: "仙人指路",
		Desc:  "一位白发苍苍的老者突然出现在你面前，他周身散发着淡淡仙光。\n'小友，你我有缘，老夫送你一场造化。'",
		Choices: []EventChoice{
			{
				Text: "恭敬请教",
				Action: func(p *Player) string {
					gain := p.MaxCult / 2
					p.Cult += gain
					if p.Cult > p.MaxCult {
						p.Cult = p.MaxCult
					}
					p.addLog("仙人指点，修为大增")
					return fmt.Sprintf("仙人传授了一段口诀，你顿悟其中玄妙，获得 %d 点修为！", gain)
				},
			},
			{
				Text: "挑战仙人，试试自己实力",
				Action: func(p *Player) string {
					if rand.Intn(100) < 30 {
						p.HP -= p.MaxHP / 2
						p.addLog("挑战仙人失败，身受重伤")
						return "仙人轻轻一拂袖，你倒飞出去，身受重伤！损失了一半生命。"
					}
					gain := p.MaxCult
					p.Cult += gain
					if p.Cult > p.MaxCult {
						p.Cult = p.MaxCult
					}
					p.addLog("挑战仙人获得认可")
					return fmt.Sprintf("仙人哈哈大笑：'好胆量！' 非但没有生气，反而赐你一场造化，获得 %d 点修为！", gain)
				},
			},
		},
	}
}

// 遭遇袭击
func ambushEvent(p *Player) *GameEvent {
	dmg := p.MaxHP / 5

	return &GameEvent{
		Title: "遭遇袭击",
		Desc:  "修炼正酣时，突然窜出几个蒙面修士！\n'交出灵石，饶你不死！'",
		Choices: []EventChoice{
			{
				Text: "奋起反击",
				Action: func(p *Player) string {
					if rand.Intn(100) < 50+int(p.Realm)*5 {
						stones := 30 + rand.Intn(50)
						p.Stones += stones
						p.TotalKills++
						p.addLog("击败了来袭的修士")
						return fmt.Sprintf("你三两下击退了来袭者，反而从他们身上搜出了 %d 灵石！", stones)
					}
					p.HP -= dmg
					p.addLog("被修士袭击受伤")
					return fmt.Sprintf("不敌对方人多势众，你受了些伤，损失 %d 生命。", dmg)
				},
			},
			{
				Text: "破财消灾",
				Action: func(p *Player) string {
					lost := p.Stones / 5
					if lost < 10 {
						lost = 10
					}
					if p.Stones < lost {
						lost = p.Stones
					}
					p.Stones -= lost
					p.addLog(fmt.Sprintf("花钱消灾，损失 %d 灵石", lost))
					return fmt.Sprintf("你扔出 %d 灵石，趁对方争抢时脱身而去。", lost)
				},
			},
		},
	}
}

// 发现灵草
func herbEvent(p *Player) *GameEvent {
	return &GameEvent{
		Title: "发现灵草",
		Desc:  "你在山间漫步时，发现了一株散发着微光的灵草。\n这株灵草灵气充沛，似乎可以直接服用，也可以拿去售卖。",
		Choices: []EventChoice{
			{
				Text: "直接服用",
				Action: func(p *Player) string {
					heal := p.MaxHP / 3
					p.HP += heal
					if p.HP > p.MaxHP {
						p.HP = p.MaxHP
					}
					p.addLog("服用了灵草，生命恢复")
					return fmt.Sprintf("你服下灵草，一股暖流涌遍全身，恢复了 %d 点生命！", heal)
				},
			},
			{
				Text: "炼化成丹药",
				Action: func(p *Player) string {
					cult := p.MaxCult / 3
					p.Cult += cult
					if p.Cult > p.MaxCult {
						p.Cult = p.MaxCult
					}
					p.addLog("将灵草炼化，获得修为")
					return fmt.Sprintf("你将灵草炼化，获得了 %d 点修为！", cult)
				},
			},
		},
	}
}

// 冥想感悟
func meditationEvent(p *Player) *GameEvent {
	return &GameEvent{
		Title: "天地感悟",
		Desc:  "你静坐修炼时，忽然感受到了天地间一丝大道韵律。\n你面临两个选择：追求力量的极致，还是追求境界的圆满？",
		Choices: []EventChoice{
			{
				Text: "追求力量 (+攻击)",
				Action: func(p *Player) string {
					p.Atk += 5 + int(p.Realm)*2
					p.addLog("感悟力量之道，攻击提升")
					return fmt.Sprintf("你明悟了力量的本质，攻击力永久提升 %d 点！", 5+int(p.Realm)*2)
				},
			},
			{
				Text: "追求境界 (+修为)",
				Action: func(p *Player) string {
					gain := p.MaxCult
					p.Cult += gain
					if p.Cult > p.MaxCult {
						p.Cult = p.MaxCult
					}
					p.addLog("感悟天地大道，修为精进")
					return fmt.Sprintf("你沉浸在天地韵律中，修为暴涨 %d 点！", gain)
				},
			},
		},
	}
}
