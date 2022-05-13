package aggregates

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
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

func GetExpensesForMonth(month AggregateDate, catId int64) ([]ExpenseDto, error) {
	log.Printf("Get expenses for month %v for id %v\n",month, catId)
	rows, err := database.Conn.Query(context.Background(),"SELECT id, description, date, price, categoryid FROM expense where categoryid = $1 and date >= $2 and date <= $3", catId, month.FirstDay, month.LastDay)
	defer rows.Close()
	if err != nil {
		return make([]ExpenseDto, 0), err
	}

	return getExpenseFromRows(rows), nil
}

func getExpenseFromRows(rows pgx.Rows) []ExpenseDto{
	var expenses []ExpenseDto
	for rows.Next() {
		var id int64
		var description string
		var date time.Time
		var price int
		var categoryid int64
		rows.Scan(&id,&description,&date,&price,&categoryid)
		expenses = append(expenses, ExpenseDto{id,date,description,price,categoryid})
	}

	return expenses
}