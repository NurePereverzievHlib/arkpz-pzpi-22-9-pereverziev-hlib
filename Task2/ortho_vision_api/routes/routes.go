package routes

import (
	"ortho_vision_api/controllers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes налаштовує маршрути для вашого API.
func SetupRoutes(app *fiber.App) {

	app.Post("/register", controllers.RegisterUser) // Реєстрація користувача

	app.Get("/login", controllers.LoginUser) // Вхід користувача

	app.Put("/users/:id", controllers.UpdateUser) // Оновлення даних користувача

	app.Delete("/admin/users/:id", controllers.DeleteUser) // Видалення користувача за ID

	app.Post("/admin/clinics", controllers.AddClinic) // Додавання клініки

	app.Get("/admin/clinics", controllers.GetAllClinics) // Отримання всіх клінік

	app.Get("/admin/clinics/:name", controllers.GetClinicByName) // Отримання клініки за назвою

	app.Delete("/admin/clinics/:id", controllers.DeleteClinic) // Видалення клініки за ID

	app.Post("/doctor/:doctor_id/appointment_times", controllers.CreateAppointmentTime) // Створення вільного часу доктора

	app.Put("/doctor/:doctor_id/appointment_times/:appointment_time_id", controllers.UpdateAppointmentTime) // Редагування вільного часу доктора

	app.Get("/doctor/:doctor_id/appointment_times", controllers.CreateAppointmentTime) // Отримання всіх вільних місць доктора

	app.Delete("/doctor/:doctor_id/appointment_times/:appointment_time_id", controllers.DeleteAppointmentTime) // Видалення конкретного вільного часу

	app.Get("/appointment-times/search", controllers.SearchAppointmentTimes) // Знайти вільні години до лікаря за часом або лікарем

	app.Post("/appointments", controllers.CreateAppointment) //Запис на прийом

	app.Delete("/appointments/:id", controllers.DeleteAppointment) // Видалення запису на прийом

	app.Get("/appointments/patient/:patientID", controllers.GetAppointmentsByPatientID) // Отримати історію всі прийомів

	app.Post("/diseases", controllers.CreateDisease) // Створення нового запису про хворобу

	app.Delete("/diseases/:id", controllers.DeleteDisease) // Видалення запису про хворобу

	app.Put("/diseases/:id", controllers.UpdateDisease) // Оновлення запису про хворобу

	app.Get("/medical-record/:patientID", controllers.GetMedicalRecord) // Отримання всіх хвороб пацієнта за його ID

	app.Get("/clinic-stats", controllers.GetClinicDiseaseStats)

	// Запити для смарт-окулярів
	app.Post("/smart-glasses", controllers.AddSmartGlassesData)

	app.Get("/smart-glasses/statistics", controllers.GetSmartGlassesStatistics)
}
