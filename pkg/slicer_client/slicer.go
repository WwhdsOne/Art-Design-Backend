package slicer_client

import (
	"bytes"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
)

type Slicer struct {
	SlicerURL string `mapstructure:"slicer_url" yaml:"slicer_url"`
}

type SlicerResponse struct {
	Chunks []string `json:"chunks"`
}

func (s *Slicer) GetChunksFromSlicer(fileURL string) ([]string, error) {
	requestData, _ := sonic.Marshal(map[string]string{
		"url": fileURL,
	})
	resp, err := http.Post(s.SlicerURL, "application/json", bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 响应
	var result SlicerResponse
	if err := sonic.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Chunks, nil
}
