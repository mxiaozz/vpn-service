package system

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/util"
	"xorm.io/builder"
	"xorm.io/xorm"
)

var UserService = &userService{}

type userService struct {
}

// 获取所有用户
func (us *userService) GetAllUsers() ([]entity.SysUser, error) {
	var users = []entity.SysUser{}
	err := gb.DB.Select("*").Where("del_flag = '0'").Find(&users)
	return users, err
}

// 用户管理（搜索/分页）
func (us *userService) GetUserListPage(user *entity.SysUser, page *model.Page[*entity.SysUser]) error {
	err := gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("u.user_id, u.dept_id, u.nick_name, u.user_name, u.email, u.avatar, u.phonenumber, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, u.valid_day").
			From("sys_user", "u").
			Where(
				builder.Eq{"u.del_flag": "0"}.
					And(builder.If(user.UserId != 0, builder.Eq{"u.user_id": user.UserId})).
					And(builder.If(user.UserName != "", builder.Like{"u.user_name", user.UserName})).
					And(builder.If(user.Status != "", builder.Eq{"u.status": user.Status})).
					And(builder.If(user.Phonenumber != "", builder.Like{"u.phonenumber", user.Phonenumber})).
					And(builder.If(func() bool { return user.Params["beginTime"] != nil }(),
						builder.Gte{"u.create_time": user.Params["beginTime"]})).
					And(builder.If(func() bool { return user.Params["endTime"] != nil }(),
						builder.Lte{"u.create_time": user.Params["endTime"]})).
					And(builder.If(user.DeptId != 0,
						builder.If(gb.DB.DriverName() == "sqlite3",
							builder.Expr(`u.dept_id = ? or u.dept_id in ( 
									select t.dept_id from sys_dept t where t.ancestors like (
										select ancestors||','||?||'%' from sys_dept where dept_id = ?
								))`, user.DeptId, user.DeptId, user.DeptId),
							builder.Expr(`u.dept_id = ? or u.dept_id in ( 
									select t.dept_id from sys_dept t where t.ancestors like (
										select concat(ancestors, ',', ?, '%') from sys_dept where dept_id = ?
								))`, user.DeptId, user.DeptId, user.DeptId)))))
		return builder.Expr("user_id")
	})
	if err != nil {
		return err
	}

	// 加载部门信息
	list := util.NewList(page.Rows).
		MapToInt64(func(user *entity.SysUser) int64 { return user.DeptId }).
		Distinct(func(deptId int64) any { return deptId })
	depts, err := DeptService.GetDepts(list...)
	if err != nil {
		return err
	}
	util.NewList(page.Rows).ForEach(func(user *entity.SysUser) {
		deptList := util.NewList(depts).Filter(func(dept entity.SysDept) bool { return dept.DeptId == user.DeptId })
		if deptList.Count() > 0 {
			user.Dept = &deptList[0]
		}
	})

	return nil
}

// 获取用户信息，不包含包含部门和角色
func (us *userService) GetSysUsers(userIds []int64) ([]entity.SysUser, error) {
	var users = []entity.SysUser{}
	err := gb.DB.Table("sys_user").In("user_id", userIds).Find(&users)
	return users, err
}

// 获取用户信息，isWithExtInfo=true 包含部门和角色列表（用户不存在时返回 nil）
func (us *userService) GetSysUser(userName string, isWithExtInfo bool) (*entity.SysUser, error) {
	return us.selectUser(map[string]any{"user_name": userName}, isWithExtInfo)
}

// 获取用户信息，isWithExtInfo=true 包含部门和角色列表（用户不存在时返回 nil）
func (us *userService) GetSysUserById(userId int64, isWithExtInfo bool) (*entity.SysUser, error) {
	return us.selectUser(map[string]any{"user_id": userId}, isWithExtInfo)
}

