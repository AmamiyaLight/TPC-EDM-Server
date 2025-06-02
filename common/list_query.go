package common

import (
	"TPC-EDM-Server/global"
	"fmt"
)

//list,count,err = common.ListQuery(models.LogModel,conf)

type PageInfo struct {
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
	Key   string `form:"key"`
	Order string `form:"order"`
}

func (p PageInfo) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}
func (p PageInfo) GetLimit() int {
	if p.Limit <= 0 || p.Limit > 100 {
		return 10
	}
	return p.Limit
}
func (p PageInfo) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

type Options struct {
	PageInfo      PageInfo
	Likes         []string
	PreLoads      []string
	Joins         []string
	Where         any
	Debug         bool
	DefaultOrder  string
	AccuracyCount bool
}

func ListQuery[T any](model T, option Options) (list []T, count int, err error) {

	//基础查询
	query := global.DB.Model(model).Where(model)

	if option.Debug {
		query = query.Debug()
	}

	if len(option.PreLoads) > 0 {
		for _, preload := range option.PreLoads {
			query = query.Preload(preload)
		}
	}
	joinMap := make(map[string]bool)
	if len(option.Joins) > 0 {
		for _, join := range option.Joins {
			if !joinMap[join] {
				query = query.Joins(join)
				joinMap[join] = true
			}
		}
	}
	if len(option.Likes) > 0 && option.PageInfo.Key != "" {
		likes := global.DB.Where("")
		for _, column := range option.Likes {
			likes.Or(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", option.PageInfo.Key))
		}
		query = query.Where(likes)
	}
	if option.Where != nil {
		query = query.Where(option.Where)
	}

	//排序
	if option.PageInfo.Order != "" {
		query = query.Order(option.PageInfo.Order)
	} else {
		if option.DefaultOrder != "" {
			query = query.Order(option.DefaultOrder)
		}
	}
	if option.AccuracyCount {
		var _c int64
		if err = query.Count(&_c).Error; err != nil {
			return nil, 0, err
		}
		count = int(_c)
	} else {
		var approxCount int64
		row := struct{ Rows int64 }{}
		stmt := global.DB.Model(model).Statement
		if stmt.Table == "" {
			if err := stmt.Parse(model); err != nil {
				return nil, 0, err
			}
		}
		tableName := stmt.Table
		quotedTableName := stmt.Quote(tableName) // 正确引用表名
		err = global.DB.Raw(fmt.Sprintf("EXPLAIN SELECT * FROM %s", quotedTableName)).Scan(&row).Error
		approxCount = row.Rows
		count = int(approxCount)
	}
	//分页
	limit := option.PageInfo.GetLimit()
	offset := option.PageInfo.GetOffset()
	err = query.Offset(offset).Limit(limit).Find(&list).Error
	return
}
