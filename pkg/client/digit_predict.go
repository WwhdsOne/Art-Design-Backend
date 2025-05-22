package client

import (
	"bytes"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
)

type DigitPredict struct {
	PredictUrl string
}

type digitPredictResult struct {
	PredictedClass int `json:"predicted_class"`
}

func (c *DigitPredict) Predict(imageUrl string) (result int, err error) {
	request := map[string]string{
		"image_url": imageUrl,
	}
	requestData, err := sonic.Marshal(request)
	if err != nil {
		// 处理错误
		return
	}
	// 使用 bytes.NewBuffer 创建一个 io.Reader
	reader := bytes.NewBuffer(requestData)
	// 发送 POST 请求
	resp, err := http.Post(c.PredictUrl, "application/json", reader)
	if err != nil {
		// 处理错误
		return
	}
	defer resp.Body.Close()
	// 读取响应体
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		// 处理错误
		return
	}
	var r digitPredictResult
	// 解析 JSON 响应体
	err = sonic.Unmarshal(responseData, &r)
	if err != nil {
		// 处理错误
		return
	}
	result = r.PredictedClass
	return
}
