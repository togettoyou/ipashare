package req

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,len=32"`
}

type UserPW struct {
	Password string `json:"password" binding:"required,len=32"`
}
