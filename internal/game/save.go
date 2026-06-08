package game

import (
	"encoding/json"
	"os"
)

const saveFileName = "xianxia_save.json"

type SaveData struct {
	Name       string
	Realm      Realm
	Layer      int
	HP         int
	MaxHP      int
	Atk        int
	Def        int
	Cult       int
	MaxCult    int
	Stones     int
	TotalKills int
	Inventory  []*Item
}

// 保存游戏
func (p *Player) Save() error {
	data := SaveData{
		Name:       p.Name,
		Realm:      p.Realm,
		Layer:      p.Layer,
		HP:         p.HP,
		MaxHP:      p.MaxHP,
		Atk:        p.Atk,
		Def:        p.Def,
		Cult:       p.Cult,
		MaxCult:    p.MaxCult,
		Stones:     p.Stones,
		TotalKills: p.TotalKills,
		Inventory:  p.Inventory,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(saveFileName, bytes, 0644)
}

// 加载游戏
func LoadGame() (*Player, error) {
	bytes, err := os.ReadFile(saveFileName)
	if err != nil {
		return nil, err
	}

	var data SaveData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	p := NewPlayer(data.Name)
	p.Realm = data.Realm
	p.Layer = data.Layer
	p.HP = data.HP
	p.MaxHP = data.MaxHP
	p.Atk = data.Atk
	p.Def = data.Def
	p.Cult = data.Cult
	p.MaxCult = data.MaxCult
	p.Stones = data.Stones
	p.TotalKills = data.TotalKills
	p.Inventory = data.Inventory

	return p, nil
}

// 检查是否有存档
func HasSave() bool {
	_, err := os.Stat(saveFileName)
	return err == nil
}

// 删除存档
func DeleteSave() {
	os.Remove(saveFileName)
}
