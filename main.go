package main

import (
	"fmt"
	"xianxia-game/internal/server"
)

func main() {
	fmt.Println("╔══════════════════════════════════╗")
	fmt.Println("║     修 仙 录 — 凡人修仙传       ║")
	fmt.Println("║     浏览器中将打开游戏界面      ║")
	fmt.Println("║     按 Ctrl+C 退出              ║")
	fmt.Println("╚══════════════════════════════════╝")

	if err := server.Start(8080); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
