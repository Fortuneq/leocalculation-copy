package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	null "gopkg.in/guregu/null.v3/zero"
	"strings"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

type Device struct {
	Name        string  `db:"name"`
	Id          int64   `db:"id"`
	Cost        float64 `db:"cost"`
	Size        string  `db:"size"`
	Power       float64 `db:"power"`
	Quantity    int     `json:"quantity"`
	Hashrate    float64 `db:"hashrate"`
	UID         string  `db:"uid"`
	Algorithm   string  `db:"algorithm"`
	VideoUrl    string  `db:"video_url"`
	CoinName    string  `db:"coin_name"`
	HashName    string  `db:"hash_name"`
	BrandName   string  `db:"brand_name"`
	OfferName   string  `db:"offer_name"`
	Recommended int     `db:"recommended"`
}

type DeviceImage struct {
	ID    int64  `db:"id"`
	Image string `db:"image"`
}
type CaseImage struct {
	ID int64 `db:"id"`

	Image string `db:"image"`
}

type ArticleImage struct {
	ID    int64  `db:"id"`
	Image string `db:"image"`
}

type Brand struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Coin struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type HashrateType struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type OfferType struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type DeviceReviews struct {
	Email    string  `db:"email"`
	Text     string  `db:"review_text"`
	Name     string  `db:"name"`
	Stars    int64   `db:"stars"`
	Id       int64   `db:"id"`
	DeviceId int64   `db:"device_id"`
	Amount   float64 `db:"amount"`
	Date     string  `db:"date"`
}

type Articles struct {
	Text     string `db:"text"`
	Name     string `db:"name"`
	Id       int64  `db:"id"`
	VideoUrl string `db:"video_url"`
}

type Cases struct {
	Text     string `db:"text"`
	Name     string `db:"name"`
	Id       int64  `db:"id"`
	VideoUrl string `db:"video_url"`
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
	DeviceID []*int64 `json:"deviceID,omitempty"`
}

type ArticleImageDTO struct {
	ArticleID []sql.NullInt64 `json:"articleID,omitempty"`
}

type CaseImageDTO struct {
	CaseID []sql.NullInt64 `json:"caseID,omitempty"`
}

func (r *Repository) GetDevice(ctx context.Context, id int) (Device, error) {

	//Абстрактный sql ,  с которого получаем данные

	q := "SELECT id,name, cost,image,size,power,hashrate,algorithm,uid,video_url FROM devices  where id = ?"

	place := Device{}

	err := r.db.GetContext(ctx, &place, q, id)
	if err != nil {
		return Device{}, err
	}
	return place, nil
}

