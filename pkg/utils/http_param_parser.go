package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ParseIDs 解析逗号分隔的 ID 字符串，返回一个 int64 切片
func ParseIDs(c *gin.Context) (ids []int64, err error) {
	idsParam, _ := c.Params.Get("ids")
	if idsParam == "" {
		return nil, errors.New("IDs 为空")
	}
	idStrings := strings.SplitSeq(idsParam, ",")
	for idStr := range idStrings {
		var id int64
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			err = errors.New("ID 解析错误")
			return
		}
		ids = append(ids, id)
	}
	return
}

func ParseID(c *gin.Context) (id int64, err error) {
	idParam, exist := c.Params.Get("id")
	if !exist {
		err = errors.New("ID 不存在")
		return
	}
	id, err = strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		err = errors.New("ID 解析错误")
		return
	}
	return
}
