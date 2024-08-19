package system

import (
	"strings"

	"github.com/pkg/errors"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/util"
	"xorm.io/builder"
	"xorm.io/xorm"
)

var RoleService = &roleService{}

type roleService struct {
}

// 获取所有角色(不包含已删除)
func (rs *roleService) GetAllRoles() ([]entity.SysRole, error) {
	var roles = []entity.SysRole{}
	err := gb.DB.Table("sys_role").Where("del_flag = '0'").Find(&roles)
	return roles, err
}

// 获取用户角色权限标识
func (rs *roleService) GetUserRolePerms(user *entity.SysUser) (map[string]int8, error) {
	roleSet := make(map[string]int8)
	if user.IsAdmin() {
		roleSet["admin"] = 1
		return roleSet, nil
	}

	roles, err := rs.GetUserRoles(user.UserId)
	if err != nil {
		return nil, err
	}

	util.NewList(roles).
		Filter(func(r entity.SysRole) bool { return r.RoleKey != "" }).
		ForEach(func(r entity.SysRole) {
			keys := strings.Split(r.RoleKey, ",")
			for _, k := range keys {
				roleSet[k] = 1
			}
		})

	return roleSet, nil
}

// 获取用户的角色列表
func (rs *roleService) GetUserRoles(userId int64) ([]entity.SysRole, error) {
	var roles = []entity.SysRole{}
	err := gb.DB.Table("sys_role").Alias("r").Select("*").
		Join("left", []string{"sys_user_role", "ur"}, "ur.role_id = r.role_id").
		Where("ur.user_id = ? and r.del_flag = '0'", userId).Find(&roles)
	return roles, err
}

// 角色管理，获取角色列表（包括角色查询/分页）
func (rs *roleService) GetRoleListPage(role *entity.SysRole, page *model.Page[entity.SysRole]) error {
	return gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_role").
			Where(builder.If(role.RoleName != "", builder.Like{"role_name", role.RoleName}).
				And(builder.If(role.Status != "", builder.Eq{"status": role.Status})).
				And(builder.If(role.RoleKey != "", builder.Like{"role_key", role.RoleKey})).
				And(builder.If(func() bool { return role.Params["beginTime"] != nil }(),
					builder.Gte{"create_time": role.Params["beginTime"]})).
				And(builder.If(func() bool { return role.Params["endTime"] != nil }(),
					builder.Lte{"create_time": role.Params["endTime"]})))
		return builder.Expr("role_sort asc")
	})
}

// 获取角色信息（角色不存在时返回 nil）
func (rs *roleService) GetRole(roleId int64) (*entity.SysRole, error) {
	var role entity.SysRole
	if exist, err := gb.DB.Table("sys_role").Where("role_id = ?", roleId).Get(&role); err != nil || !exist {
		return nil, err
	}
	return &role, nil
}

// 添加角色
func (rs *roleService) AddRole(role *entity.SysRole) error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		if exist, err := rs.checkRoleNameUnique(role.RoleName, role.RoleId); err != nil {
			return err
		} else if exist {
			return errors.New("角色名称已存在")
		}
		if exist, err := rs.checkRoleKeyUnique(role.RoleKey, role.RoleId); err != nil {
			return err
		} else if exist {
			return errors.New("角色权限标识已存在")
		}

		// insert role
		if _, err := dbSession.Insert(role); err != nil {
			return err
		}

		// insert role-menus
		if len(role.MenuIds) > 0 {
			menuList := util.Convert(role.MenuIds, func(id int64) entity.SysRoleMenu {
				return entity.SysRoleMenu{RoleId: role.RoleId, MenuId: id}
			})
			if _, err := dbSession.Insert(menuList); err != nil {
				return errors.Wrap(err, "增加角色时增加角色菜单关联失败")
			}
		}
		return nil
	})
}

