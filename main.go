package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"Go-Server/database"
	"Go-Server/routes"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	_ = godotenv.Load()
	database.ConnectDB()

	app := fiber.New()

	// ✅ Enable CORS (allow all origins)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("📞 Phonebook API is running")
	})

	routes.ContactRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8800"
	}

	log.Fatal(app.Listen(":" + port))
}
