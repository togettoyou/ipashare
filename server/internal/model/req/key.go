package req

type KeyQuery struct {
	Username string `json:"username" binding:"required,min=4,max=16"`
}

type KeyCr struct {
	Username string `json:"username" binding:"required,min=4,max=16"`
	Password string `json:"password" binding:"required,min=4,max=16"`
	Num      int    `json:"num"`
}

type KeyUp struct {
	Username string `json:"username" binding:"required,min=4,max=16"`
	Num      int    `json:"num"`
}
