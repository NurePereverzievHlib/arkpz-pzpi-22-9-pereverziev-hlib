package controllers

import (
	"log"
	"ortho_vision_api/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser - функція для реєстрації нового користувача
func RegisterUser(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Створюємо нову структуру для користувача
	var user models.User

	// Парсимо тіло запиту в структуру user
	if err := c.BodyParser(&user); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Перевірка на унікальність email
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		log.Println("User with this email already exists")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "User with this email already exists",
		})
	}

	// Хешуємо пароль перед збереженням
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}

	// Присвоюємо хешоване значення паролю
	user.PasswordHash = string(hashedPassword)

	// Зберігаємо нового користувача в базі
	if err := db.Create(&user).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving user to database",
		})
	}

	// Відповідь про успішну реєстрацію
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

// LoginUser - функція для авторизації користувача
func LoginUser(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Створюємо структуру для даних, які прийдуть у запиті
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Парсимо тіло запиту
	if err := c.BodyParser(&loginData); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Знаходимо користувача за email
	var user models.User
	if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		log.Println("User not found:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Перевіряємо пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password)); err != nil {
		log.Println("Password mismatch")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Створюємо анонімну структуру для відповіді, щоб не включати PasswordHash
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			// Додайте інші необхідні поля тут, які потрібно відправити у відповіді
		},
	})
}

// UpdateUser - функція для оновлення даних користувача
func UpdateUser(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID користувача з параметрів URL
	userID := c.Params("id")

	// Перевіряємо, чи існує користувач із цим ID
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		log.Println("User not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Структура для отримання оновлених даних
	var updatedData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Парсимо тіло запиту
	if err := c.BodyParser(&updatedData); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Оновлюємо поля
	if updatedData.Name != "" {
		user.Name = updatedData.Name
	}
	if updatedData.Email != "" {
		user.Email = updatedData.Email
	}
	if updatedData.Password != "" {
		// Хешуємо новий пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Password hashing error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error hashing password",
			})
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Зберігаємо оновлення в базу
	if err := db.Save(&user).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating user",
		})
	}

	// Відповідь про успішне оновлення
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    user,
	})
}

// DeleteUser - функція для видалення користувача за ID
func DeleteUser(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID користувача з параметрів
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User ID is required",
		})
	}

	// Перевірка, чи існує користувач
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("User not found:", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		log.Println("Error finding user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding user",
		})
	}

	// Видалення користувача
	if err := db.Delete(&user).Error; err != nil {
		log.Println("Error deleting user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting user",
		})
	}

	// Відповідь про успішне видалення
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
