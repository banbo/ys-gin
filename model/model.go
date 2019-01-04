package model

import (
	"math"
)

type Model struct {
}

// 列表
type ModelList struct {
	IsPage    bool        `json:"is_page"`    // 是否分页
	PageIndex int         `json:"page_index"` // 当前页码
	PageSize  int         `json:"page_size"`  // 每页数据条数
	PageCount int         `json:"page_count"` // 总页数
	Total     int         `json:"total"`      // 总数据条数
	Items     interface{} `json:"items"`      // 数据数组
}

// 分页列表
func (m *Model) Paging(pageIndex, pageSize, total int) (limit, offset int, modelList *ModelList) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	limit = pageSize
	offset = (pageIndex - 1) * pageSize

	modelList = &ModelList{
		IsPage:    true,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		PageCount: int(math.Ceil(float64(total) / float64(pageSize))),
		Total:     total,
	}

	return
}

// 不分页列表，也包装成一个ModelList
func (m *Model) NoPaging(total int, list interface{}) *ModelList {
	return &ModelList{
		IsPage:    false,
		PageIndex: 1,
		PageSize:  total,
		PageCount: 1,
		Total:     total,
		Items:     list,
	}

}
