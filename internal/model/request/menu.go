package request

import "Art-Design-Backend/internal/model/base"

type Menu struct {
	ID        *base.LongStringID `json:"id" label:"菜单ID"`
	Type      int8               `json:"type" binding:"required" label:"菜单类型"`
	Name      string             `json:"name" binding:"required,min=2,max=20" label:"组件名称"`
	Path      string             `json:"path" label:"路由地址"`
	Component string             `json:"component" label:"组件路径"`
	ParentID  base.LongStringID  `json:"parentID" binding:"required" label:"上级菜单ID"`
	Meta      `json:"meta"`
	Sort      int `json:"sort" label:"排序"`
}

type Meta struct {
	Title             string `json:"title" binding:"required,min=2,max=20" label:"菜单名称"`
	Icon              string `json:"icon" label:"菜单图标"`
	ShowBadge         bool   `json:"showBadge" label:"显示徽标"`
	ShowTextBadge     string `json:"showTextBadge" label:"文本徽标内容"`
	IsHide            bool   `json:"isHide" label:"菜单隐藏"`
	IsHideTab         bool   `json:"isHideTab" label:"标签页隐藏"`
	Link              string `json:"link" label:"外部链接"`
	IsIframe          bool   `json:"isIframe" label:"iframe嵌入"`
	KeepAlive         bool   `json:"keepAlive" label:"缓存组件"`
	IsInMainContainer bool   `json:"isInMainContainer" label:"主容器内显示"`
}
