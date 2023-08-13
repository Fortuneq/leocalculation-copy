package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

type Devices struct {
	Name      string `db:"name"`
	Id        string `db:"id"`
	Cost      string `db:"cost"`
	Image     string `db:"image"`
	Size      string `db:"size"`
	Power     string `db:"power"`
	Hashrate  string `db:"hashrate"`
	UID       string `db:"uid"`
	Algorithm string `db:"algorithm"`
	VideoUrl  string `db:"video_url"`
}

func (r *Repository) GetDevices(ctx context.Context, id int) (Devices, error) {

	//Абстрактный sql ,  с которого получаем данные

	q := "SELECT id,name, cost,image,size,power,hashrate,algorithm,uid,video_url FROM devices  where id = ?"

	place := Devices{}

	err := r.db.GetContext(ctx, &place, q, id)
	if err != nil {
		return Devices{}, err
	}
	return place, nil
}
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}
