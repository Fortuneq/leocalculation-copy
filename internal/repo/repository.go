package repo

import (
	"context"
	"database/sql"
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
	DeviceID     sql.NullInt64   `json:"deviceID,omitempty"`
	PriceLow     sql.NullFloat64 `json:"priceLow,omitempty"`
	PriceHigh    sql.NullFloat64 `json:"priceHigh,omitempty"`
	PowerLow     sql.NullFloat64 `json:"powerLow,omitempty"`
	PowerHigh    sql.NullFloat64 `json:"powerHigh,omitempty"`
	HashrateLow  sql.NullFloat64 `json:"hashrateLow,omitempty"`
	HashrateHigh sql.NullFloat64 `json:"hashrateHigh,omitempty"`
	HashrateID   sql.NullInt64   `json:"hashrateID,omitempty"`
	BrandID      sql.NullInt64   `json:"brandID,omitempty"`
	OfferID      sql.NullInt64   `json:"offerID,omitempty"`
	CoinID       sql.NullInt64   `json:"coinID,omitempty"`
	Recommended  sql.NullInt64   `json:"recommended,omitempty"`
}

type DeviceImageDTO struct {
	DeviceID []int64 `json:"deviceID,omitempty"`
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
	q := ""
	//Абстрактный sql ,  с которого получаем данные
	q = "SELECT DISTINCT devices.id,devices.name as name, cost,size,power,hashrate,algorithm,uid,video_url,c.name as coin_name,     " +
		"           h.name as hash_name,ot.name as offer_name,recommended,dp.name as brand_name FROM devices    JOIN device_coin dc on devices.id = dc.device_id  " +
		"  join coins c on dc.coin_id = c.id join device_producers dp on dp.id = devices.producer_id  " +
		"  join hashrate h on h.id = devices.hashrate_id   " +
		" join offer_types ot on devices.offer_type = ot.id WHERE (? is null or  cost >= ?)  and (? is null or cost <= ?) and( ? is null or power >= ?) and( ? is null or power <=?)" +
		"and (? is null or  hashrate >= ? )and (? is null or hashrate <= ? )and ( ? is null or ot.id in(?)) and( ? is null or h.id =? )and (? is null or c.id in(?)) and( ? is null or dp.id in(?))" +
		"  and( ? is null or recommended = ?) and( ? is null or devices.id = ?)"

	err = r.db.SelectContext(ctx, &result, q, p.PriceLow, p.PriceLow, p.PriceHigh, p.PriceHigh, p.PowerLow, p.PowerLow, p.PowerHigh, p.PowerHigh, p.HashrateLow, p.HashrateLow, p.HashrateHigh, p.HashrateHigh,
		p.OfferID, p.OfferID, p.HashrateID, p.HashrateID, p.CoinID, p.CoinID, p.BrandID, p.BrandID, p.Recommended, p.Recommended, p.DeviceID, p.DeviceID)
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
	if len(p.DeviceID) == 0 {
		stmt := `SELECT id,image from devices`
		err = r.db.SelectContext(ctx, &result, stmt)
		if err != nil {
			return []DeviceImage{}, err
		}
		return result, err
	}
	stmt := `SELECT id,image from devices where id in(?` + strings.Repeat(",?", len(args)-1) + `)`
	err = r.db.SelectContext(ctx, &result, stmt, args...)
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
