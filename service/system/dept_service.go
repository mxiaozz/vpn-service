package system

import (
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/util"
	"xorm.io/xorm"
)

var DeptService = &deptService{}

type deptService struct {
}

// 获取部门信息，不存在时返回 nil
func (ds *deptService) GetDept(deptId int64) (*entity.SysDept, error) {
	var dept entity.SysDept
	if exist, err := gb.DB.Where("dept_id = ?", deptId).Get(&dept); err != nil {
		return nil, err
	} else if !exist {
		return nil, nil
	}
	return &dept, nil
}

func (ds *deptService) GetDepts(deptIds ...int64) ([]entity.SysDept, error) {
	var dept = []entity.SysDept{}
	if err := gb.DB.Table("sys_dept").In("dept_id", deptIds).Find(&dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (ds *deptService) GetDeptTree() ([]entity.SysDept, error) {
	if depts, err := ds.selectDeptList(&entity.SysDept{}); err != nil {
		return nil, err
	} else {
		return ds.buildDeptTree(depts, 0), nil
	}
}

func (ds *deptService) GetDeptList(dept *entity.SysDept) ([]entity.SysDept, error) {
	return ds.selectDeptList(dept)
}

func (ds *deptService) selectDeptList(dept *entity.SysDept) ([]entity.SysDept, error) {
	dbSession := gb.DB.Where("del_flag = 0")
	if dept.DeptId != 0 {
		dbSession.And("dept_id = ?", dept.DeptId)
	}
	if dept.ParentId != 0 {
		dbSession.And("parent_id = ?", dept.ParentId)
	}
	if dept.DeptName != "" {
		dbSession.And("dept_name like ?", "%"+dept.DeptName+"%")
	}
	if dept.Status != "" {
		dbSession.And("status = ?", dept.Status)
	}
	dbSession.Asc("parent_id", "order_num")

	var depts = []entity.SysDept{}
	if err := dbSession.Find(&depts); err != nil {
		return nil, err
	}
	return depts, nil
}

func (ds *deptService) buildDeptTree(depts []entity.SysDept, parentId int64) []entity.SysDept {
	list := make([]entity.SysDept, 0)
	for _, dept := range depts {
		if dept.ParentId == parentId {
			c := ds.buildDeptTree(depts, dept.DeptId)
			dept.Children = c
			// tree select
			dept.Id = dept.DeptId
			dept.Label = dept.DeptName
			list = append(list, dept)
		}
	}
	return list
}

func (ds *deptService) GetDeptListByRoleId(roleId int64) ([]int64, error) {
	role, err := RoleService.GetRole(roleId)
	if err != nil {
		return nil, err
	}

	dbSession := gb.DB.Table("sys_dept").Alias("d").
		Join("left", []string{"sys_role_dept", "rd"}, "d.dept_id = rd.dept_id").
		Where("rd.role_id = ?", role.RoleId)
	if role.DeptCheckStrictly {
		dbSession.And(`d.dept_id not in (
			select d.parent_id from sys_dept d 
			inner join sys_role_dept rd on d.dept_id = rd.dept_id and rd.role_id = ?
		)`, role.RoleId)
	}
	dbSession.OrderBy("d.parent_id, d.order_num")

	var deptIds = []int64{}
	err = dbSession.Find(&deptIds)
	return deptIds, err
}

func (ds *deptService) AddDept(dept *entity.SysDept) error {
	parent, err := ds.GetDept(dept.ParentId)
	if err != nil {
		return errors.Wrap(err, "查询上级部门失败")
	} else if parent == nil {
		return errors.New("上级部门不存在")
	}
	dept.Ancestors = parent.Ancestors + "," + strconv.FormatInt(parent.DeptId, 10)

	if exist, err := ds.checkDeptNameUnique(dept.DeptName, dept.ParentId, dept.DeptId); err != nil {
		return errors.Wrap(err, "查询部门名称是否唯一失败")
	} else if exist {
		return errors.New("部门名称已存在")
	}

	_, err = gb.DB.Insert(dept)
	return err
}

func (ds *deptService) checkDeptNameUnique(deptName string, parentId, deptId int64) (bool, error) {
	return gb.DB.Table("sys_dept").Where("dept_name = ? and parent_id = ? and dept_id <> ?", deptName, parentId, deptId).Exist()
}

func (ds *deptService) GetDeptsExcludeChild(deptId int64) ([]entity.SysDept, error) {
	depts, err := ds.GetDeptList(&entity.SysDept{})
	if err != nil {
		return nil, err
	}

	deptIdStr := strconv.FormatInt(deptId, 10)
	depts = util.NewList(depts).Filter(func(d entity.SysDept) bool {
		if d.DeptId == deptId {
			return false
		}
		ancestors := strings.Split(d.Ancestors, ",")
		return !slices.Contains(ancestors, deptIdStr)
	})
	return depts, nil
}

func (ds *deptService) UpdateDept(dept *entity.SysDept) error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		// 检查
		if exist, err := ds.checkDeptNameUnique(dept.DeptName, dept.ParentId, dept.DeptId); err != nil {
			return errors.Wrap(err, "查询部门名称是否唯一失败")
		} else if exist {
			return errors.New("部门名称已存在")
		}

		// 说明修改的是顶级部门(一般是单位组织)
		// 其他部门如果要修改上级部门，其 parentId > 0
		if dept.ParentId == 0 {
			_, err := dbSession.Where("dept_id = ?", dept.DeptId).Update(dept)
			return err
		}

		// 修改子部门，要校验 parentId 是否存在
		parent, err := ds.GetDept(dept.ParentId)
		if err != nil {
			return errors.Wrap(err, "查询上级部门失败")
		} else if parent == nil {
			return errors.New("上级部门不存在")
		}

		// 修正子部门 ancestors
		dbDept, err := ds.GetDept(dept.DeptId)
		if err != nil {
			return err
		}
		deptIdStr := strconv.FormatInt(dept.DeptId, 10)
		newParentIdStr := strconv.FormatInt(dept.ParentId, 10)
		oldChildrenAncestors := dbDept.Ancestors + "," + deptIdStr
		newChildrenAncestors := parent.Ancestors + "," + newParentIdStr + "," + deptIdStr
		dept.Ancestors = parent.Ancestors + "," + newParentIdStr

		if children, err := ds.getDeptChild(dept.DeptId); err != nil {
			return err
		} else if len(children) > 0 {
			for _, d := range children {
				d.Ancestors = strings.Replace(d.Ancestors, oldChildrenAncestors, newChildrenAncestors, -1)
				if _, err := dbSession.Where("dept_id = ?", d.DeptId).Update(d); err != nil {
					return err
				}
			}
		}

		// 更新部门信息
		_, err = dbSession.Where("dept_id = ?", dept.DeptId).Update(dept)
		return err
	})
}

func (ds *deptService) getDeptChild(deptId int64) ([]*entity.SysDept, error) {
	depts, err := ds.GetDeptList(&entity.SysDept{})
	if err != nil {
		return nil, err
	}

	deptIdStr := strconv.FormatInt(deptId, 10)
	depts = util.NewList(depts).Filter(func(d entity.SysDept) bool {
		ancestors := strings.Split(d.Ancestors, ",")
		return slices.Contains(ancestors, deptIdStr)
	})
	return util.Convert(depts, func(d entity.SysDept) *entity.SysDept {
		return &d
	}), nil
}

func (ds *deptService) DeleteDept(deptId int64) error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		if exist, err := gb.DB.Table("sys_user").Where("dept_id = ?", deptId).Exist(); err != nil {
			return err
		} else if exist {
			return errors.New("部门存在关联用户，不能删除")
		}

		if _, err := dbSession.Table("sys_role_dept").Where("dept_id = ?", deptId).Delete(); err != nil {
			return errors.Wrap(err, "删除部门角色关联失败")
		}

		if _, err := dbSession.Table("sys_dept").Where("dept_id = ?", deptId).Delete(); err != nil {
			return errors.Wrap(err, "删除部门失败")
		}
		return nil
	})
}
