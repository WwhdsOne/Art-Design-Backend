package base

import (
	"Art-Design-Backend/pkg/utils"
	"context"
	"github.com/dromara/carbon/v2"
	"gorm.io/gorm"
)

// BaseModel 是一个通用的模型，包含 ID、CreatedAt、UpdatedAt 和 DeletedAt 字段
type BaseModel struct {
	ID        ID              // 雪花ID
	CreatedAt carbon.DateTime `gorm:"autoCreateTime" json:"createdAt"`   // 创建记录时自动设置为当前时间
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime" json:"updatedAt"`   // 更新记录时自动设置为当前时间
	UpdateBy  int64           `gorm:"column:updated_by" json:"updateBy"` // 修改人字段，记录最后一次更新操作者的标识
	CreateBy  int64           `gorm:"column:created_by" json:"createBy"` // 创建人字段，记录创建操作者的标识
}

// BeforeCreate 在创建记录之前自动生成雪花 ID
func (b *BaseModel) BeforeCreate(db *gorm.DB) (err error) {
	b.fillAddReq(db.Statement.Context)
	return
}

func (b *BaseModel) BeforeUpdate(db *gorm.DB) (err error) {
	b.fillUpdateReq(db.Statement.Context)
	return
}

func (b *BaseModel) fillAddReq(c context.Context) {
	userID := utils.GetUserID(c)
	// 如果存在 claims，正常提取用户 ID
	b.CreateBy = userID
	b.UpdateBy = userID
}

func (b *BaseModel) fillUpdateReq(c context.Context) {
	userID := utils.GetUserID(c)
	b.UpdateBy = userID
}
