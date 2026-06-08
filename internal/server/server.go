package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"

	"xianxia-game/internal/game"
)

//go:embed web/*
var webFiles embed.FS

// 响应结构
type GameState struct {
	Screen     string            `json:"screen"` // main, combat, shop, inventory, cultivating, event, gameover, victory
	Player     *PlayerData       `json:"player"`
	Monster    *MonsterData      `json:"monster,omitempty"`
	CombatLog  []string          `json:"combatLog,omitempty"`
	ShopItems  []ShopItemData    `json:"shopItems,omitempty"`
	Event      *EventData        `json:"event,omitempty"`
	Message    string            `json:"message"`
	HasSave    bool              `json:"hasSave"`
}

type PlayerData struct {
	Name       string      `json:"name"`
	Realm      string      `json:"realm"`
	RealmIdx   int         `json:"realmIdx"`
	Layer      int         `json:"layer"`
	HP         int         `json:"hp"`
	MaxHP      int         `json:"maxHp"`
	Atk        int         `json:"atk"`
	Def        int         `json:"def"`
	Cult       int         `json:"cult"`
	MaxCult    int         `json:"maxCult"`
	Stones     int         `json:"stones"`
	TotalKills int         `json:"totalKills"`
	Inventory  []ItemData  `json:"inventory"`
	Logs       []string    `json:"logs"`
}

type ItemData struct {
	Name   string `json:"name"`
	Type   int    `json:"type"` // 0=丹药 1=武器 2=防具
	Desc   string `json:"desc"`
	HP     int    `json:"hp"`
	Cult   int    `json:"cult"`
	Atk    int    `json:"atk"`
	Def    int    `json:"def"`
	Amount int    `json:"amount"`
}

type MonsterData struct {
	Name   string `json:"name"`
	HP     int    `json:"hp"`
	MaxHP  int    `json:"maxHp"`
	Atk    int    `json:"atk"`
	Def    int    `json:"def"`
	Reward int    `json:"reward"`
}

type ShopItemData struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
	Type  int    `json:"type"`
	HP    int    `json:"hp"`
	Cult  int    `json:"cult"`
	Atk   int    `json:"atk"`
	Def   int    `json:"def"`
}

type EventData struct {
	Title   string        `json:"title"`
	Desc    string        `json:"desc"`
	Choices []ChoiceData  `json:"choices"`
}

type ChoiceData struct {
	Text string `json:"text"`
}

// 全局游戏实例（简化设计，单用户）
var (
	currentPlayer  *game.Player
	currentMonster *game.Monster
	currentEvent   *game.GameEvent
	combatLog      []string
	currentMessage string
)

// 启动服务器
func Start(port int) error {
	webFS, _ := fs.Sub(webFiles, "web")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(webFS)))

	// API
	mux.HandleFunc("/api/state", handleState)
	mux.HandleFunc("/api/new-game", handleNewGame)
	mux.HandleFunc("/api/load", handleLoad)
	mux.HandleFunc("/api/cultivate", handleCultivate)
	mux.HandleFunc("/api/combat-start", handleCombatStart)
	mux.HandleFunc("/api/combat-action", handleCombatAction)
	mux.HandleFunc("/api/shop", handleShop)
	mux.HandleFunc("/api/shop-buy", handleShopBuy)
	mux.HandleFunc("/api/inventory-use", handleInventoryUse)
	mux.HandleFunc("/api/breakthrough", handleBreakthrough)
	mux.HandleFunc("/api/rest", handleRest)
	mux.HandleFunc("/api/event-choose", handleEventChoose)
	mux.HandleFunc("/api/save", handleSave)
	mux.HandleFunc("/api/has-save", handleHasSave)

	addr := fmt.Sprintf(":%d", port)
	go openBrowser(fmt.Sprintf("http://localhost:%d", port))
	return http.ListenAndServe(addr, mux)
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		fmt.Printf("请手动打开浏览器访问: %s\n", url)
	}
}

// 辅助：构建游戏状态
func buildState(screen string) GameState {
	gs := GameState{
		Screen:  screen,
		Message: currentMessage,
		HasSave: game.HasSave(),
	}

	if currentPlayer != nil {
		gs.Player = playerToData(currentPlayer)
	}

	if currentMonster != nil {
		gs.Monster = &MonsterData{
			Name:   currentMonster.Name,
			HP:     currentMonster.HP,
			MaxHP:  currentMonster.MaxHP,
			Atk:    currentMonster.Atk,
			Def:    currentMonster.Def,
			Reward: currentMonster.Reward,
		}
	}

	if len(combatLog) > 0 {
		gs.CombatLog = combatLog
	}

	if currentEvent != nil {
		choices := make([]ChoiceData, len(currentEvent.Choices))
		for i, c := range currentEvent.Choices {
			choices[i] = ChoiceData{Text: c.Text}
		}
		gs.Event = &EventData{
			Title:   currentEvent.Title,
			Desc:    currentEvent.Desc,
			Choices: choices,
		}
	}

	return gs
}

