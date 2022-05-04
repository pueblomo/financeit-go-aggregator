package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"main.go/conf"
)

var Conn *pgxpool.Pool

var createCat = "CREATE TABLE IF NOT EXISTS category (id serial PRIMARY KEY,categoryid int, name text, icon text, valid_from timestamp, valid_to timestamp)"
var createExp = "CREATE TABLE IF NOT EXISTS expense (id serial PRIMARY KEY, description text, date timestamp, price int, categoryid int)"

func Connect() {
	var err error
	pw := conf.GetDBPw()
	url := conf.GetDBUrl()
	conString := "postgresql://postgres:"+pw+"@"+url+":5432/?sslmode=disable"
	Conn,err = pgxpool.Connect(context.Background(), conString)
	conf.CheckErrorFatal(err)

	initializeDb()
}

func initializeDb() {
	_, err := Conn.Exec(context.Background(),createCat)
	conf.CheckErrorFatal(err)

	_, err = Conn.Exec(context.Background(), createExp)
	conf.CheckErrorFatal(err)
}

