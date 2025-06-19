package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/db"
	"context"
)

type DigitPredictRepo struct {
	DigitPredictDB *db.DigitPredictDB
}

func NewDigitPredictRepo(digitPredictDB *db.DigitPredictDB) *DigitPredictRepo {
	return &DigitPredictRepo{
		DigitPredictDB: digitPredictDB,
	}
}

func (d *DigitPredictRepo) GetDigitPredictList(c context.Context, createdBy int64) (digitPredictList []*entity.DigitPredict, err error) {
	return d.DigitPredictDB.GetDigitPredictList(c, createdBy)
}

func (d *DigitPredictRepo) IsLabeled(c context.Context, id int64) (labeled bool, err error) {
	return d.DigitPredictDB.IsLabeled(c, id)
}

func (d *DigitPredictRepo) UpdateLabelByID(c context.Context, id int64, label int) (err error) {
	return d.DigitPredictDB.UpdateLabelByID(c, id, label)
}

func (d *DigitPredictRepo) Create(c context.Context, e *entity.DigitPredict) (err error) {
	return d.DigitPredictDB.Create(c, e)
}

func (d *DigitPredictRepo) GetDigitById(c context.Context, id int64) (res *entity.DigitPredict, err error) {
	return d.DigitPredictDB.GetDigitById(c, id)
}
