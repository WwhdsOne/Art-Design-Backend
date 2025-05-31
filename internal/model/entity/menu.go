package entity

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/pkg/constant/tablename"
)

type Menu struct {
	base.BaseModel
	Name      *string `gorm:"type:varchar(30);unique;comment:组件名称"`
	Type      int8    `gorm:"type:smallint;not null;default:1;comment:类型（1：目录；2：菜单；3：按钮）"`
	Path      *string `gorm:"type:varchar(255);unique;comment:路由地址"`
	Component string  `gorm:"type:varchar(255);comment:组件路径"`
	ParentID  int64   `gorm:"type:bigint;not null;default:0;index:idx_parent_id;comment:上级菜单ID,顶级菜单父ID为-1"`
	Meta              // 页面元信息
	Sort      int     `gorm:"type:integer;default:999;comment:排序"`
	Children  []Menu  `gorm:"-"` // 子页面
}

type Meta struct {
	Title             string `gorm:"type:varchar(100);not null;unique;comment:菜单名称"`
	Icon              string `gorm:"type:varchar(20);comment:菜单图标"`
	ShowBadge         bool   `gorm:"type:boolean;default:false;comment:是否显示徽标"`
	ShowTextBadge     string `gorm:"type:varchar(20);comment:文本徽标内容"`
	IsHide            bool   `gorm:"type:boolean;default:false;comment:是否在菜单中隐藏"`
	IsHideTab         bool   `gorm:"type:boolean;default:false;comment:是否在上方标签页中隐藏"`
	Link              string `gorm:"type:varchar(255);comment:外部链接地址"`
	IsIframe          bool   `gorm:"type:boolean;default:false;comment:是否是iframe嵌入"`
	KeepAlive         bool   `gorm:"type:boolean;default:false;comment:是否缓存组件"`
	AuthCode          string `gorm:"type:varchar(255);comment:权限标识码,只有按钮会填充这列"`
	IsInMainContainer bool   `gorm:"type:boolean;default:false;comment:是否在主容器中"`
}

// TableName 设置表名
func (m *Menu) TableName() string {
	return tablename.MenuTableName
}
