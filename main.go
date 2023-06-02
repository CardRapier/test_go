package main

import (
	"github.com/cardrapier/hello-fiber/database"
	"github.com/cardrapier/hello-fiber/motel"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.RunDB()
	motel.SetupRoutes(app)

	app.Listen(":3000")
}
