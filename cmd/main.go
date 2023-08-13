package main

import (
	"BTCcalc/repo"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type CalculateDTO struct {
	DeviceID        int     `json:"deviceID"`
	DaysWork        int     `json:"daysWork"`
	ElectricityCost float64 `json:"electricityCost"`
}
type countryRub struct {
	USD float64 `json:"USD"`
}

type dollar struct {
	Rates countryRub `json:"rates"`
}

type values struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type hashValues struct {
	Values []values `json:"values"`
}

func main() {
	app := fiber.New()

	db, err := sqlx.Open("mysql", "root:dCmd5e5A6hUN8Yv@(193.109.84.90:3306)/leomine_schema")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	repo := repo.NewRepository(db)

	// GET /api/register
	app.Post("/api/calculate", func(c *fiber.Ctx) error {
		return calculate(c, repo)
	}).Name("api")

	data, _ := json.MarshalIndent(app.GetRoute("api"), "", "  ")
	fmt.Print(string(data))

	log.Fatal(app.Listen(":3000"))

}

func GetBTCHashrate(some string) (b hashValues) {
	resp, err := http.Get(some)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	b = hashValues{}

	json.Unmarshal(body, &b)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func calculate(ctx *fiber.Ctx, repo *repo.Repository) error {
	var p CalculateDTO
	if err := ctx.BodyParser(&p); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON((err.Error()))
	}

	b := GetDollarCourse("https://www.cbr-xml-daily.ru/latest.js")
	j := fmt.Sprintf("https://api.blockchain.info/charts/hash-rate?timespan=%ddays&format=json", p.DaysWork)
	h := GetBTCHashrate(j)
	x := 0.0
	for _, i := range h.Values {
		x += math.Floor(i.Y)
	}
	x = x / float64(p.DaysWork)

	dev, err := repo.GetDevices(context.Background(), p.DeviceID)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(dev.Hashrate, 64)
	if err != nil {
		return err
	}
	btc, _ := GetBitcoinPrice()
	pow, err := strconv.ParseFloat(dev.Power, 64)
	if err != nil {
		return err
	}
	f := p.ElectricityCost

	s := ((900/float64(x))*float64(y)*float64(p.DaysWork))*btc - ((float64(pow) / 1000) * float64(p.DaysWork) * 24 * (float64(f) / (1 / b.Rates.USD)))

	return ctx.Status(fiber.StatusOK).JSON((s))
}

func GetBitcoinPrice() (price float64, err error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://blockchain.info/tobtc?currency=USD&value=500", nil)
	if err != nil {
		fmt.Printf("Got error %s", err.Error())
		return
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Got error %s", err.Error())
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	s, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return
	}
	price = 500 / s
	return
}

func GetDollarCourse(some string) (b dollar) {
	resp, err := http.Get(some)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	b = dollar{}

	json.Unmarshal(body, &b)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
