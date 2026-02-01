package digit_client

import (
	"bytes"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
)

type DigitPredict struct {
	PredictURL string
}

func (c *DigitPredict) Predict(imageURL string) (result int, err error) {
	request := map[string]string{
		"image_url": imageURL,
	}
	requestData, err := sonic.Marshal(request)
	if err != nil {
		return
	}

	resp, err := http.Post(c.PredictURL, "application/json", bytes.NewBuffer(requestData))
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
