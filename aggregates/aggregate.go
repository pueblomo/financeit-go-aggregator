package aggregates

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"main.go/util"
)

var validate = validator.New()

type AggregateDate struct {
	FirstDay string
	LastDay string
}

type AggregateByDate struct {
	Date string `json:"date"`
	Categories []AggregateByCategory `json:"categories"`
}

type AggregateByCategory struct {
	Category CategoryDto `json:"category"`
	Expenses []ExpenseDto `json:"expenses"`
}


func GetAggregatesInMonths(ctx *fiber.Ctx) error{
	log.Println("Get Aggregates")
	page := ctx.Query("page")
	items := ctx.Query("items")
	if page == "" || items == "" {
		log.Println("can't get page or items from path")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	log.Printf("for Page %v with %v items \n", page, items)

	intPage, errPage := strconv.Atoi(page)
	intItems, errItems := strconv.Atoi(items)

	if errPage != nil || errItems != nil {
		log.Printf("error on parsing %v, %v", page, items)
	}

	aggregateDates := getMonths(intPage, intItems)
	result := []AggregateByDate{}
	for _, v := range aggregateDates {
		categories, err := GetCategoriesForMonth(v)
		if err != nil {
			log.Println(err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		aggregatesByCategory := []AggregateByCategory{}

		for _, cat := range categories {
			expenses, err := GetExpensesForMonth(v,cat.Id)
			if err != nil {
				log.Println(err)
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}
			aggregatesByCategory = append(aggregatesByCategory, AggregateByCategory{cat,expenses})
		}
		
		result = append(result, AggregateByDate{v.FirstDay, aggregatesByCategory})
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func getMonths(page int, items int) []AggregateDate {
	months := make([]AggregateDate,items)
	t := time.Now().AddDate(0,-1,0)
	firstday := util.GetDateWithFirstDay(t)
	offset := page * items
	for i:=0; i < items; i++{
		dateWithOffset := firstday.AddDate(0,-i-offset,0)
		firstDayWithOffset := strings.Split(dateWithOffset.String(), " ")
		lastDayWithOffset := strings.Split(util.GetDateWithLastDay(dateWithOffset).String(), " ")
		months[i] = AggregateDate{
			firstDayWithOffset[0] + " " + firstDayWithOffset[1],
			lastDayWithOffset[0] + " " + lastDayWithOffset[1],
		}
	}

	return months
}