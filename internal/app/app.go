package app

import (
	"BTCcalc/internal/repo"
	"BTCcalc/internal/usecase"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	DeviceID     *int64   `json:"deviceID,omitempty"`
	PriceLow     *float64 `json:"priceLow,omitempty"`
	PriceHigh    *float64 `json:"priceHigh,omitempty"`
	PowerLow     *float64 `json:"powerLow,omitempty"`
	PowerHigh    *float64 `json:"powerHigh,omitempty"`
	HashrateLow  *float64 `json:"hashrateLow,omitempty"`
	HashrateHigh *float64 `json:"hashrateHigh,omitempty"`
	HashrateID   []*int64 `json:"hashrateID,omitempty"`
	BrandID      []*int64 `json:"brandID,omitempty"`
	OfferID      []*int64 `json:"offerID,omitempty"`
	CoinID       []*int64 `json:"coinID,omitempty"`
	Recommended  *int64   `json:"recommended,omitempty"`
}

type DeviceImageDTO struct {
	DeviceID []*int64 `query:"deviceID,omitempty"`
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

type GetCallDTO struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

type WriteReviewDTO struct {
	Stars       int64  `json:"stars"`
	Text        string `json:"text"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type WriteDeviceReviewDTO struct {
	DeviceID    int64  `json:"deviceID"`
	Stars       int64  `json:"stars"`
	Amount      int64  `json:"amount"`
	Text        string `json:"text"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
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

	app.Use(cors.New())

	db, err := sqlx.Open("mysql", "u2333338_root:Y3S4G9taZvwsYtU2@(31.31.196.165:3306)/u2333338_some")
	if err != nil {
		log.Fatal(err.Error())
	}
	db.SetMaxIdleConns(0)
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

	app.Get("/api/get_device", func(c *fiber.Ctx) error {
		var p DeviceDTO
		if err := c.QueryParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		some := usecase.DeviceDTO{DeviceID: p.DeviceID, PriceLow: p.PriceLow, PriceHigh: p.PriceHigh, PowerLow: p.PowerLow,
			PowerHigh: p.PowerHigh, HashrateHigh: p.HashrateHigh, HashrateLow: p.HashrateLow,
			HashrateID: p.HashrateID, BrandID: p.BrandID, OfferID: p.OfferID, CoinID: p.CoinID, Recommended: p.Recommended}

		result, err := uc.GetDevices(c.Context(), some)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Get("/api/get_powerful_device", func(c *fiber.Ctx) error {
		result, err := uc.GetPowerfulDevices(c.Context())
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Get("/api/get_device_image", func(c *fiber.Ctx) error {
		var p DeviceImageDTO
		if err := c.QueryParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		d := usecase.DeviceImageDTO{DeviceID: p.DeviceID}

		result, err := uc.GetDeviceImage(c.Context(), d)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Get("/api/get_article_image", func(c *fiber.Ctx) error {
		var p ArticleImageDTO
		if err := c.QueryParser(&p); err != nil {
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

	app.Get("/api/get_case_image", func(c *fiber.Ctx) error {
		var p CaseImageDTO
		if err := c.QueryParser(&p); err != nil {
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

	app.Get("/api/get_device_review", func(c *fiber.Ctx) error {
		var p ReviewDTO
		if err := c.QueryParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		result, err := uc.GetReviews(c.Context(), p.DeviceID)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Get("/api/get_call", func(c *fiber.Ctx) error {
		var p GetCallDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		_, err := uc.GetCall(c.Context(), p.Name, p.PhoneNumber)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	}).Name("api")

	app.Post("/api/write_review", func(c *fiber.Ctx) error {
		var p WriteReviewDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		err := uc.WriteReview(c.Context(), p.Email, p.PhoneNumber, p.Text, p.Stars)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	}).Name("api")

	app.Post("/api/write_device_review", func(c *fiber.Ctx) error {
		var p WriteDeviceReviewDTO
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		err := uc.WriteDeviceReview(c.Context(), p.Email, p.PhoneNumber, p.Text, p.DeviceID, p.Stars, p.Amount)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	}).Name("api")

	app.Get("/api/get_article", func(c *fiber.Ctx) error {
		var p ArticleDTO
		if err := c.QueryParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON((err.Error()))
		}
		result, err := uc.GetArticles(c.Context(), p.ArticleID)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON((result))
	}).Name("api")

	app.Get("/api/get_case", func(c *fiber.Ctx) error {
		var p CaseDTO
		if err := c.QueryParser(&p); err != nil {
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
	log.Fatal(app.Listen(":80"))
	//	csr := []byte(`-----BEGIN CERTIFICATE-----
	//MIIFLTCCAxWgAwIBAgIUczzQqI/P9Kt7dnvfEzn4QDMWgCQwDQYJKoZIhvcNAQEL
	//BQAwEzERMA8GA1UEAwwIYWxvZGEucnUwHhcNMjMxMTE1MTkyMDA0WhcNMzMxMTEy
	//MTkyMDA0WjATMREwDwYDVQQDDAhhbG9kYS5ydTCCAiIwDQYJKoZIhvcNAQEBBQAD
	//ggIPADCCAgoCggIBAL6cfBGqK7qMhESQzeawZfGIHHc9lZx/wKxqoJmLPQdGTbKV
	//L+9m8DYOP3o53088mQ58UyYN9A27tfWYqDWjJW55wcUl2e10OlW4ox7T5jvM95mI
	//QFLopPp1WaMR6DJuackuAGmIK7qWi9xPwJyao35tKoOzVeah23dJpf/tMc6rgOqx
	//OIB3b1Wr/yXnNgnVAdWRIL0Kufrf3M4TuP+mfHdoAXn8wcf/IoPl4+BoIXmn7XVB
	//mGv2hTh/W59mcPGVC0UKIb6kjhPl84SepyBPoT1XjPA7HV7yqYqFIbv/YGmmgxaa
	//qr7aV68xgB89Gf8PQdYx53vz/1yE6BdVDX4BDupl4I1MNFvs1qH+y+pUJ79fi39D
	//1liNoL4bcije2kghebA82MmeQsLNc8OZp/Wq/QJuhM91/QYXI0ERFbkc/edmTfGn
	//RZLHC8L8aJ7YlZYDzI7DHlaEbyO6uDHimNU75KpTz/sS6Z76KausEQ4DmA06doc8
	//n1eEvxGSzwV4ohH2ny4IJRR4po10U76Fm19/0F3VME935Fz7hVCozMzOe2Utzmtx
	//at7A4yj1bT/9/A6GRWd6vdkyraQerpa9L9cx/CkbCuPq3+vmszy3AP0fbj1mt/ok
	//ksJwUaVlrr6SJatayrSnGDKVMyydJcPrFasBDvLp4NXGQglNydMuB70EOrEPAgMB
	//AAGjeTB3MB0GA1UdDgQWBBQY26NT/ROqMTEAeTd+gkPlu4E/tjAfBgNVHSMEGDAW
	//gBQY26NT/ROqMTEAeTd+gkPlu4E/tjAPBgNVHRMBAf8EBTADAQH/MCQGA1UdEQQd
	//MBuCCGFsb2RhLnJ1ggkqYWxvZGEucnWHBMFtVFowDQYJKoZIhvcNAQELBQADggIB
	//AAPHGCHhucIOGGr5n1LnMPWDHrlzZouDUguVM5YuhI3522wfe3vx6ntdkcIFZ48C
	//0DyqfDDUHTLHewMrgprb2Jo+OEhBf+69bn7GlUYtDrOe3M0J3VjdawRkAdGxy3EY
	//YZQ/81KsG9yu7YilvpMnNfalMWdD5qDd46ZsJz7cQTIV/oW8httCAC0VjfifEmNa
	//d22qd1cVpxGCCphbFAootEV4imbY0lAC3aFmwNoajTBzL3+SlyTYs7E1svj6Dvyy
	//kMFnEkrUB1L0R7LzfNw4OxPBrnuKIzXzuYCykdS33l/W1MKtYNTUPVcXFEVPtSWf
	//DfJLSVEdFUOuNF0xlcsE+X9BfSvWRzB0odW9iuLrzmSG6bVJM0r6GphffWc5e172
	//jK9Q7S3cXhGlTiVBbbDKKtcQhbp1zlsuTXiFH/uKHJOLfZD6hKrrTYf/W5sIh2Zq
	//eSL0W4KznOjm0HEuxIOgLGYTIR5JjfOV5zYu5omRgkJtqYPPWMgkzbT7oBSLAIiK
	//to/1yk+UmerfgcYQW+08QxF4Ln59OBdPG+shw4/V5c2oM9SuCnlIjnVGmRDgHkjn
	//q/gBf662DD5DUMJRR6Eds7I2omctUPNScxE7FjwHSb/nmt0F4SDUalKZIbrR7oZn
	//EW5istNSfEn/KSU/0q4cFNJNsbluTzLTr7bahvimLc0E
	//-----END CERTIFICATE-----`)
	//	key := []byte(`-----BEGIN PRIVATE KEY-----
	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
	//SL++7m2xTlSJtLIi0U9f4t4B16YZqg==
	//-----END PRIVATE KEY-----`)
	//	cert, err := tls.X509KeyPair(csr, key)
	//	log.Fatal(app.ListenTLSWithCertificate(":80", cert))
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
