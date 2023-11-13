package app

import (
	"BTCcalc/internal/repo"
	"BTCcalc/internal/usecase"
	"context"
	"crypto/tls"
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

	db, err := sqlx.Open("mysql", "root:dCmd5e5A6hUN8Yv@(193.109.84.90:3306)/leomine_schema")
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
	csr := []byte(`-----BEGIN CERTIFICATE-----
MIIFOzCCAyOgAwIBAgIURRqN5q8eXZ0yRsSErpBv5v8H3HswDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wHhcNMjMxMTEzMTkxNjA3WhcNMzMx
MTEwMTkxNjA3WjAWMRQwEgYDVQQDDAtleGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcN
AQEBBQADggIPADCCAgoCggIBAPe1GjioJMmQWNDmQUubI2RYhpcyDm7ZMpcae459
Ws9FXQhyplAIl5V73p8VxJK2xY3/1KZ5gGr7SZ5n5d+552oIXygz5ojkLj0rkxpr
3CCpjo3APmAFoQ+Nj8VGSx4mDh/+27hET70L3/s6ZNBoEdqrrfd20hYQ9O2FsdnH
r3traMohkTldLyT67L+mVgGewP80ktUfmlFb87k5shT0Thps4sPzMiHkmU+FNhM1
gfWgWzEe+1DOE3JUtyOF9++nXZuGWP3oGv+l/P8M4OGeL9UA1//6miHATV7xpV8P
f6JnPKth1qJjQPF3xc24cSyFuLBOEY0z2q7AGB5AoUzMPo7yHRxjsYjMWiqcfzfb
WWCTY/95SQO23Xiy0IFcLJpgBQIIkRzVJD7KpgZxvcUZDoL5vLSk1f/zKKX6R2l7
mtFIxD7j3trUnXErEnPH9vv+X7o83zkUbvV5yiC9J+xAj25VKznAS9PPzzgJl1BR
gixvBsK3EPnmYz3DddWDuSjK8HX4M5o5mj7VLgIzml+jEFDW8MIKoDgtiWmzCKgG
Z4DKzH7ICRw+tcpyGke8p+vC+cbAlFb98nsAV5GtDxKTYDf+1eF0bNyYjLS5Al3U
UPsAHrutsOSbwmiodZqx/Dr48ZCGHzDMCqIImVL5dEfoLCeU2dNYk+L9IAPm4JzQ
xgNRAgMBAAGjgYAwfjAdBgNVHQ4EFgQUrASqx2GusQrW2cevsqBid59hKwMwHwYD
VR0jBBgwFoAUrASqx2GusQrW2cevsqBid59hKwMwDwYDVR0TAQH/BAUwAwEB/zAr
BgNVHREEJDAiggtleGFtcGxlLmNvbYINKi5leGFtcGxlLmNvbYcECgAAATANBgkq
hkiG9w0BAQsFAAOCAgEA8vqiDEcnTrP8CkdUk7vWdmNEocInGHducQiP2sscgzmr
3MIqE/0DJihfmWbWZGvfP++P5twtO/6MX0OBDf8VNw1z9AUYMaiApZVkudb7YozE
T6AAAqfzFmaLTEeX08eynN2UV/ijffAAWQYTGWUL5e4dPRBaT8JIl+fewxBEJflz
PIvaX/61wAiqa17F0g+D9G3n0qzEdufWxDjYO/hsM5aI8j1qGN8goZBqP+VYAUBh
eMZVp6y55+Us+oO/z2ADmqAYIFJJXvwqAHMeyPKyJewd+tkE8xYePOmKjE4NCQ+9
oKqEtgyDr2JGNnf+wqXORnr1NOAXvQWu+OuX2eJlfRHC7ug8Had3QnvfEoUtY7Xa
gwRWe9g2aW3UKiOb1Lo3oKAlldoIIzBX5IWof9YopjVR2Nqsc1IFyWM1nw5xHzP5
uvpmjR1t2yb118wsQgkV4kDmdEJFEsvxWIsIV6L9EdaUyn8yEfAR4zdxIAKcEy3S
pgslfLDm9i6T+KDZIer2N8eP8qTR0wpkz38bbE16Dw0BYdSiUcHQfjC4EBHYLDFN
eF9xlqeDrFRxvNQfjqKYtxxKwt77+uYoRsOig9SP7zKK80kGGCnugKctOFXvYyB7
MC2AGxWYuD/xzHbiFsQ3HO/yZ9ZerQmkL7z7Q7cqSa131LqlEgzRufQcW9NnRiw=
-----END CERTIFICATE-----`)
	key := []byte(`-----BEGIN PRIVATE KEY-----
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
WPEqaNHG/ch/4OAnOg0yKtJOzMfHFfe91g==
-----END PRIVATE KEY-----`)
	cert, err := tls.X509KeyPair(csr, key)
	log.Fatal(app.ListenTLSWithCertificate(":80", cert))
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
