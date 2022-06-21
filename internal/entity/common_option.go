package entity

import (
// "time"
)

// ================================================================
// ================================================================ Struct
type ListOption struct {
	// offset
	Page    int64 `json:"page,omitempty" form:"page" uri:"page"`
	PerPage int64 `json:"per_page,omitempty" form:"per_page" uri:"per_page"`
	Offset  int64 `json:"offset,omitempty" form:"offset" uri:"offset"`

	// order by
	OrderBy string `json:"order_by,omitempty" form:"order_by" uri:"order_by"`
	Order   string `json:"order,omitempty" form:"order" uri:"order"`
}

type ListResponse struct {
	TotalCount int64 `json:"total_count"`
	TotalPage  int64 `json:"total_page"`
	Page       int64 `json:"page"`
	PerPage    int64 `json:"per_page"`
}

// type JSONTime struct {
// 	time.Time
// }

// const TimeFormat = "2006-01-02T15:04:05Z07:00" // 这是个奇葩,必须是这个时间点, 据说是go诞生之日

// ================================================================
// ================================================================ Function
func (option *ListOption) CalcOffset() {
	if option.Page == 0 || option.PerPage == 0 {
		option.Offset = 0
	}
	option.Offset = (option.Page - 1) * option.PerPage
}

func (response *ListResponse) CalcTotalPage() {
	if response.PerPage == 0 {
		response.TotalPage = 1
		return
	}
	response.TotalPage = response.TotalCount / response.PerPage
	if response.TotalCount%response.PerPage != 0 {
		response.TotalPage = response.TotalPage + 1
	}
}

// func (t *JSONTime) MarshalJSON() ([]byte, error) {
// 	return []byte(t.Format(TimeFormat)), nil
// }

// func (t *JSONTime) UnmarshalJSON(data []byte) error {
// 	var err error
// 	t.Time, err = time.Parse(TimeFormat, string(data))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
