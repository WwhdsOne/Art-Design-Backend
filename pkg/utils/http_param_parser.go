package utils

import (
	"Art-Design-Backend/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// ParseIDs 解析逗号分隔的 ID 字符串，返回一个 int64 切片
func ParseIDs(c *gin.Context) ([]int64, error) {
	idsParam, _ := c.Params.Get("ids")
	if idsParam == "" {
		return nil, errors.New("IDs 为空")
	}
	idStrings := strings.Split(idsParam, ",")
	var ids []int64
	for _, idStr := range idStrings {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, errors.New("ID 解析错误")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func ParseID(c *gin.Context) (int64, error) {
	idParam, exist := c.Params.Get("id")
	if exist != true {
		return 0, errors.New("ID 不存在")
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.FailWithMessage("ID 解析错误", c)
	}
	return id, nil
}
