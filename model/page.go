package model

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

type Page[T any] struct {
	PageNum       int    `json:"pageNum" form:"pageNum"`             // 当前页面
	PageSize      int    `json:"pageSize" form:"pageSize"`           // 每页数量
	OrderByColumn string `json:"orderByColumn" form:"orderByColumn"` // 排序字段
	IsAsc         string `json:"isAsc" form:"isAsc"`

	Offset int   `json:"-" form:"-"` // 自动根据 PageNum & PageSize 计算偏移量
	Total  int64 `json:"-" form:"-"` // 总记录数
	Rows   []T   `json:"-" form:"-"` // 返回数据
}

func NewPage[T any](ctx *gin.Context) (*Page[T], error) {
	// 缓存流支持多次读取
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var page Page[T]
	if ctx.Request.Method == "GET" {
		if err = ctx.ShouldBindQuery(&page); err != nil {
			return nil, err
		}
	} else {
		if err = ctx.ShouldBindJSON(&page); err != nil {
			return nil, err
		}
	}
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 10
	}
	page.Offset = (page.PageNum - 1) * page.PageSize
	page.Rows = make([]T, 0, 0)

	// 支持后续流读取
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return &page, nil
}

func (p *Page[T]) ToMap() map[string]any {
	m := make(map[string]any, 2)
	m["rows"] = p.Rows
	m["total"] = p.Total
	return m
}
