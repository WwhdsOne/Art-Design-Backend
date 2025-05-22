package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/aliyun"
	"Art-Design-Backend/pkg/client"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"mime/multipart"
)

type DigitPredictService struct {
	DigitPredictRepo   *repository.DigitPredictRepository
	DigitPredictClient *client.DigitPredict
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

func (d *DigitPredictService) SubmitMission(c context.Context, req *request.DigitPredict) error {
	go func() {
		var result int
		var err error
		if result, err = d.DigitPredictClient.Predict(req.Image); err != nil {
			zap.L().Error("预测任务失败", zap.Error(err))
		}
		zap.L().Info("预测结果", zap.Int("result", result))
		err = d.DigitPredictRepo.UpdateLabelByID(c, int64(req.ID), result)
		if err != nil {
			zap.L().Error("更新任务结果失败", zap.Error(err))
		}
	}()
	return nil
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
