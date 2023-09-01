package main

import (
	app "BTCcalc/internal/app"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	app.Run()

}