func playerToData(p *game.Player) *PlayerData {
	items := make([]ItemData, len(p.Inventory))
	for i, item := range p.Inventory {
		items[i] = ItemData{
			Name:   item.Name,
			Type:   int(item.Type),
			Desc:   item.Desc,
			HP:     item.HP,
			Cult:   item.Cult,
			Atk:    item.Atk,
			Def:    item.Def,
			Amount: item.Amount,
		}
	}

	return &PlayerData{
		Name:       p.Name,
		Realm:      p.Realm.String(),
		RealmIdx:   int(p.Realm),
		Layer:      p.Layer,
		HP:         p.HP,
		MaxHP:      p.MaxHP,
		Atk:        p.Atk,
		Def:        p.Def,
		Cult:       p.Cult,
		MaxCult:    p.MaxCult,
		Stones:     p.Stones,
		TotalKills: p.TotalKills,
		Inventory:  items,
		Logs:       p.RecentLogs(20),
	}
}

func respond(w http.ResponseWriter, gs GameState) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gs)
}

// ============================================================
// API 处理器
// ============================================================

func handleState(w http.ResponseWriter, r *http.Request) {
	screen := "main"
	if currentPlayer == nil {
		screen = "title"
	} else if currentPlayer.IsDead() {
		screen = "gameover"
	} else if currentPlayer.IsAscended() {
		screen = "victory"
	}
	respond(w, buildState(screen))
}

func handleNewGame(w http.ResponseWriter, r *http.Request) {
	var req struct{ Name string }
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name == "" {
		req.Name = "无名修士"
	}

	currentPlayer = game.NewPlayer(req.Name)
	currentMonster = nil
	currentEvent = nil
	combatLog = nil
	currentMessage = fmt.Sprintf("修士【%s】踏上修仙之路！", req.Name)

	respond(w, buildState("main"))
}

func handleLoad(w http.ResponseWriter, r *http.Request) {
	p, err := game.LoadGame()
	if err != nil {
		respond(w, GameState{Screen: "title", Message: "没有找到存档"})
		return
	}
	currentPlayer = p
	currentMonster = nil
	currentEvent = nil
	combatLog = nil
	currentMessage = "存档加载成功！"

	screen := "main"
	if p.IsDead() {
		screen = "gameover"
	}
	if p.IsAscended() {
		screen = "victory"
	}
	respond(w, buildState(screen))
}

func handleCultivate(w http.ResponseWriter, r *http.Request) {
	if currentPlayer == nil {
		respond(w, GameState{Screen: "title"})
		return
	}

	currentMessage = currentPlayer.Cultivate()

	// 随机事件
	currentEvent = game.GetRandomEvent(currentPlayer)
	if currentEvent != nil {
		respond(w, buildState("event"))
		return
	}

	respond(w, buildState("cultivating"))
}

func handleCombatStart(w http.ResponseWriter, r *http.Request) {
	if currentPlayer == nil {
		respond(w, GameState{Screen: "title"})
		return
	}

	currentMonster = game.GenerateMonster(currentPlayer)
	combatLog = []string{
		fmt.Sprintf("遭遇了 %s！准备战斗！", currentMonster.Name),
	}
	currentMessage = ""
	respond(w, buildState("combat"))
}

