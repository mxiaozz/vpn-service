package gb

import (
	"strings"

	"github.com/pkg/errors"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/util"
	"xorm.io/builder"
	"xorm.io/xorm/schemas"
)

func SelectPage[T any](page *model.Page[T], action func(sql *builder.Builder) builder.Cond) error {
	subSql := builder.Select()
	mainSql := builder.Dialect(string(DB.Dialect().URI().DBType)).
		Select("*").From(subSql, "t")

	// 业务 select sql
	order := action(subSql)

	// 查询总数
	var err error
	if page.Total, err = DB.SQL(mainSql.Select("count(*)")).Count(); err != nil {
		return errors.Wrap(err, "分页查询总数失败")
	} else if page.Total == 0 {
		return nil
	}

	// order by
	customOrderBy(page, mainSql, order)

	var reasonable = true
REASONABLE:
	// limit
	mainSql.Limit(page.PageSize, page.Offset)

	// query
	if err := DB.SQL(mainSql.Select("*")).Find(&page.Rows); err != nil {
		return errors.Wrap(err, "分页查询记录失败")
	}

	if reasonable && page.Total > 0 && len(page.Rows) == 0 {
		reasonable = false

		num := page.Total / int64(page.PageSize)
		if page.Total%int64(page.PageSize) != 0 {
			num++
		}
		page.PageNum = int(num)
		page.Offset = int(num-1) * page.PageSize

		goto REASONABLE
	}

	return nil
}

func customOrderBy[T any](page *model.Page[T], mainSql *builder.Builder, order builder.Cond) {
	if page.OrderByColumn == "" {
		mainSql.OrderBy(order)
		return
	}

	// 默认排序
	columnName := page.OrderByColumn
	mode := "asc"

	// 从对象中映射字段
	if table, err := DB.TableInfo(new(T)); err != nil {
		Logger.Error(err)
	} else {
		list := util.NewList(table.Columns()).Filter(func(c *schemas.Column) bool {
			return strings.EqualFold(c.FieldName, page.OrderByColumn)
		})
		if len(list) > 0 {
			columnName = list[0].Name
		}
	}

	// 排序方式
	if len(page.IsAsc) >= 3 {
		flag := page.IsAsc[0:3]
		if strings.EqualFold("des", flag) {
			mode = "desc"
		}
	}

	mainSql.OrderBy(columnName + " " + mode)
}
