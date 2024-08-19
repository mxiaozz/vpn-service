package system

import (
	"strings"

	"github.com/pkg/errors"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/model/response"
	"vpn-web.funcworks.net/util"
)

var MenuService = &menuService{}

type menuService struct {
}

// 用户登录获取左侧菜单树结构
func (ms *menuService) GetUserMenuTree(user *model.LoginUser) ([]response.RouterVo, error) {
	var menus []entity.SysMenu
	var err error
	if user.User.IsAdmin() {
		menus, err = ms.selectMenuAll(false)
	} else {
		menus, err = ms.selectMenuByUserId(user.UserId, false)
	}
	if err != nil {
		return nil, err
	}
	menus = ms.buildMenuTree(menus, 0)
	return ms.buildRouters(menus), nil
}

// 读取所有菜单(M & C)，isWithFuncAction true时包含 F
func (ms *menuService) selectMenuAll(isWithFuncAction bool) ([]entity.SysMenu, error) {
	dbSession := gb.DB.Table("sys_menu").Alias("m").Where("m.status = '0'").
		Asc("m.parent_id", "m.order_num")
	if !isWithFuncAction {
		dbSession.And("m.menu_type in ('M', 'C')")
	}

	var menus = []entity.SysMenu{}
	err := dbSession.Find(&menus)
	return menus, err
}

// 读取用户(非admin)关联的所有菜单(M & C)，isWithFuncAction true时包含 F
func (ms *menuService) selectMenuByUserId(userId int64, isWithFuncAction bool) ([]entity.SysMenu, error) {
	dbSession := gb.DB.Table("sys_menu").Alias("m").
		Select("distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time").
		Join("left", []string{"sys_role_menu", "rm"}, "m.menu_id = rm.menu_id").
		Join("left", []string{"sys_user_role", "ur"}, "rm.role_id = ur.role_id").
		Join("left", []string{"sys_role", "ro"}, "ur.role_id = ro.role_id").
		Where("ur.user_id = ? and m.status = 0 and ro.status = 0", userId).
		Asc("m.parent_id", "m.order_num")
	if !isWithFuncAction {
		dbSession.And("m.menu_type in ('M', 'C')")
	}

	var menus = []entity.SysMenu{}
	err := dbSession.Find(&menus)
	return menus, err
}

// 构造菜单树结构
func (ms *menuService) buildMenuTree(menus []entity.SysMenu, parentId int64) []entity.SysMenu {
	list := make([]entity.SysMenu, 0)
	for _, menu := range menus {
		if menu.ParentId == parentId {
			c := ms.buildMenuTree(menus, menu.MenuId)
			menu.Children = c
			// tree select
			menu.Id = menu.MenuId
			menu.Label = menu.MenuName
			list = append(list, menu)
		}
	}
	return list
}

// 菜单树转换为路由结构
func (ms *menuService) buildRouters(menus []entity.SysMenu) []response.RouterVo {
	routers := make([]response.RouterVo, 0)
	for _, menu := range menus {
		router := response.RouterVo{}
		router.Hidden = (menu.Visible == "1")
		router.Name = ms.getRouteName(&menu)
		router.Path = ms.getRouterPath(&menu)
		router.Component = ms.getComponent(&menu)
		router.Query = menu.Query
		router.Meta = &response.MetaVo{
			Title:   menu.MenuName,
			Icon:    menu.Icon,
			NoCache: menu.IsCache == "1",
		}
		if util.IsHttp(menu.Path) {
			router.Meta.Link = menu.Path
		}

		children := menu.Children
		// 目录
		if len(children) > 0 && cst.MENU_TYPE_DIR == menu.MenuType {
			router.AlwaysShow = true
			router.Redirect = "noRedirect"
			router.Children = ms.buildRouters(children)
		} else
		// 一级菜单
		if ms.isMenuFrame(&menu) {
			c := response.RouterVo{
				Path:      menu.Path,
				Component: menu.Component,
				Name:      menu.Path,
				Query:     menu.Query,
				Meta: &response.MetaVo{
					Title:   menu.MenuName,
					Icon:    menu.Icon,
					NoCache: menu.IsCache == "1",
					Link:    menu.Path,
				},
			}
			router.Meta = nil
			router.Children = []response.RouterVo{c}
		} else
		// 一级菜单链接或一级目录链接
		if menu.ParentId == 0 && ms.isInnerLink(&menu) {
			path := ms.innerLinkReplaceEach(menu.Path)
			c := response.RouterVo{
				Path:      path,
				Component: cst.MENU_INNER_LINK,
				Name:      path,
				Meta: &response.MetaVo{
					Title: menu.MenuName,
					Icon:  menu.Icon,
					Link:  menu.Path,
				},
			}

			router.Meta = &response.MetaVo{
				Title: menu.MenuName,
				Icon:  menu.Icon,
			}
			router.Path = "/"
			router.Children = []response.RouterVo{c}
		}
		routers = append(routers, router)
	}
	return routers
}

