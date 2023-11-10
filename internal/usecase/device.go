package usecase

import (
	"BTCcalc/internal/repo"
	"context"
	"database/sql"
	null "gopkg.in/guregu/null.v3/zero"
)

type Usecase struct {
	db *repo.Repository
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
	ArticleID []sql.NullInt64 `json:"articleID,omitempty"`
}
type CaseImageDTO struct {
	CaseID []sql.NullInt64 `json:"caseID,omitempty"`
}

func (r *Usecase) GetDevices(ctx context.Context, p DeviceDTO) ([]repo.Device, error) {
	some := repo.DeviceDTO{DeviceID: p.DeviceID, PriceLow: p.PriceLow, PriceHigh: p.PriceHigh, PowerLow: p.PowerLow, PowerHigh: p.PowerHigh, HashrateHigh: p.HashrateHigh, HashrateLow: p.HashrateLow, HashrateID: p.HashrateID, BrandID: p.BrandID, OfferID: p.OfferID, CoinID: p.CoinID, Recommended: p.Recommended}

	//Абстрактный sql ,  с которого получаем данные
	result, err := r.db.GetDevices(ctx, some)
	if err != nil {
		return []repo.Device{}, err
	}
	return result, nil
}
func (r *Usecase) GetPowerfulDevices(ctx context.Context) ([]repo.Device, error) {

	//Абстрактный sql ,  с которого получаем данные
	result, err := r.db.GetPowerfulDevices(ctx)
	if err != nil {
		return []repo.Device{}, err
	}
	return result, nil
}

func (r *Usecase) GetDeviceImage(ctx context.Context, p DeviceImageDTO) ([]repo.DeviceImage, error) {
	some := repo.DeviceImageDTO{DeviceID: p.DeviceID}

	//Абстрактный sql ,  с которого получаем данные
	result, err := r.db.GetDeviceImage(ctx, some)
	if err != nil {
		return []repo.DeviceImage{}, err
	}
	return result, nil
}

func (r *Usecase) GetArticleImage(ctx context.Context, p ArticleImageDTO) ([]repo.ArticleImage, error) {
	some := repo.ArticleImageDTO{ArticleID: p.ArticleID}

	//Абстрактный sql ,  с которого получаем данные
	result, err := r.db.GetArticleImage(ctx, some)
	if err != nil {
		return []repo.ArticleImage{}, err
	}
	return result, nil
}

func (r *Usecase) GetCaseImage(ctx context.Context, p CaseImageDTO) ([]repo.CaseImage, error) {
	some := repo.CaseImageDTO{CaseID: p.CaseID}

	//Абстрактный sql ,  с которого получаем данные
	result, err := r.db.GetCaseImages(ctx, some)
	if err != nil {
		return []repo.CaseImage{}, err
	}
	return result, nil
}

func (r *Usecase) GetReviews(ctx context.Context, id null.Int) ([]repo.DeviceReviews, error) {

	//Абстрактный sql ,  с которого получаем данные

	result, err := r.db.GetDeviceReviews(ctx, id)
	if err != nil {
		return []repo.DeviceReviews{}, err
	}
	return result, nil
}
func (r *Usecase) GetCall(ctx context.Context, name string, phone string) ([]repo.DeviceReviews, error) {
	//somecode
	return nil, nil
}

func (r *Usecase) WriteReview(ctx context.Context, email string, phone string, text string, stars int64) (err error) {
	err = r.db.InsertReview(ctx, email, text, phone, stars)
	if err != nil {
		return err
	}
	return nil
}

func (r *Usecase) WriteDeviceReview(ctx context.Context, email string, phone string, text string, deviceID, stars, amount int64) (err error) {
	err = r.db.InsertDeviceReview(ctx, email, text, phone, deviceID, stars, amount)
	if err != nil {
		return err
	}
	return nil
}

func (r *Usecase) GetArticles(ctx context.Context, id null.Int) ([]repo.Articles, error) {

	//Абстрактный sql ,  с которого получаем данные

	result, err := r.db.GetArticles(ctx, id)
	if err != nil {
		return []repo.Articles{}, err
	}
	return result, nil
}

func (r *Usecase) GetCases(ctx context.Context, id null.Int) ([]repo.Cases, error) {

	//Абстрактный sql ,  с которого получаем данные

	result, err := r.db.GetCases(ctx, id)
	if err != nil {
		return []repo.Cases{}, err
	}
	return result, nil
}

func (r *Usecase) GetBrand(ctx context.Context) ([]repo.Brand, error) {

	//Абстрактный sql ,  с которого получаем данные

	result, err := r.db.GetBrands(ctx)
	if err != nil {
		return []repo.Brand{}, err
	}
	return result, nil
}
func (r *Usecase) GetCoin(ctx context.Context) ([]repo.Coin, error) {

	//Абстрактный sql ,  с которого получаем данные

	result, err := r.db.GetCoins(ctx)
	if err != nil {
		return []repo.Coin{}, err
	}
	return result, nil
}
func (r *Usecase) GetHashrate(ctx context.Context) ([]repo.HashrateType, error) {

	//Абстрактный sql ,  с которого получаем данные

	result, err := r.db.GetHashrate(ctx)
	if err != nil {
		return []repo.HashrateType{}, err
	}
	return result, nil
}
func (r *Usecase) GetOffer(ctx context.Context) ([]repo.OfferType, error) {

	//Абстрактный sql ,  с которого получаем данные

	result, err := r.db.GetOffers(ctx)
	if err != nil {
		return []repo.OfferType{}, err
	}
	return result, nil
}

func NewUsecase(db *repo.Repository) *Usecase {
	return &Usecase{db}
}
