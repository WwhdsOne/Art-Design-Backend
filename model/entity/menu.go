package entity

import "Art-Design-Backend/model/base"

type Menu struct {
	base.BaseModel
	Name      string `gorm:"column:name;type:varchar(50);unique;comment:组件名称"`
	Type      int8   `gorm:"column:type;type:tinyint;not null;default:1;comment:类型（1：目录；2：菜单；3：按钮）"`
	Path      string `gorm:"column:path;type:varchar(255);unique;comment:路由地址"`
	Component string `gorm:"column:component;type:varchar(255);unique;comment:组件路径"`
	ParentID  int64  `gorm:"column:parent_id;type:bigint;not null;default:0;index:idx_parent_id;comment:上级菜单ID"`
	Meta             // 页面元信息
	Redirect  string `gorm:"column:redirect;type:varchar(255);comment:重定向地址"`
	Sort      int    `gorm:"column:sort;type:int;not null;default:999;comment:排序"`
	Status    int8   `gorm:"column:status;type:tinyint;not null;default:1;comment:状态（1：启用；2：禁用）"`
	Children  []Menu `gorm:"-"` // 子页面
}

type Meta struct {
	Title             string   `gorm:"column:title;type:varchar(100);not null;unique;comment:菜单名称"`
	Icon              string   `gorm:"column:icon;type:varchar(50);comment:菜单图标"`
	ShowBadge         bool     `gorm:"column:show_badge;type:boolean;default:false;comment:是否显示徽标"`
	ShowTextBadge     string   `gorm:"column:show_text_badge;type:varchar(20);comment:文本徽标内容"`
	IsHide            bool     `gorm:"column:is_hide;type:boolean;default:false;comment:是否在菜单中隐藏"`
	IsHideTab         bool     `gorm:"column:is_hide_tab;type:boolean;default:0;comment:是否在上方标签页中隐藏"`
	Link              string   `gorm:"column:link;type:varchar(255);comment:外部链接地址"`
	IsIframe          bool     `gorm:"column:is_iframe;type:boolean;default:false;comment:是否是iframe嵌入"`
	IsCache           bool     `gorm:"column:is_cache;type:boolean;default:false;comment:是否缓存组件"`
	AuthList          []string `gorm:"-"` // 权限标识码列表,用于鉴权前端某些按钮是否能显示
	IsInMainContainer bool     `gorm:"column:is_in_main_container;type:boolean;default:1;comment:是否在主容器中"`
}

// TableName 设置表名
func (Menu) TableName() string {
	return "menu"
}
