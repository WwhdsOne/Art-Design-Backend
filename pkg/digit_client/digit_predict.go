package digit_client

import (
	"bytes"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
)

type DigitPredict struct {
	PredictUrl string
}

func (c *DigitPredict) Predict(imageUrl string) (result int, err error) {
	request := map[string]string{
		"image_url": imageUrl,
	}
	requestData, err := sonic.Marshal(request)
	if err != nil {
		return
	}

	resp, err := http.Post(c.PredictUrl, "application/json", bytes.NewBuffer(requestData))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var respMap map[string]int
	err = sonic.Unmarshal(responseData, &respMap)
	if err != nil {
		return
	}

	result = respMap["predicted_class"]
	return
}