func (ms *menuService) getRouteName(menu *entity.SysMenu) string {
	if ms.isMenuFrame(menu) {
		return ""
	}
	return menu.Path
}

// 是否为菜单内部跳转
func (ms *menuService) isMenuFrame(menu *entity.SysMenu) bool {
	return menu.ParentId == 0 &&
		cst.MENU_TYPE_MENU == menu.MenuType &&
		menu.IsFrame == cst.MENU_NO_FRAME
}

// 是否为内链组件
func (ms *menuService) isInnerLink(menu *entity.SysMenu) bool {
	return menu.IsFrame == cst.MENU_NO_FRAME &&
		util.IsHttp(menu.Path)
}

func (ms *menuService) innerLinkReplaceEach(path string) string {
	oldList := []string{cst.SYS_HTTP, cst.SYS_HTTPS, cst.SYS_WWW, ".", ":"}
	newList := []string{"", "", "", "/", "/"}
	for i, o := range oldList {
		path = strings.ReplaceAll(path, o, newList[i])
	}
	return path
}

func (ms *menuService) getRouterPath(menu *entity.SysMenu) string {
	routerPath := menu.Path
	// 内链打开外网方式
	if menu.ParentId != 0 && ms.isInnerLink(menu) {
		routerPath = ms.innerLinkReplaceEach(routerPath)
	}
	// 非外链并且是一级目录（类型为目录）
	if menu.ParentId == 0 && cst.MENU_TYPE_DIR == menu.MenuType &&
		cst.MENU_NO_FRAME == menu.IsFrame {
		routerPath = "/" + menu.Path
	} else
	// 非外链并且是一级目录（类型为菜单）
	if ms.isMenuFrame(menu) {
		routerPath = "/"
	}
	return routerPath
}

func (ms *menuService) getComponent(menu *entity.SysMenu) string {
	// 顶级目录/菜单
	component := cst.MENU_LAYOUT

	// 只要 Component 不为空都认为是菜单
	// 子菜单（目录的 Component 认为是空的）
	if menu.Component != "" && !ms.isMenuFrame(menu) {
		component = menu.Component
	} else
	// 子菜单或子目录（只要设置 Component 为空），并且是内部链接
	if menu.Component == "" && menu.ParentId != 0 && ms.isInnerLink(menu) {
		component = cst.MENU_INNER_LINK
	} else
	// 子目录（目录的 Component 认为是空的）
	if menu.Component == "" && ms.isParentView(menu) {
		component = cst.MENU_PARENT_VIEW
	}

	return component
}

func (ms *menuService) isParentView(menu *entity.SysMenu) bool {
	return menu.ParentId != 0 && cst.MENU_TYPE_DIR == menu.MenuType
}

// 获取用户菜单权限标识列表，admin为 *：*：*
func (ms *menuService) GetMenuPermission(user *entity.SysUser) (map[string]int8, error) {
	permSet := make(map[string]int8)
	if user.IsAdmin() {
		permSet[cst.SYS_ALL_PERMISSION] = 1
		return permSet, nil
	}

	permList, err := ms.selectMenuPermsByUserId(user.UserId)
	if err != nil {
		return nil, err
	}

	for _, perm := range permList {
		if perm == "" {
			continue
		}
		keys := strings.Split(perm, ",")
		for _, k := range keys {
			permSet[k] = 1
		}
	}
	return permSet, nil
}

// 获取用户所关联的菜单(M/C/F)权限perms标识(未解析)
func (ms *menuService) selectMenuPermsByUserId(userId int64) ([]string, error) {
	var perms = []string{}
	err := gb.DB.Table("sys_menu").Alias("m").
		Select("distinct ifnull(m.perms, '')").
		Join("left", []string{"sys_role_menu", "rm"}, "m.menu_id = rm.menu_id").
		Join("left", []string{"sys_user_role", "ur"}, "rm.role_id = ur.role_id").
		Join("left", []string{"sys_role", "r"}, "r.role_id = ur.role_id").
		Where("m.status = '0' and r.status = '0' and ur.user_id = ?", userId).Find(&perms)
	if err != nil {
		return nil, err
	}
	return perms, nil
}

