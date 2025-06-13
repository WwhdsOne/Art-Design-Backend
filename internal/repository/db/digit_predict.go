package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DigitPredictRepository struct {
	db *gorm.DB
}

func NewDigitPredictRepository(db *gorm.DB) *DigitPredictRepository {
	return &DigitPredictRepository{
		db: db,
	}
}

func (r *DigitPredictRepository) GetDigitPredictList(c context.Context, createdBy int64) (digitPredictList []*entity.DigitPredict, err error) {
	if err = DB(c, r.db).
		Where("created_by = ?", createdBy).
		Find(&digitPredictList).Error; err != nil {
		zap.L().Error("查询用户预测列表失败", zap.Error(err))
		err = errors.NewDBError("查询用户预测列表失败")
		return
	}
	return
}

func (r *DigitPredictRepository) UpdateLabelByID(c context.Context, id int64, label int) (err error) {
	if err = DB(c, r.db).Model(&entity.DigitPredict{}).
		Where("id = ?", id).
		Update("label", label).Error; err != nil {
		zap.L().Error("更新数字预测结果失败", zap.Error(err))
		err = errors.NewDBError("更新数字预测结果失败")
		return
	}
	return
}

func (r *DigitPredictRepository) Create(c context.Context, predict *entity.DigitPredict) (err error) {
	if err = DB(c, r.db).Create(predict).Error; err != nil {
		err = errors.NewDBError("创建数字识别任务失败")
		return
	}
	return
}

func (r *DigitPredictRepository) IsLabeled(c context.Context, id int64) (bool, error) {
	var count int64
	if err := DB(c, r.db).Model(&entity.DigitPredict{}).
		Where("id = ?", id).
		Where("label != null").
		Count(&count).Error; err != nil {
		zap.L().Error("查询数字识别结果失败", zap.Error(err))
		return false, errors.NewDBError("查询数字识别结果失败")
	}
	return count > 0, nil
}

func (r *DigitPredictRepository) GetDigitById(c *gin.Context, id int64) (res *entity.DigitPredict, err error) {
	if err = DB(c, r.db).Where("id = ?", id).First(&res).Error; err != nil {
		zap.L().Error("查询数字识别结果失败", zap.Error(err))
		err = errors.NewDBError("查询数字识别结果失败")
		return
	}
	return
}