func (us *userService) selectUser(wh map[string]any, isWithExtInfo bool) (*entity.SysUser, error) {
	var user entity.SysUser
	if exist, err := gb.DB.Select("*").Where(wh).Desc("create_time").Limit(1).Get(&user); err != nil || !exist {
		return nil, errors.Wrap(err, "读取用户信息失败")
	}

	if isWithExtInfo {
		if err := us.getUserExtInfo(&user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (us *userService) getUserExtInfo(user *entity.SysUser) error {
	dept, err := DeptService.GetDept(user.DeptId)
	if err != nil {
		return errors.Wrap(err, "读取用户部门信息失败")
	}
	user.Dept = dept

	roles, err := RoleService.GetUserRoles(user.UserId)
	if err != nil {
		return errors.Wrap(err, "读取用户角色列表失败")
	}
	user.Roles = roles
	return nil
}

// 更新用户证书剩余有效期
func (us *userService) UpdateUserCertValidDay(user *entity.SysUser) error {
	_, err := gb.DB.Cols("valid_day").Where("user_name = ?", user.UserName).Update(user)
	return err
}

// 增加用户（包含关联岗位/角色）
func (us *userService) AddUser(user *entity.SysUser) error {
	// 加密密码
	pdata, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(pdata)

	return gb.Tx(func(dbSession *xorm.Session) error {
		if exist, err := us.isExistUserName(dbSession, user.UserName); err != nil {
			return errors.Wrap(err, "检查用户名是否已存在失败")
		} else if exist {
			return errors.New("用户名已存在")
		}

		if _, err := dbSession.InsertOne(user); err != nil {
			return errors.Wrap(err, "增加用户信息失败")
		}
		if err := us.addUserPosts(dbSession, user.UserId, user.PostIds); err != nil {
			return errors.Wrap(err, "增加用户关联岗位失败")
		}
		if err := us.addUserRoles(dbSession, user.UserId, user.RoleIds); err != nil {
			return errors.Wrap(err, "增加用户关联角色失败")
		}
		return nil
	})
}

func (us *userService) isExistUserName(db *xorm.Session, userName string) (bool, error) {
	return db.Table("sys_user").Where("user_name = ?", userName).Exist()
}

func (us *userService) UpdateUserLoginInfo(user *entity.SysUser) error {
	_, err := gb.DB.Table("sys_user").Cols("login_ip", "login_date").
		Where("user_id = ?", user.UserId).Update(user)
	return err
}

// 更新用户信息（包含关联岗位/角色）
func (us *userService) UpdateUser(user *entity.SysUser) error {
	if util.IsAdminId(user.UserId) {
		return errors.New("不允许修改超级管理员信息")
	}

	return gb.Tx(func(dbSession *xorm.Session) error {
		// 不允许修改用户名
		user.UserName = ""

		if _, err := dbSession.Table("sys_user").
			Cols("dept_id", "nick_name", "email", "phonenumber", "sex", "status", "update_by", "update_time", "remark").
			Where("user_id = ?", user.UserId).
			Update(user); err != nil {
			return errors.Wrap(err, "更新用户信息失败")
		}

		// update posts
		if err := us.delUserPosts(dbSession, user.UserId); err != nil {
			return errors.Wrap(err, "删除用户关联岗位失败")
		}
		if err := us.addUserPosts(dbSession, user.UserId, user.PostIds); err != nil {
			return errors.Wrap(err, "增加用户关联岗位失败")
		}

		// update roles
		if err := us.delUserRoles(dbSession, user.UserId); err != nil {
			return errors.Wrap(err, "删除用户关联角色失败")
		}
		if err := us.addUserRoles(dbSession, user.UserId, user.RoleIds); err != nil {
			return errors.Wrap(err, "增加用户关联角色失败")
		}
		return nil
	})
}

func (us *userService) addUserPosts(db *xorm.Session, userId int64, postIds []int64) error {
	if len(postIds) == 0 {
		return nil
	}

	postIds = util.NewList(postIds).Distinct(func(postId int64) any { return postId })
	upList := util.Convert(postIds, func(postId int64) entity.SysUserPost {
		return entity.SysUserPost{UserId: userId, PostId: postId}
	})

	_, err := db.Insert(upList)
	return err
}

func (us *userService) delUserPosts(db *xorm.Session, userIds ...int64) error {
	_, err := db.Table("sys_user_post").In("user_id", userIds).Delete()
	return err
}

func (us *userService) addUserRoles(db *xorm.Session, userId int64, roleIds []int64) error {
	if len(roleIds) == 0 {
		return nil
	}

	roleIds = util.NewList(roleIds).Distinct(func(roleId int64) any { return roleId })
	urList := util.Convert(roleIds, func(roleId int64) entity.SysUserRole {
		return entity.SysUserRole{UserId: userId, RoleId: roleId}
	})

	_, err := db.Insert(urList)
	return err
}

func (us *userService) delUserRoles(db *xorm.Session, userIds ...int64) error {
	_, err := db.Table("sys_user_role").In("user_id", userIds).Delete()
	return err
}

// 删除用户（包含删除岗位/角色关联信息）
func (us *userService) DeleteUser(userIds []int64) error {
	for _, id := range userIds {
		if util.IsAdminId(id) {
			return errors.New("不允许删除超级管理员")
		}
	}

	return gb.Tx(func(dbSession *xorm.Session) error {
		if err := us.delUserPosts(dbSession, userIds...); err != nil {
			return errors.Wrap(err, "删除用户关联岗位失败")
		}
		if err := us.delUserRoles(dbSession, userIds...); err != nil {
			return errors.Wrap(err, "删除用户关联角色失败")
		}
		if _, err := dbSession.In("user_id", userIds).Delete(&entity.SysUser{}); err != nil {
			return errors.Wrap(err, "删除用户失败")
		}
		return nil
	})
}

// 重置密码
func (us *userService) ResetPassword(user *entity.SysUser) error {
	if util.IsAdminId(user.UserId) {
		return errors.New("不允许修改超级管理员密码")
	}

	pdata, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(pdata)

	_, err := gb.DB.Table("sys_user").Cols("password").Where("user_id = ?", user.UserId).Update(user)
	return err
}

// 变更状态
func (us *userService) ChangeStatus(user *entity.SysUser) error {
	if util.IsAdminId(user.UserId) {
		return errors.New("不允许修改超级管理员状态")
	}

	_, err := gb.DB.Table("sys_user").Cols("status").Where("user_id = ?", user.UserId).Update(user)
	return err
}

// 更新用户关联的角色信息
func (us *userService) ChangeUserRoles(userId int64, roleIds []int64) error {
	if util.IsAdminId(userId) {
		return errors.New("不允许修改超级管理员角色")
	}

	return gb.Tx(func(dbSession *xorm.Session) error {
		if err := us.delUserRoles(dbSession, userId); err != nil {
			return errors.Wrap(err, "删除用户关联角色失败")
		}
		if err := us.addUserRoles(dbSession, userId, roleIds); err != nil {
			return errors.Wrap(err, "变更用户关联角色失败")
		}
		return nil
	})
}

// 角色管理，分配用户
func (us *userService) GetRoleUserPage(roleId int64, user *entity.SysUser, page *model.Page[entity.SysUser]) error {
	// 查询语句
	sql := builder.Dialect(string(gb.DB.Dialect().URI().DBType)).
		Select("*").From(
		builder.Select("distinct u.*").From("sys_user", "u").
			LeftJoin("sys_dept d", "u.dept_id = d.dept_id").
			LeftJoin("sys_user_role ur", "u.user_id = ur.user_id").
			LeftJoin("sys_role r", "r.role_id = ur.role_id").
			Where(builder.Expr("u.del_flag = '0'").
				And(builder.Eq{"r.role_id": roleId}).
				And(builder.If(user.UserName != "", builder.Like{"u.user_name", user.UserName})).
				And(builder.If(user.Phonenumber != "", builder.Like{"u.phonenumber", user.Phonenumber}))),
		"t")

	// 查询总数
	var err error
	if page.Total, err = gb.DB.SQL(sql.Select("count(*)")).Count(); err != nil {
		return err
	} else if page.Total == 0 {
		return nil
	}

	// 分页查询记录
	if page.OrderByColumn != "" {
		sql.OrderBy(page.OrderByColumn)
	} else {
		sql.OrderBy("user_id")
	}
	sql.Limit(page.PageSize, page.Offset)

	if err := gb.DB.SQL(sql.Select("*")).Find(&page.Rows); err != nil {
		return err
	}
	return nil
}

// 角色管理，分配用户
func (us *userService) GetNotRoleUserPage(roleId int64, user *entity.SysUser, page *model.Page[entity.SysUser]) error {
	// 查询语句
	sql := builder.Dialect(string(gb.DB.Dialect().URI().DBType)).
		Select("*").From("sys_user").Where(builder.Expr("user_id > 1").
		And(builder.NotIn("user_id", builder.Expr(`select u.user_id from sys_user u
		left join sys_user_role ur on ur.user_id = u.user_id
		where ur.role_id = ?`, roleId))))

	// 查询总数
	var err error
	if page.Total, err = gb.DB.SQL(sql.Select("count(*)")).Count(); err != nil {
		return err
	} else if page.Total == 0 {
		return nil
	}

	// 分页查询记录
	if page.OrderByColumn != "" {
		sql.OrderBy(page.OrderByColumn)
	} else {
		sql.OrderBy("user_id")
	}
	sql.Limit(page.PageSize, page.Offset)

	if err := gb.DB.SQL(sql.Select("*")).Find(&page.Rows); err != nil {
		return err
	}
	return nil
}
