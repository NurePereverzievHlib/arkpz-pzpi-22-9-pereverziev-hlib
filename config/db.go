package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB — це глобальна змінна для зберігання підключення до БД.
var DB *gorm.DB

// InitDB ініціалізує з'єднання з базою даних.
func InitDB() *gorm.DB {
	// Прямо в коді вставляємо конкретні дані для підключення:
	dbUser := "postgres"      // Твій користувач PostgreSQL
	dbPass := "31415"         // Твій пароль для PostgreSQL (заміни на справжній)
	dbName := "orthovisiondb" // Назва твоєї бази даних
	dbHost := "localhost"     // Адреса сервера
	dbPort := "5432"          // Порт сервера PostgreSQL

	// Створюємо рядок з'єднання для PostgreSQL.
	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", dbHost, dbUser, dbName, dbPass, dbPort)

	// Підключаємося до бази даних.
	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Перевірка з'єднання
	sqlDB, err := DB.DB() // Отримуємо доступ до звичайного SQL драйвера для виконання запиту
	if err != nil {
		log.Fatal("Failed to get DB connection:", err)
	}

	// Перевірка з'єднання з базою даних
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	} else {
		log.Println("Successfully connected to the database.")
	}

	// Автоматичне створення таблиць при запуску програми (якщо їх немає).
	// Якщо потрібно зробити тільки міграцію, можна замінити db.AutoMigrate() на інші міграційні інструменти.
	DB.AutoMigrate()

	// Повертаємо підключення до БД для використання в інших частинах програми.
	return DB
}
