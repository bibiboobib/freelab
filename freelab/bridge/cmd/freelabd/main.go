package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"freelab-bridge/internal/api"
	"freelab-bridge/internal/qemu" // Подключаем наш новый пакет QEMU
	"freelab-bridge/internal/state"
)

func main() {
	log.Println("Starting freelab Bridge Daemon on 0.0.0.0:8080...")

	// 1. Инициализируем центральное состояние
	stateManager := state.NewManager()

	// 2. Инициализируем QEMU менеджер (указываем папки в WSL)
	qemuManager := qemu.NewQemuManager("/tmp/freelab/base", "/tmp/freelab/overlays")

	// 3. Передаем ОБА менеджера в API роутер
	router := api.NewRouter(stateManager, qemuManager)

	// Запускаем сервер...
	err := http.ListenAndServe("0.0.0.0:8080", router)
	
	if err != nil {
		fmt.Printf("\n[!] КРИТИЧЕСКАЯ ОШИБКА СЕРВЕРА: %v\n", err)
		fmt.Println("Нажмите Enter, чтобы выйти...")
		fmt.Scanln() // Окно остановится здесь и будет ждать вашего нажатия
		os.Exit(1)
	}
}
