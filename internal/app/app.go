package app

import (
	"BTCcalc/internal/repo"
	"BTCcalc/internal/usecase"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	null "gopkg.in/guregu/null.v3/zero"
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

type DeviceDTO struct {
	DeviceID     null.Int   `json:"deviceID,omitempty"`
	PriceLow     null.Float `json:"priceLow,omitempty"`
	PriceHigh    null.Float `json:"priceHigh,omitempty"`
	PowerLow     null.Float `json:"powerLow,omitempty"`
	PowerHigh    null.Float `json:"powerHigh,omitempty"`
	HashrateLow  null.Float `json:"hashrateLow,omitempty"`
	HashrateHigh null.Float `json:"hashrateHigh,omitempty"`
	HashrateID   null.Int   `json:"hashrateID,omitempty"`
	BrandID      null.Int   `json:"brandID,omitempty"`
	OfferID      null.Int   `json:"offerID,omitempty"`
	CoinID       null.Int   `json:"coinID,omitempty"`
	Recommended  null.Int   `json:"recommended,omitempty"`
}

type DeviceImageDTO struct {
	DeviceID []null.Int `json:"deviceID,omitempty"`
}

type ArticleImageDTO struct {
	ArticleID []null.Int `json:"articleID,omitempty"`
}

type CaseImageDTO struct {
	CaseID []null.Int `json:"caseID,omitempty"`
}

type ReviewDTO struct {
	DeviceID null.Int `json:"deviceID,omitempty"`
	ReviewID int      `json:"reviewID,omitempty"`
}
type ArticleDTO struct {
	ArticleID null.Int `json:"articleID,omitempty"`
}

type CaseDTO struct {
	CaseID null.Int `json:"CaseID,omitempty"`
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

func Run() {
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

	uc := usecase.NewUsecase(repo)
	// GET /api/register
	app.Post("/api/calculate", func(c *fiber.Ctx) error {
		return calculate(c, repo)
	}).Name("api")

	app.Post("/api/get_device", func(c *fiber.Ctx) error {
		var p DeviceDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		some := usecase.DeviceDTO{DeviceID: p.DeviceID.NullInt64, PriceLow: p.PriceLow.NullFloat64, PriceHigh: p.PriceHigh.NullFloat64, PowerLow: p.PowerLow.NullFloat64,
			PowerHigh: p.PowerHigh.NullFloat64, HashrateHigh: p.HashrateHigh.NullFloat64, HashrateLow: p.HashrateLow.NullFloat64,
			HashrateID: p.HashrateID.NullInt64, BrandID: p.BrandID.NullInt64, OfferID: p.OfferID.NullInt64, CoinID: p.CoinID.NullInt64, Recommended: p.Recommended.NullInt64}

		result, err := uc.GetDevices(c.Context(), some)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")
	app.Post("/api/get_device_image", func(c *fiber.Ctx) error {
		var p DeviceImageDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		some := make([]sql.NullInt64, 1)
		for _, v := range p.DeviceID {
			some = append(some, v.NullInt64)
		}
		d := usecase.DeviceImageDTO{DeviceID: some}

		result, err := uc.GetDeviceImage(c.Context(), d)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Post("/api/get_article_image", func(c *fiber.Ctx) error {
		var p ArticleImageDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		some := make([]sql.NullInt64, 1)
		for _, v := range p.ArticleID {
			some = append(some, v.NullInt64)
		}
		d := usecase.ArticleImageDTO{ArticleID: some}

		result, err := uc.GetArticleImage(c.Context(), d)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Post("/api/get_case_image", func(c *fiber.Ctx) error {
		var p CaseImageDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		some := make([]sql.NullInt64, 1)
		for _, v := range p.CaseID {
			some = append(some, v.NullInt64)
		}
		d := usecase.CaseImageDTO{CaseID: some}

		result, err := uc.GetCaseImage(c.Context(), d)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Post("/api/get_device_review", func(c *fiber.Ctx) error {
		var p ReviewDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		result, err := uc.GetReviews(c.Context(), p.DeviceID)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Post("/api/get_article", func(c *fiber.Ctx) error {
		var p ArticleDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		result, err := uc.GetArticles(c.Context(), p.ArticleID)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Post("/api/get_case", func(c *fiber.Ctx) error {
		var p CaseDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		result, err := uc.GetCases(c.Context(), p.CaseID)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Get("/api/fetch_brands", func(c *fiber.Ctx) error {
		result, err := uc.GetBrand(c.Context())
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Get("/api/fetch_offers", func(c *fiber.Ctx) error {
		result, err := uc.GetOffer(c.Context())
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")
	app.Get("/api/fetch_coins", func(c *fiber.Ctx) error {
		result, err := uc.GetCoin(c.Context())
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Get("/api/fetch_hashrate", func(c *fiber.Ctx) error {
		result, err := uc.GetHashrate(c.Context())
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
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

	dev, err := repo.GetDevice(context.Background(), p.DeviceID)
	if err != nil {
		return err
	}
	btc, _ := GetBitcoinPrice()
	f := p.ElectricityCost

	s := ((900/float64(x))*float64(dev.Hashrate)*float64(p.DaysWork))*btc - ((float64(dev.Power) / 1000) * float64(p.DaysWork) * 24 * (float64(f) / (1 / b.Rates.USD)))

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
