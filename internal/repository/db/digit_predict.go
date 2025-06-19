package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type DigitPredictDB struct {
	db *gorm.DB
}

func NewDigitPredictDB(db *gorm.DB) *DigitPredictDB {
	return &DigitPredictDB{
		db: db,
	}
}

func (r *DigitPredictDB) GetDigitPredictList(c context.Context, createdBy int64) (digitPredictList []*entity.DigitPredict, err error) {
	if err = DB(c, r.db).
		Where("created_by = ?", createdBy).
		Find(&digitPredictList).Error; err != nil {
		err = errors.WrapDBError(err, "查询用户预测列表失败")
		return
	}
	return
}

func (r *DigitPredictDB) UpdateLabelByID(c context.Context, id int64, label int) (err error) {
	if err = DB(c, r.db).Model(&entity.DigitPredict{}).
		Where("id = ?", id).
		Update("label", label).Error; err != nil {
		return errors.WrapDBError(err, "更新数字预测结果失败")
	}
	return
}

func (r *DigitPredictDB) Create(c context.Context, predict *entity.DigitPredict) (err error) {
	if err = DB(c, r.db).Create(predict).Error; err != nil {
		return errors.WrapDBError(err, "创建数字识别任务失败")
	}
	return
}

func (r *DigitPredictDB) IsLabeled(c context.Context, id int64) (bool, error) {
	var count int64
	if err := DB(c, r.db).Model(&entity.DigitPredict{}).
		Where("id = ?", id).
		Where("label != null").
		Count(&count).Error; err != nil {
		return false, errors.WrapDBError(err, "查询数字识别结果失败")
	}
	return count > 0, nil
}

func (r *DigitPredictDB) GetDigitById(c context.Context, id int64) (res *entity.DigitPredict, err error) {
	if err = DB(c, r.db).Where("id = ?", id).First(&res).Error; err != nil {
		err = errors.WrapDBError(err, "查询数字识别结果失败")
		return
	}
	return
}
