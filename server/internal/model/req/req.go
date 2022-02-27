package req

type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type Find struct {
	Pagination
	Content string `form:"content"`
}
