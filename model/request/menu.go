package request

type Menu struct {
	Type      int8   `json:"type"`
	Path      string `json:"path"`
	Component string `json:"component"`
	ParentID  int64  `json:"parentID"`
	Meta      `json:"meta"`
	Sort      int `json:"sort"`
}

type Meta struct {
	Title             string `json:"title"`
	Icon              string `json:"icon"`
	ShowBadge         bool   `json:"showBadge"`
	ShowTextBadge     string `json:"showTextBadge"`
	IsHide            bool   `json:"isHide"`
	IsHideTab         bool   `json:"isHideTab"`
	Link              string `json:"link"`
	IsIframe          bool   `json:"isIframe"`
	KeepAlive         bool   `json:"keepAlive"`
	AuthList          string `json:"authList"`
	IsInMainContainer bool   `json:"isInMainContainer"`
}
