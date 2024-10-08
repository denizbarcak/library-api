package main

import (
	"library-api/configs"
	"library-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
    app := fiber.New()

    // CORS middleware
    app.Use(cors.New())

    // Veritabanı bağlantısı
    configs.ConnectDB()

    // Rotalar
    routes.UserRoutes(app)
    routes.BookRoutes(app)

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendFile("./frontend/login.html")
    })

    // Frontend dosyasını sunma
    app.Static("/", "./frontend")

    // Sunucu başlatma
    app.Listen(":3000")
}
