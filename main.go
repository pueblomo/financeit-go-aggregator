package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"main.go/api"
	"main.go/conf"
	"main.go/database"
)


func main() {
	conf.InitEnv()

	database.Connect()
	defer database.Conn.Close()
	app := fiber.New()

	api.InitRoutes(app)

	log.Fatalln(app.Listen(":7080"))
}