// 菜单管理列表，menu.Params["userId"] 指定具体用户菜单列表，未指定则返回全部菜单(M & C & F)
func (ms *menuService) GetMenuList(menu *entity.SysMenu) ([]entity.SysMenu, error) {
	dbSession := gb.DB.Table("sys_menu").Alias("m").
		Select("m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time")
	if userId, ok := menu.Params["userId"]; ok {
		dbSession = dbSession.Join("inner", []string{"sys_role_menu", "rm"}, "m.menu_id = rm.menu_id")
		dbSession = dbSession.Join("inner", []string{"sys_user_role", "ur"}, "rm.role_id = ur.role_id")
		dbSession = dbSession.Join("inner", []string{"sys_role", "r"}, "r.role_id = ur.role_id")
		dbSession = dbSession.And("ur.user_id = ? and r.status = '0'", userId)
	}
	if menu.MenuName != "" {
		dbSession = dbSession.And("m.menu_name like ?", "%"+menu.MenuName+"%")
	}
	if menu.Visible != "" {
		dbSession = dbSession.And("m.visible = ?", menu.Visible)
	}
	if menu.Status != "" {
		dbSession = dbSession.And("m.status = ?", menu.Status)
	}

	var menus = []entity.SysMenu{}
	err := dbSession.OrderBy("m.parent_id, m.order_num").Find(&menus)
	return menus, err
}

func (ms *menuService) GetMenu(menuId int64) (*entity.SysMenu, error) {
	var menu entity.SysMenu
	if exist, err := gb.DB.Table("sys_menu").Where("menu_id = ?", menuId).Get(&menu); err != nil || !exist {
		return nil, err
	}
	return &menu, nil
}

func (ms *menuService) AddMenu(menu *entity.SysMenu) error {
	if exist, err := ms.checkMenuNameUnique(menu.MenuName, menu.MenuId, menu.ParentId); err != nil {
		return err
	} else if exist {
		return errors.New("已存在相同菜单名称")
	}

	_, err := gb.DB.InsertOne(menu)
	return err
}

func (ms *menuService) UpdateMenu(menu *entity.SysMenu) error {
	if exist, err := ms.checkMenuNameUnique(menu.MenuName, menu.MenuId, menu.ParentId); err != nil {
		return err
	} else if exist {
		return errors.New("已存在相同菜单名称")
	}

	_, err := gb.DB.Where("menu_id = ?", menu.MenuId).Update(menu)
	return err
}

func (ms *menuService) checkMenuNameUnique(menuName string, menuId, parentId int64) (bool, error) {
	return gb.DB.Table("sys_menu").Where("menu_name = ? and parent_id = ? and menu_id <> ?", menuName, parentId, menuId).Exist()
}

func (ms *menuService) DeleteMenu(menuId int64) error {
	if exist, err := ms.existChild(menuId); err != nil {
		return errors.Wrap(err, "删除菜单失败")
	} else if exist {
		return errors.New("存在子菜单,不允许删除")
	}

	if exist, err := ms.existMenuRole(menuId); err != nil {
		return errors.Wrap(err, "删除菜单失败")
	} else if exist {
		return errors.New("菜单已分配,不允许删除")
	}

	_, err := gb.DB.Table("sys_menu").Where("menu_id = ?", menuId).Delete()
	return err
}

func (ms *menuService) existChild(menuId int64) (bool, error) {
	return gb.DB.Table("sys_menu").Where("parent_id = ?", menuId).Exist()
}

func (ms *menuService) existMenuRole(menuId int64) (bool, error) {
	return gb.DB.Table("sys_role_menu").Where("menu_id = ?", menuId).Exist()
}

// 角色管理，角色编辑菜单树选择框列表
func (ms *menuService) GetRolerMenuTreeSelect(user *model.LoginUser) ([]entity.SysMenu, error) {
	var menus []entity.SysMenu
	var err error
	if user.User.IsAdmin() {
		menus, err = ms.selectMenuAll(true)
	} else {
		menus, err = ms.selectMenuByUserId(user.UserId, true)
	}
	if err != nil {
		return nil, err
	}
	return ms.buildMenuTree(menus, 0), nil
}

// 角色关联的所有菜单列表
func (ms *menuService) GetMenuListByRoleId(roleId int64) ([]int64, error) {
	role, err := RoleService.GetRole(roleId)
	if err != nil {
		return nil, err
	}

	dbSession := gb.DB.Table("sys_menu").Alias("m").
		Select("m.menu_id").
		Join("left", []string{"sys_role_menu", "rm"}, "m.menu_id = rm.menu_id").
		Where("m.status = '0' and rm.role_id = ?", roleId)
	if role.MenuCheckStrictly {
		dbSession.And(`m.menu_id not in (
			select m.parent_id from sys_menu m 
			inner join sys_role_menu rm on m.menu_id = rm.menu_id and rm.role_id = ?
		)`, roleId)
	}
	dbSession.OrderBy("m.parent_id, m.order_num")

	var menuIds = []int64{}
	err = dbSession.Find(&menuIds)
	return menuIds, err
}
