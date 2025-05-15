package request

import "Art-Design-Backend/internal/model/base"

type Menu struct {
	ID        base.LongStringID `json:"id"`
	Type      int8              `json:"type" binding:"required"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Component string            `json:"component"`
	ParentID  base.LongStringID `json:"parentID"`
	Meta      `json:"meta"`
	Sort      int `json:"sort"`
}

type Meta struct {
	Title             string `json:"title" binding:"required"`
	Icon              string `json:"icon"`
	ShowBadge         bool   `json:"showBadge"`
	ShowTextBadge     string `json:"showTextBadge"`
	IsHide            bool   `json:"isHide"`
	IsHideTab         bool   `json:"isHideTab"`
	Link              string `json:"link"`
	IsIframe          bool   `json:"isIframe"`
	KeepAlive         bool   `json:"keepAlive"`
	AuthCode          string `json:"authCode"`
	IsInMainContainer bool   `json:"isInMainContainer"`
}
