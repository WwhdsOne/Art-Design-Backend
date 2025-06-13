package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/aliyun"
	"Art-Design-Backend/pkg/digit_client"
	"Art-Design-Backend/pkg/errors"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"mime/multipart"
)

type DigitPredictService struct {
	DigitPredictRepo   *db.DigitPredictRepository
	DigitPredictClient *digit_client.DigitPredict
	OssClient          *aliyun.OssClient // 阿里云OSS
}

func (d *DigitPredictService) GetDigitPredictList(c context.Context, createdBy int64) (digitPredictList []*response.DigitPredict, err error) {
	list, err := d.DigitPredictRepo.GetDigitPredictList(c, createdBy)
	if err != nil {
		return
	}
	for _, predict := range list {
		var item response.DigitPredict
		_ = copier.Copy(&item, predict)
		digitPredictList = append(digitPredictList, &item)
	}
	return
}

func (d *DigitPredictService) SubmitMission(c context.Context, req *request.DigitPredict) (err error) {
	id := int64(req.ID)
	labeled, err := d.DigitPredictRepo.IsLabeled(c, id)
	if err != nil {
		zap.L().Error("查询任务状态失败", zap.Error(err))
		return
	}
	if labeled {
		zap.L().Error("该图片已经识别，请勿重复提交", zap.Int64("id", id))
		return errors.NewBusinessError("该图片已经识别，请勿重复提交")
	}
	go func() {
		var result int
		if result, err = d.DigitPredictClient.Predict(req.Image); err != nil {
			zap.L().Error("预测任务失败", zap.Error(err))
			return
		}
		zap.L().Info("预测结果", zap.Int("result", result))
		err = d.DigitPredictRepo.UpdateLabelByID(c, id, result)
		if err != nil {
			zap.L().Error("更新任务结果失败", zap.Error(err))
		}
	}()
	return
}

func (d *DigitPredictService) UploadDigitImage(c *gin.Context, filename string, src multipart.File) (fileUrl string, err error) {
	url, err := d.OssClient.UploadDigitImage(c, filename, src)
	if err != nil {
		zap.L().Error("上传数字图片失败", zap.Error(err))
		return
	}
	var digitDo entity.DigitPredict
	digitDo.Image = url
	err = d.DigitPredictRepo.Create(c, &digitDo)
	if err != nil {
		zap.L().Error("创建数字识别失败", zap.Error(err))
		return
	}
	return
}

func (d *DigitPredictService) GetDigitById(c *gin.Context, id int64) (res *response.DigitPredict, err error) {
	res = &response.DigitPredict{}
	digitDo, err := d.DigitPredictRepo.GetDigitById(c, id)
	if err != nil {
		zap.L().Error("查询数字识别失败", zap.Error(err))
		return
	}
	_ = copier.Copy(res, digitDo)
	return
}
