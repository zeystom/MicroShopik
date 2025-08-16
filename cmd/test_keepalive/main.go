package main

import (
	"MicroShopik/internal/services"
	"log"
	"time"
)

func main() {
	log.Println("Testing keep-alive service...")

	// Создаем сервис с интервалом 30 секунд для тестирования
	keepAliveService := services.NewKeepAliveService(
		"http://localhost:8080",
		30*time.Second,
	)

	// Запускаем сервис
	keepAliveService.Start()

	// Ждем 2 минуты для демонстрации
	log.Println("Keep-alive service is running. Press Ctrl+C to stop...")
	time.Sleep(2 * time.Minute)

	// Останавливаем сервис
	keepAliveService.Stop()
	log.Println("Test completed.")
}
