package main

import (
	"log"
	"ortho_vision_api/config"
	"ortho_vision_api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Налаштування підключення до БД
	_ = config.InitDB()

	// Створення нового серверу на Fiber
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		// Додаємо з'єднання з базою даних у контекст
		c.Locals("db", config.DB)
		return c.Next()
	})

	// Налаштування маршрутів
	routes.SetupRoutes(app)

	// Запуск сервера
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
