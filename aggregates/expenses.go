package aggregates

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"main.go/database"
)


type ExpenseDto struct {
	Id int64 `json:"id"`
	Date time.Time `json:"date" validate:"required"`
	Description string `json:"description" validate:"required,max=60"`
	Price int `json:"price" validate:"required,number,min=0"`
	CategoryId int64 `json:"categoryId" validate:"required,min=0"`
}

func InsertExpense (c *fiber.Ctx) error {
	c.Accepts("application/json")
	var newExp ExpenseDto
	log.Println("Post expense")
	if err:= c.BodyParser(&newExp); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err:= validate.Struct(newExp); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	log.Printf("Post expense: %v\n", newExp)
	_, err := database.Conn.Exec(context.Background(),"INSERT INTO expense (description, date, price, categoryid) VALUES ($1, $2, $3, $4)", newExp.Description, newExp.Date, newExp.Price, newExp.CategoryId)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusCreated)
}