// 编辑更新角色
func (rs *roleService) UpdateRole(role *entity.SysRole) error {
	if role.IsAdmin() {
		return errors.New("不允许操作超级管理员角色")
	}

	return gb.Tx(func(dbSession *xorm.Session) error {
		if exist, err := rs.checkRoleNameUnique(role.RoleName, role.RoleId); err != nil {
			return err
		} else if exist {
			return errors.New("角色名称已存在")
		}
		if exist, err := rs.checkRoleKeyUnique(role.RoleKey, role.RoleId); err != nil {
			return err
		} else if exist {
			return errors.New("角色权限已存在")
		}

		// update role
		if _, err := dbSession.Where("role_id = ?", role.RoleId).Update(role); err != nil {
			return errors.Wrap(err, "更新角色信息失败")
		}
		// update role-menus
		if _, err := dbSession.Table("sys_role_menu").In("role_id", role.RoleId).Delete(); err != nil {
			return errors.Wrap(err, "删除角色关联菜单失败")
		}
		if len(role.MenuIds) > 0 {
			menuList := util.Convert(role.MenuIds, func(id int64) entity.SysRoleMenu {
				return entity.SysRoleMenu{RoleId: role.RoleId, MenuId: id}
			})
			if _, err := dbSession.Insert(menuList); err != nil {
				return errors.Wrap(err, "插入角色关联菜单失败")
			}
		}
		return nil
	})
}

func (rs *roleService) checkRoleNameUnique(roleName string, roleId int64) (bool, error) {
	return gb.DB.Table("sys_role").Where("role_name = ? and role_id <> ?", roleName, roleId).Exist()
}

func (rs *roleService) checkRoleKeyUnique(roleKey string, roleId int64) (bool, error) {
	return gb.DB.Table("sys_role").Where("role_key = ? and role_id <> ?", roleKey, roleId).Exist()
}

// 批量删除角色
func (rs *roleService) DeleteRoles(roleIds []int64) error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		for _, id := range roleIds {
			r := entity.SysRole{RoleId: id}
			if r.IsAdmin() {
				return errors.New("不允许删除超级管理员角色")
			}

			if exist, err := rs.checkRoleRsUser(id); err != nil {
				return err
			} else if exist {
				return errors.New("角色已分配用户,不允许删除")
			}
		}

		if _, err := dbSession.Table("sys_role_menu").In("role_id", roleIds).Delete(); err != nil {
			return errors.Wrap(err, "删除角色关联菜单失败")
		}
		if _, err := dbSession.Table("sys_role_dept").In("role_id", roleIds).Delete(); err != nil {
			return errors.Wrap(err, "删除角色关联部门失败")
		}
		if _, err := dbSession.Table("sys_role").In("role_id", roleIds).Delete(); err != nil {
			return errors.Wrap(err, "删除角色失败")
		}
		return nil
	})
}

func (rs *roleService) checkRoleRsUser(roleId int64) (bool, error) {
	return gb.DB.Table("sys_user_role").Where("role_id = ?", roleId).Exist()
}

// 角色状态修改
func (rs *roleService) ChangeStatus(role *entity.SysRole) error {
	if role.IsAdmin() {
		return errors.New("不允许操作超级管理员角色")
	}
	_, err := gb.DB.Table("sys_role").Cols("status").Where("role_id = ?", role.RoleId).Update(role)
	return err
}

// 修改角色数据权限范围
func (rs *roleService) ChangeRoleDataScope(role *entity.SysRole) error {
	if role.IsAdmin() {
		return errors.New("不允许操作超级管理员角色")
	}
	_, err := gb.DB.Table("sys_role").Cols("data_scope").Where("role_id = ?", role.RoleId).Update(role)
	return err
}

// 增加角色关联用户列表（为批量用户设置相同角色）
func (rs *roleService) AddRoleUsers(roleId int64, userIds []int64) error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		var userIdList = []int64{}
		if err := dbSession.Table("sys_user").Select("user_id").In("user_id", userIds).
			NotIn("user_id", builder.Select("u.user_id").From("sys_user u").
				InnerJoin("sys_user_role ur", "u.user_id = ur.user_id").
				Where(builder.Eq{"ur.role_id": roleId}.And(builder.In("u.user_id", userIds)))).
			Find(&userIdList); err != nil {
			return err
		}

		// 已由其他并发操作执行
		if len(userIdList) == 0 {
			return nil
		}

		urList := util.Convert(userIdList, func(id int64) entity.SysUserRole {
			return entity.SysUserRole{UserId: id, RoleId: roleId}
		})
		if _, err := dbSession.Insert(urList); err != nil {
			return err
		}
		return nil
	})
}

func (rs *roleService) DeleteRoleUsers(roleId int64, userIds []int64) error {
	_, err := gb.DB.Table("sys_user_role").
		Where("role_id = ?", roleId).
		In("user_id", userIds).
		Delete()
	return err
}
