package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/resp"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/loginUtils"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type MenuService struct {
	MenuRepo *repository.MenuRepository // 用户Repo
	Redis    *redisx.RedisWrapper       // redis
}

func (m *MenuService) CreateMenu(c context.Context, menu *request.Menu) (err error) {
	var menuDo entity.Menu
	err = copier.Copy(&menuDo, &menu)
	// 默认父级为顶级菜单
	if menu.ParentID == 0 {
		menuDo.ParentID = -1
	}
	if err != nil {
		zap.L().Error("菜单属性复制失败", zap.Error(err))
		return
	}
	err = m.MenuRepo.CreateMenu(c, &menuDo)
	if err != nil {
		return
	}
	return
}

func (m *MenuService) GetMenuList(c context.Context) (res []*resp.Menu, err error) {
	roleIds := loginUtils.GetUserRoleIDs(c)
	menuList, err := m.MenuRepo.GetMenuListByRoleIDList(c, roleIds)
	// 先用 map 存储所有菜单，方便查找
	menuMap := make(map[int64]*resp.Menu)
	for _, menuDo := range menuList {
		var menuResp resp.Menu
		// 如果是不是按钮类型，则初始化 AuthList 和 Children
		if menuDo.Type != 3 {
			menuResp.Meta.AuthList = make([]string, 0)
			menuResp.Children = make([]resp.Menu, 0)
		}
		err = copier.Copy(&menuResp, &menuDo)
		if err != nil {
			return
		}
		menuMap[menuDo.ID] = &menuResp
	}

	// 遍历所有菜单，构建树形结构
	for _, dbMenu := range menuList {
		frontendMenu := menuMap[dbMenu.ID]

		// 跳过按钮类型（按钮的 AuthCode 会挂到父菜单上）
		if dbMenu.Type == 3 { // 按钮类型
			if parent, exists := menuMap[dbMenu.ParentID]; exists {
				// 将按钮的 AuthCode 添加到父菜单的 AuthList
				parent.Meta.AuthList = append(parent.Meta.AuthList, dbMenu.Meta.AuthCode)
			}
			continue
		}

		// 如果是顶级菜单，直接添加到结果列表
		if dbMenu.ParentID == -1 {
			res = append(res, frontendMenu)
		} else {
			// 否则挂到父菜单的 Children 下
			if parent, exists := menuMap[dbMenu.ParentID]; exists {
				parent.Children = append(parent.Children, *frontendMenu)
			}
		}
	}

	return
}
