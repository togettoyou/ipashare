package resp

// Pagination 分页结构体
type Pagination struct {
	//每页显示的数量
	PageSize int `json:"page_size"`
	//当前页码
	Page int `json:"page"`
	//分页的数据内容
	Data interface{} `json:"data"`
	//全部的页码数量
	Total int64 `json:"total"`
}