func (r *Repository) GetDevices(ctx context.Context, p DeviceDTO) (result []Device, err error) {
	//Абстрактный sql ,  с которого получаем данные
	q := `
SELECT DISTINCT devices.id, devices.name AS name, cost, size, power, hashrate, algorithm, uid, video_url, c.name AS coin_name,
	h.name AS hash_name, ot.name AS offer_name, recommended, dp.name AS brand_name
FROM devices
JOIN device_coin dc ON devices.id = dc.device_id
JOIN coins c ON dc.coin_id = c.id
JOIN device_producers dp ON dp.id = devices.producer_id
JOIN hashrate h ON h.id = devices.hashrate_id
JOIN offer_types ot ON devices.offer_type = ot.id
WHERE 1=1 
`
	if p.DeviceID != nil {
		q += " AND   (? IS NULL OR devices.id = ?)\n"
	}
	if p.PriceLow != nil {
		q += " AND   (? IS NULL OR cost >= ?)\n"
	}
	if p.PowerLow != nil {
		q += " AND   (? IS NULL OR power >= ?)\n"
	}
	if p.PriceHigh != nil {
		q += " AND   (? IS NULL OR cost <= ?)\n"
	}
	if p.PowerHigh != nil {
		q += "  AND  (? IS NULL OR power <= ?)\n"
	}
	if p.HashrateLow != nil {
		q += " AND   (? IS NULL OR hashrate >= ?)\n"
	}
	if p.HashrateHigh != nil {
		q += " AND   (? IS NULL OR hashrate <= ?)\n"
	}
	if p.Recommended != nil {
		q += "  AND  (? IS NULL OR recommended = ?)\n"
	}

	params := make([]interface{}, 0)
	if p.DeviceID != nil {
		params = append(params, *p.DeviceID)
		params = append(params, *p.DeviceID)
	}
	if p.PriceLow != nil {
		params = append(params, *p.PriceLow)
		params = append(params, *p.PriceLow)
	}
	if p.PriceHigh != nil {
		params = append(params, *p.PriceHigh)
		params = append(params, *p.PriceHigh)
	}
	if p.PowerLow != nil {
		params = append(params, *p.PowerLow)
		params = append(params, *p.PowerLow)
	}
	if p.PowerHigh != nil {
		params = append(params, *p.PowerHigh)
		params = append(params, *p.PowerHigh)
	}
	if p.HashrateLow != nil {
		params = append(params, *p.HashrateLow)
		params = append(params, *p.HashrateLow)
	}
	if p.HashrateHigh != nil {
		params = append(params, *p.HashrateHigh)
		params = append(params, *p.HashrateHigh)
	}
	if p.Recommended != nil {
		params = append(params, *p.Recommended)
		params = append(params, *p.Recommended)
	}

	var coinIDs []string
	var offerIDs []string
	var hashIDs []string
	var brandIDs []string

	// Формируем значения и подстановки из поля CoinID структуры
	for _, v := range p.CoinID {
		placeholder := fmt.Sprintf("%d", *v)
		coinIDs = append(coinIDs, placeholder)
		//params = append(params, *coinID)
	}
	if p.CoinID != nil && len(p.CoinID) > 0 {
		q += fmt.Sprintf("AND ((COALESCE(%s IS NULL, 1) OR c.id IN (%s)))\n", strings.Join(coinIDs, ","), strings.Join(coinIDs, ","))
	}

	for _, v := range p.HashrateID {
		placeholder := fmt.Sprintf("%d", *v)
		hashIDs = append(hashIDs, placeholder)
		//params = append(params, *coinID)
	}
	if p.HashrateID != nil && len(p.HashrateID) > 0 {
		q += fmt.Sprintf("AND ((COALESCE(%s IS NULL, 1) OR h.id IN (%s)))\n", strings.Join(hashIDs, ","), strings.Join(hashIDs, ","))
	}

	for _, v := range p.OfferID {
		placeholder := fmt.Sprintf("%d", *v)
		offerIDs = append(offerIDs, placeholder)
		//params = append(params, *coinID)
	}
	if p.OfferID != nil && len(p.OfferID) > 0 {
		q += fmt.Sprintf("AND ((COALESCE(%s IS NULL, 1) OR ot.id IN (%s)))\n", strings.Join(offerIDs, ","), strings.Join(offerIDs, ","))
	}

	for _, v := range p.BrandID {
		placeholder := fmt.Sprintf("%d", *v)
		brandIDs = append(brandIDs, placeholder)
		//params = append(params, *coinID)
	}
	if p.BrandID != nil && len(p.BrandID) > 0 {
		q += fmt.Sprintf("AND ((COALESCE(%s IS NULL, 1) OR dp.id IN (%s)))\n", strings.Join(brandIDs, ","), strings.Join(brandIDs, ","))
	}

	err = r.db.SelectContext(ctx, &result, q, params...)
	if err != nil {
		return []Device{}, err
	}

	return result, nil
}

func (r *Repository) GetPowerfulDevices(ctx context.Context) (result []Device, err error) {
	q := "SELECT DISTINCT devices.id,devices.name as name, cost,size,power,hashrate,algorithm,uid,video_url,c.name as coin_name,     " +
		"           h.name as hash_name,ot.name as offer_name,recommended,dp.name as brand_name FROM devices    JOIN device_coin dc on devices.id = dc.device_id  " +
		"  join coins c on dc.coin_id = c.id join device_producers dp on dp.id = devices.producer_id  " +
		"  join hashrate h on h.id = devices.hashrate_id   " +
		" join offer_types ot on devices.offer_type  ORDER BY  power  DESC LIMIT 6 "

	err = r.db.SelectContext(ctx, &result, q)
	if err != nil {
		return []Device{}, err
	}

	return result, nil
}

func (r *Repository) GetDeviceImage(ctx context.Context, p DeviceImageDTO) (result []DeviceImage, err error) {
	args := make([]interface{}, len(p.DeviceID))
	for i, id := range p.DeviceID {
		args[i] = id
	}
	var deviceIDS []string
	if p.DeviceID == nil {
		stmt := `SELECT id,image from devices`
		err = r.db.SelectContext(ctx, &result, stmt)
		if err != nil {
			return []DeviceImage{}, err
		}
		return result, err
	}
	for _, v := range p.DeviceID {
		placeholder := fmt.Sprintf("%d", *v)
		deviceIDS = append(deviceIDS, placeholder)
		//params = append(params, *coinID)
	}
	stmt := fmt.Sprintf("SELECT id,image from devices where id in(%s)", strings.Join(deviceIDS, ","))
	err = r.db.SelectContext(ctx, &result, stmt)
	if err != nil {
		return []DeviceImage{}, err
	}
	return result, nil
}

