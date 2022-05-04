package api

import (
	"github.com/gofiber/fiber/v2"
	"main.go/aggregates"
)

func InitRoutes(app *fiber.App){
	aggGroup := app.Group("/aggregate/month")
	catGroup := app.Group("/categories")
	expGroup := app.Group("/expenses")

	initAggRoutes(aggGroup)
	initCatRoutes(catGroup)
	initExpRoutes(expGroup)
}

func initAggRoutes(grp fiber.Router){
	grp.Get("/:page/:items", func(c *fiber.Ctx) error {
		return aggregates.GetAggregatesInMonths(c)
	})
}

func initCatRoutes(grp fiber.Router){
	grp.Post("/", func(c *fiber.Ctx) error {
		return aggregates.InsertCategory(c)
	})
}

func initExpRoutes(grp fiber.Router){
	grp.Post("/", func(c *fiber.Ctx) error {
		return aggregates.InsertExpense(c)
	})
}