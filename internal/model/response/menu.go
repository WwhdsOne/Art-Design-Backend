package response

import "github.com/dromara/carbon/v2"

type Menu struct {
	ID        int64           `json:"id,string"`
	Name      string          `json:"name"`
	Path      string          `json:"path"`
	Component string          `json:"component"`
	ParentID  int64           `json:"parentID,string"`
	*Meta     `json:"meta"`   // 不用指针则无法将修改映射到本体
	Sort      int             `json:"sort"`
	Children  []Menu          `json:"children"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type Meta struct {
	Title             string     `json:"title"`
	Icon              string     `json:"icon"`
	ShowBadge         bool       `json:"showBadge"`
	ShowTextBadge     string     `json:"showTextBadge"`
	IsHide            bool       `json:"isHide"`
	IsHideTab         bool       `json:"isHideTab"`
	Link              string     `json:"link"`
	IsIframe          bool       `json:"isIframe"`
	KeepAlive         bool       `json:"keepAlive"`
	AuthList          []MenuAuth `json:"authList"`
	IsInMainContainer bool       `json:"isInMainContainer"`
}

type MenuAuth struct {
	ID       int64  `json:"id,string"`
	Name     string `json:"name"`
	AuthCode string `json:"authCode"`
}