func (r *Repository) GetArticleImage(ctx context.Context, p ArticleImageDTO) (result []ArticleImage, err error) {
	args := make([]interface{}, len(p.ArticleID))
	for i, id := range p.ArticleID {
		args[i] = id
	}
	stmt := `SELECT id,image from articles where id in(?` + strings.Repeat(",?", len(args)-1) + `)`
	err = r.db.SelectContext(ctx, &result, stmt, args...)
	if err != nil {
		return []ArticleImage{}, err
	}
	return result, nil
}

func (r *Repository) GetCaseImages(ctx context.Context, p CaseImageDTO) (result []CaseImage, err error) {
	args := make([]interface{}, len(p.CaseID))
	for i, id := range p.CaseID {
		args[i] = id
	}
	stmt := `SELECT id,image from cases where id in(?` + strings.Repeat(",?", len(args)-1) + `)`
	err = r.db.SelectContext(ctx, &result, stmt, args...)
	if err != nil {
		return []CaseImage{}, err
	}
	return result, nil
}

func (r *Repository) GetCaseImage(ctx context.Context, id int) (result CaseImage, err error) {
	q := "SELECT image from cases where id = ?"
	err = r.db.GetContext(ctx, &result, q, id)
	if err != nil {
		return CaseImage{}, err
	}
	return result, nil
}

func (r *Repository) GetDeviceReviews(ctx context.Context, id null.Int) (result []DeviceReviews, err error) {
	q := ""
	//Абстрактный sql ,  с которого получаем данные

	q = "SELECT id,name,date,amount,device_id,stars,amount,email,review_text FROM device_reviews where ? is null or id = ?"

	err = r.db.SelectContext(ctx, &result, q, id, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (r *Repository) InsertReview(ctx context.Context, email, text, phone string, stars int64) (err error) {
	q := ""
	//Абстрактный sql ,  с которого получаем данные

	q = "INSERT INTO company_reviews(email,text,name,stars) VALUES (?,?,?,?) "

	_, err = r.db.ExecContext(ctx, q, email, text, phone, stars)
	if err != nil {
		return err
	}
	return err
}

func (r *Repository) InsertDeviceReview(ctx context.Context, email, text, phone string, deviceID, stars, amount int64) (err error) {
	q := ""
	//Абстрактный sql ,  с которого получаем данные

	q = "INSERT INTO device_reviews(email,review_text,name,device_id,stars,amount,date) VALUES (?,?,?,?,?,?,?)"

	_, err = r.db.ExecContext(ctx, q, email, text, phone, deviceID, stars, amount, time.Now())
	if err != nil {
		return err
	}
	return err
}
func (r *Repository) GetArticles(ctx context.Context, id null.Int) (result []Articles, err error) {

	//Абстрактный sql ,  с которого получаем данные
	q := ""
	q = "SELECT id,name,text,video_url FROM articles where( ? is null or id = ?)"

	err = r.db.SelectContext(ctx, &result, q, id, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetCases(ctx context.Context, id null.Int) (result []Cases, err error) {

	//Абстрактный sql ,  с которого получаем данные
	q := ""
	q = "SELECT id,name,text,video_url FROM cases where( ? is null or id = ?)"

	err = r.db.SelectContext(ctx, &result, q, id, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetBrands(ctx context.Context) (result []Brand, err error) {
	q := "SELECT id,name FROM device_producers"

	err = r.db.SelectContext(ctx, &result, q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) GetHashrate(ctx context.Context) (result []HashrateType, err error) {
	q := "SELECT id,name FROM hashrate"

	err = r.db.SelectContext(ctx, &result, q)
	if err != nil {
		return nil, err
	}
	return result, nil

}
func (r *Repository) GetCoins(ctx context.Context) (result []Coin, err error) {
	q := "SELECT id,name FROM coins"

	err = r.db.SelectContext(ctx, &result, q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) GetOffers(ctx context.Context) (result []OfferType, err error) {
	q := "SELECT id,name FROM offer_types"

	err = r.db.SelectContext(ctx, &result, q)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}
