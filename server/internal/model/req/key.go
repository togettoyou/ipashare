package req

type KeyQuery struct {
	Username string `form:"username" binding:"required,min=4,max=50"`
}

type KeyCr struct {
	Username string `json:"username" binding:"required,min=4,max=50"`
	Password string `json:"password" binding:"required,min=4,max=50"`
	Num      int    `json:"num"`
}

type KeyUp struct {
	Username string `json:"username" binding:"required,min=4,max=50"`
	Num      int    `json:"num"`
}
