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
	Hashrate        float64 `json:"hashrate,omitempty"`
	Power           float64 `json:"power,omitempty"`
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

	app.Get("/api/fetch_usd_to_btc", func(c *fiber.Ctx) error {

		b := GetDollarCourse("https://www.cbr-xml-daily.ru/latest.js")

		btc, _ := GetBitcoinPrice()
		result := struct {
			Result float64 `json:"result"`
		}{Result: (1 / b.Rates.USD) / btc}

		return c.Status(fiber.StatusOK).JSON(result)
	}).Name("api")

	app.Get("/api/fetch_btc_to_usd", func(c *fiber.Ctx) error {

		btc, _ := GetBitcoinPrice()
		result := struct {
			Result float64 `json:"result"`
		}{Result: btc}

		return c.Status(fiber.StatusOK).JSON(result)
	}).Name("api")

	app.Get("/api/fetch_usd_to_rub", func(c *fiber.Ctx) error {
		b := GetDollarCourse("https://www.cbr-xml-daily.ru/latest.js")
		result := struct {
			Result float64 `json:"result"`
		}{Result: 1 / b.Rates.USD}
		return c.Status(fiber.StatusOK).JSON(result)
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
	log.Fatal(app.Listen(":3000"))
	//	csr := []byte(`-----BEGIN CERTIFICATE-----
	//
	// MIIEFDCCAvygAwIBAgISBDmBcV2bXWYqrr6UXguHsVGmMA0GCSqGSIb3DQEBCwUA
	// MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
	// EwJSMzAeFw0yMzExMTkwOTI4MzNaFw0yNDAyMTcwOTI4MzJaMBQxEjAQBgNVBAMT
	// CXVyYWxkYy5ydTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABJhQnMEByee/4fB0
	// w6BJ11vZkhETQNaPmWG+SQFeFVRfjrDRRV0q0auVuo0ITXu/Z5JaYIq0SYAf0nkB
	// HwISItajggILMIICBzAOBgNVHQ8BAf8EBAMCB4AwHQYDVR0lBBYwFAYIKwYBBQUH
	// AwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFLP6PpTGCZrcOJOE
	// 9AgJ76VxL69EMB8GA1UdIwQYMBaAFBQusxe3WFbLrlAJQOYfr52LFMLGMFUGCCsG
	// AQUFBwEBBEkwRzAhBggrBgEFBQcwAYYVaHR0cDovL3IzLm8ubGVuY3Iub3JnMCIG
	// CCsGAQUFBzAChhZodHRwOi8vcjMuaS5sZW5jci5vcmcvMBQGA1UdEQQNMAuCCXVy
	// YWxkYy5ydTATBgNVHSAEDDAKMAgGBmeBDAECATCCAQQGCisGAQQB1nkCBAIEgfUE
	// gfIA8AB2AEiw42vapkc0D+VqAvqdMOscUgHLVt0sgdm7v6s52IRzAAABi+cf9tsA
	// AAQDAEcwRQIgWmNg3+hc1lysOxzqIdAgzGJBF6Y7IDuk76u1ChAXc9cCIQDivzfw
	// 64y7hHPtJPQZYikfBO6J2vAVgrTMhsfacSFkMgB2AO7N0GTV2xrOxVy3nbTNE6Iy
	// h0Z8vOzew1FIWUZxH7WbAAABi+cf9ukAAAQDAEcwRQIhAJ1px1qv8lmA0X0zqVvs
	// Tm7hPlnbT1KLp3MxpkwULjPxAiBdqA3Pgm2FzmI7wwWcN1HJvlVr/SfWhyRQssw6
	// Ud05uTANBgkqhkiG9w0BAQsFAAOCAQEAfxL2jFF97bUX744icLbhm2ZES+nErSNk
	// e5RNYJVRJ1D9loXLMxbwE33ctCTv34jiuAVBWA97kgUJZI0FeI0UTZ3ro+udv0qE
	// cr88X9PkwzMHnf8fjSYhketyK13+0tFsjqSjP/qx6bKX9XBh27gYC8njEJXln2aj
	// rj08RX/cmbioirRsavq62ueO5hCtsGrmlSCo7eDcPm9CmbyMDT/ohsDoZdTeVd94
	// CvOfbqj5MnOTZWMGUji3KcpSBpoxjOA7EsY97QSEtRYTHGhddF2mSjJOMm0+RBMl
	// RKDcBcz5gym5s25egS3z5jCEu1CKWcE6w9eQ/PqwSwbhswBemztoHg==
	// -----END CERTIFICATE-----
	// -----BEGIN CERTIFICATE-----
	// MIIFFjCCAv6gAwIBAgIRAJErCErPDBinU/bWLiWnX1owDQYJKoZIhvcNAQELBQAw
	// TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
	// cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjAwOTA0MDAwMDAw
	// WhcNMjUwOTE1MTYwMDAwWjAyMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg
	// RW5jcnlwdDELMAkGA1UEAxMCUjMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
	// AoIBAQC7AhUozPaglNMPEuyNVZLD+ILxmaZ6QoinXSaqtSu5xUyxr45r+XXIo9cP
	// R5QUVTVXjJ6oojkZ9YI8QqlObvU7wy7bjcCwXPNZOOftz2nwWgsbvsCUJCWH+jdx
	// sxPnHKzhm+/b5DtFUkWWqcFTzjTIUu61ru2P3mBw4qVUq7ZtDpelQDRrK9O8Zutm
	// NHz6a4uPVymZ+DAXXbpyb/uBxa3Shlg9F8fnCbvxK/eG3MHacV3URuPMrSXBiLxg
	// Z3Vms/EY96Jc5lP/Ooi2R6X/ExjqmAl3P51T+c8B5fWmcBcUr2Ok/5mzk53cU6cG
	// /kiFHaFpriV1uxPMUgP17VGhi9sVAgMBAAGjggEIMIIBBDAOBgNVHQ8BAf8EBAMC
	// AYYwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMBIGA1UdEwEB/wQIMAYB
	// Af8CAQAwHQYDVR0OBBYEFBQusxe3WFbLrlAJQOYfr52LFMLGMB8GA1UdIwQYMBaA
	// FHm0WeZ7tuXkAXOACIjIGlj26ZtuMDIGCCsGAQUFBwEBBCYwJDAiBggrBgEFBQcw
	// AoYWaHR0cDovL3gxLmkubGVuY3Iub3JnLzAnBgNVHR8EIDAeMBygGqAYhhZodHRw
	// Oi8veDEuYy5sZW5jci5vcmcvMCIGA1UdIAQbMBkwCAYGZ4EMAQIBMA0GCysGAQQB
	// gt8TAQEBMA0GCSqGSIb3DQEBCwUAA4ICAQCFyk5HPqP3hUSFvNVneLKYY611TR6W
	// PTNlclQtgaDqw+34IL9fzLdwALduO/ZelN7kIJ+m74uyA+eitRY8kc607TkC53wl
	// ikfmZW4/RvTZ8M6UK+5UzhK8jCdLuMGYL6KvzXGRSgi3yLgjewQtCPkIVz6D2QQz
	// CkcheAmCJ8MqyJu5zlzyZMjAvnnAT45tRAxekrsu94sQ4egdRCnbWSDtY7kh+BIm
	// lJNXoB1lBMEKIq4QDUOXoRgffuDghje1WrG9ML+Hbisq/yFOGwXD9RiX8F6sw6W4
	// avAuvDszue5L3sz85K+EC4Y/wFVDNvZo4TYXao6Z0f+lQKc0t8DQYzk1OXVu8rp2
	// yJMC6alLbBfODALZvYH7n7do1AZls4I9d1P4jnkDrQoxB3UqQ9hVl3LEKQ73xF1O
	// yK5GhDDX8oVfGKF5u+decIsH4YaTw7mP3GFxJSqv3+0lUFJoi5Lc5da149p90Ids
	// hCExroL1+7mryIkXPeFM5TgO9r0rvZaBFOvV2z0gp35Z0+L4WPlbuEjN/lxPFin+
	// HlUjr8gRsI3qfJOQFy/9rKIJR0Y/8Omwt/8oTWgy1mdeHmmjk7j1nYsvC9JSQ6Zv
	// MldlTTKB3zhThV1+XWYp6rjd5JW1zbVWEkLNxE7GJThEUG3szgBVGP7pSWTUTsqX
	// nLRbwHOoq7hHwg==
	// -----END CERTIFICATE-----
	// -----BEGIN CERTIFICATE-----
	// MIIFYDCCBEigAwIBAgIQQAF3ITfU6UK47naqPGQKtzANBgkqhkiG9w0BAQsFADA/
	// MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMT
	// DkRTVCBSb290IENBIFgzMB4XDTIxMDEyMDE5MTQwM1oXDTI0MDkzMDE4MTQwM1ow
	// TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
	// cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwggIiMA0GCSqGSIb3DQEB
	// AQUAA4ICDwAwggIKAoICAQCt6CRz9BQ385ueK1coHIe+3LffOJCMbjzmV6B493XC
	// ov71am72AE8o295ohmxEk7axY/0UEmu/H9LqMZshftEzPLpI9d1537O4/xLxIZpL
	// wYqGcWlKZmZsj348cL+tKSIG8+TA5oCu4kuPt5l+lAOf00eXfJlII1PoOK5PCm+D
	// LtFJV4yAdLbaL9A4jXsDcCEbdfIwPPqPrt3aY6vrFk/CjhFLfs8L6P+1dy70sntK
	// 4EwSJQxwjQMpoOFTJOwT2e4ZvxCzSow/iaNhUd6shweU9GNx7C7ib1uYgeGJXDR5
	// bHbvO5BieebbpJovJsXQEOEO3tkQjhb7t/eo98flAgeYjzYIlefiN5YNNnWe+w5y
	// sR2bvAP5SQXYgd0FtCrWQemsAXaVCg/Y39W9Eh81LygXbNKYwagJZHduRze6zqxZ
	// Xmidf3LWicUGQSk+WT7dJvUkyRGnWqNMQB9GoZm1pzpRboY7nn1ypxIFeFntPlF4
	// FQsDj43QLwWyPntKHEtzBRL8xurgUBN8Q5N0s8p0544fAQjQMNRbcTa0B7rBMDBc
	// SLeCO5imfWCKoqMpgsy6vYMEG6KDA0Gh1gXxG8K28Kh8hjtGqEgqiNx2mna/H2ql
	// PRmP6zjzZN7IKw0KKP/32+IVQtQi0Cdd4Xn+GOdwiK1O5tmLOsbdJ1Fu/7xk9TND
	// TwIDAQABo4IBRjCCAUIwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYw
	// SwYIKwYBBQUHAQEEPzA9MDsGCCsGAQUFBzAChi9odHRwOi8vYXBwcy5pZGVudHJ1
	// c3QuY29tL3Jvb3RzL2RzdHJvb3RjYXgzLnA3YzAfBgNVHSMEGDAWgBTEp7Gkeyxx
	// +tvhS5B1/8QVYIWJEDBUBgNVHSAETTBLMAgGBmeBDAECATA/BgsrBgEEAYLfEwEB
	// ATAwMC4GCCsGAQUFBwIBFiJodHRwOi8vY3BzLnJvb3QteDEubGV0c2VuY3J5cHQu
	// b3JnMDwGA1UdHwQ1MDMwMaAvoC2GK2h0dHA6Ly9jcmwuaWRlbnRydXN0LmNvbS9E
	// U1RST09UQ0FYM0NSTC5jcmwwHQYDVR0OBBYEFHm0WeZ7tuXkAXOACIjIGlj26Ztu
	// MA0GCSqGSIb3DQEBCwUAA4IBAQAKcwBslm7/DlLQrt2M51oGrS+o44+/yQoDFVDC
	// 5WxCu2+b9LRPwkSICHXM6webFGJueN7sJ7o5XPWioW5WlHAQU7G75K/QosMrAdSW
	// 9MUgNTP52GE24HGNtLi1qoJFlcDyqSMo59ahy2cI2qBDLKobkx/J3vWraV0T9VuG
	// WCLKTVXkcGdtwlfFRjlBz4pYg1htmf5X6DYO8A4jqv2Il9DjXA6USbW1FzXSLr9O
	// he8Y4IWS6wY7bCkjCWDcRQJMEhg76fsO3txE+FiYruq9RUWhiF1myv4Q6W+CyBFC
	// Dfvp7OOGAN6dEOM4+qR9sdjoSYKEBpsr6GtPAQw4dy753ec5
	// -----END CERTIFICATE-----`)
	//
	//	key := []byte(`-----BEGIN PRIVATE KEY-----
	//
	// MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgY0wn0e0wmhj+BZUk
	// L7byAH4WzBS9giSJIPu0xyuFHAOhRANCAASYUJzBAcnnv+HwdMOgSddb2ZIRE0DW
	// j5lhvkkBXhVUX46w0UVdKtGrlbqNCE17v2eSWmCKtEmAH9J5AR8CEiLW
	// -----END PRIVATE KEY-----`)
	//
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
	var result struct {
		Result float64 `json:"result"`
	}
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
	if p.DeviceID != 0 {
		dev, err := repo.GetDevice(context.Background(), p.DeviceID)
		if err != nil {
			return err
		}
		btc, _ := GetBitcoinPrice()
		f := p.ElectricityCost

		s := ((900/float64(x))*float64(dev.Hashrate)*float64(p.DaysWork))*btc - ((float64(dev.Power) / 1000) * float64(p.DaysWork) * 24 * (float64(f) / (1 / b.Rates.USD)))

		result.Result = s
	} else {
		btc, _ := GetBitcoinPrice()
		f := p.ElectricityCost

		s := ((900/float64(x))*float64(p.Hashrate)*float64(p.DaysWork))*btc - ((float64(p.Power) / 1000) * float64(p.DaysWork) * 24 * (float64(f) / (1 / b.Rates.USD)))

		result.Result = s
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
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