func handleCombatAction(w http.ResponseWriter, r *http.Request) {
	if currentPlayer == nil || currentMonster == nil {
		respond(w, GameState{Screen: "main"})
		return
	}

	var req struct {
		Action    string `json:"action"`
		ItemIndex int    `json:"itemIndex"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	switch req.Action {
	case "attack":
		msg := game.PlayerAttack(currentPlayer, currentMonster)
		combatLog = append(combatLog, msg)

		if currentMonster.HP <= 0 {
			reward := game.CombatReward(currentMonster, currentPlayer)
			combatLog = append(combatLog, reward)
			currentMonster = nil
			currentMessage = reward
			currentPlayer.Save()
			respond(w, buildState("main"))
			return
		}

		// 妖兽反击
		msg = game.MonsterAttack(currentPlayer, currentMonster)
		combatLog = append(combatLog, msg)

		if currentPlayer.IsDead() {
			currentMonster = nil
			respond(w, buildState("gameover"))
			return
		}

	case "item":
		if req.ItemIndex >= 0 && req.ItemIndex < len(currentPlayer.Inventory) {
			item := currentPlayer.Inventory[req.ItemIndex]
			if item.Type != game.ItemPill {
				combatLog = append(combatLog, "战斗中只能使用丹药！")
			} else {
				msg := currentPlayer.UseItem(req.ItemIndex)
				combatLog = append(combatLog, msg)

				// 妖兽反击
				msg = game.MonsterAttack(currentPlayer, currentMonster)
				combatLog = append(combatLog, msg)

				if currentPlayer.IsDead() {
					currentMonster = nil
					respond(w, buildState("gameover"))
					return
				}
			}
		}

	case "flee":
		if game.TryFlee(currentPlayer, currentMonster) {
			combatLog = append(combatLog, "你成功逃脱了！")
			currentMonster = nil
			currentMessage = "你逃离了战斗。"
			respond(w, buildState("main"))
			return
		}
		combatLog = append(combatLog, "逃跑失败！")
		msg := game.MonsterAttack(currentPlayer, currentMonster)
		combatLog = append(combatLog, msg)
		if currentPlayer.IsDead() {
			currentMonster = nil
			respond(w, buildState("gameover"))
			return
		}
	}

	respond(w, buildState("combat"))
}

func handleShop(w http.ResponseWriter, r *http.Request) {
	if currentPlayer == nil {
		respond(w, GameState{Screen: "title"})
		return
	}

	items := game.GetShopItems(currentPlayer)
	shopItems := make([]ShopItemData, len(items))
	for i, si := range items {
		shopItems[i] = ShopItemData{
			Name:  si.Item.Name,
			Desc:  si.Item.Desc,
			Price: si.Price,
			Type:  int(si.Item.Type),
			HP:    si.Item.HP,
			Cult:  si.Item.Cult,
			Atk:   si.Item.Atk,
			Def:   si.Item.Def,
		}
	}

	gs := buildState("shop")
	gs.ShopItems = shopItems
	respond(w, gs)
}

func handleShopBuy(w http.ResponseWriter, r *http.Request) {
	var req struct{ Index int }
	json.NewDecoder(r.Body).Decode(&req)

	items := game.GetShopItems(currentPlayer)
	msg, _ := game.BuyItem(currentPlayer, items, req.Index)
	currentMessage = msg

	items = game.GetShopItems(currentPlayer)
	shopItems := make([]ShopItemData, len(items))
	for i, si := range items {
		shopItems[i] = ShopItemData{
			Name:  si.Item.Name,
			Desc:  si.Item.Desc,
			Price: si.Price,
			Type:  int(si.Item.Type),
			HP:    si.Item.HP,
			Cult:  si.Item.Cult,
			Atk:   si.Item.Atk,
			Def:   si.Item.Def,
		}
	}

	gs := buildState("shop")
	gs.ShopItems = shopItems
	respond(w, gs)
}

func handleInventoryUse(w http.ResponseWriter, r *http.Request) {
	var req struct{ Index int }
	json.NewDecoder(r.Body).Decode(&req)

	currentMessage = currentPlayer.UseItem(req.Index)
	respond(w, buildState("inventory"))
}

func handleBreakthrough(w http.ResponseWriter, r *http.Request) {
	if currentPlayer == nil {
		respond(w, GameState{Screen: "title"})
		return
	}

	currentMessage = currentPlayer.Breakthrough()

	if currentPlayer.IsDead() {
		respond(w, buildState("gameover"))
		return
	}
	if currentPlayer.IsAscended() {
		respond(w, buildState("victory"))
		return
	}

	currentPlayer.Save()
	respond(w, buildState("cultivating"))
}

func handleRest(w http.ResponseWriter, r *http.Request) {
	heal := currentPlayer.MaxHP / 3
	currentPlayer.HP += heal
	if currentPlayer.HP > currentPlayer.MaxHP {
		currentPlayer.HP = currentPlayer.MaxHP
	}
	currentMessage = fmt.Sprintf("打坐调息，恢复了 %d 点生命。", heal)
	respond(w, buildState("cultivating"))
}

func handleEventChoose(w http.ResponseWriter, r *http.Request) {
	var req struct{ ChoiceIndex int }
	json.NewDecoder(r.Body).Decode(&req)

	if currentEvent != nil && req.ChoiceIndex >= 0 && req.ChoiceIndex < len(currentEvent.Choices) {
		choice := currentEvent.Choices[req.ChoiceIndex]
		currentMessage = choice.Action(currentPlayer)
		currentEvent = nil
		currentPlayer.Save()
	}

	respond(w, buildState("main"))
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	if currentPlayer != nil {
		currentPlayer.Save()
		currentMessage = "保存成功！"
	}
	respond(w, buildState("main"))
}

func handleHasSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"hasSave": game.HasSave()})
}
