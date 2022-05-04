package aggregates

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"main.go/conf"
	"main.go/database"
	"main.go/util"
)

type CategoryDto struct {
	Id int64 `json:"id"`
	Name string `json:"name" validate:"required,max=60"`
	Icon string `json:"icon" validate:"required,max=20"`
	Amount int64 `json:"amount"`
}

func getCategoriesFromRows(rows pgx.Rows, month AggregateDate) []CategoryDto {
	var categories []CategoryDto
	for rows.Next() {
		var id int64
		var name string
		var icon string
		rows.Scan(&id, &name, &icon)
		categories = append(categories, CategoryDto{id,name,icon, int64(getAmount(id, month))})
	}
	return categories
}

func GetCategoriesForMonth(month AggregateDate) ([]CategoryDto,error){
	log.Printf("Get Categories for %v\n", month)
	rows, err := database.Conn.Query(context.Background(),"SELECT categoryid, name, icon FROM category where valid_from <= $1 and valid_to >= $1 order by name", month.LastDay)
	defer rows.Close()
	if err != nil {
		return make([]CategoryDto, 0), err
	}

	return getCategoriesFromRows(rows, month), nil
}

func getAmount(id int64, month AggregateDate) int {
	rows, err := database.Conn.Query(context.Background(),"SELECT SUM (price) FROM expense WHERE categoryid = $1 and date >= $2 and date <= $3", strconv.FormatInt(id,10), month.FirstDay, month.LastDay)
	conf.CheckError(err)
	var amount int
	for rows.Next() {
		rows.Scan(&amount)
	}

	return amount
}

func InsertCategory(ctx *fiber.Ctx) error{
	ctx.Accepts("application/json")
	var newCat CategoryDto
	log.Println("Post category")
	if err:= ctx.BodyParser(&newCat); err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if err:= validate.Struct(newCat); err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	log.Printf("Post Categorie : %v\n", newCat)
	t := strings.Split(util.GetDateWithFirstDay(time.Now()).String(), " ") 
	_, err := database.Conn.Exec(context.Background(),"INSERT INTO category (categoryid, name, icon, valid_from, valid_to) VALUES ($1, $2, $3, $4, $5)", newCat.Id, newCat.Name, newCat.Icon, t[0] + " " + t[1],"2999-04-30 11:42:28.000000")
	if err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
