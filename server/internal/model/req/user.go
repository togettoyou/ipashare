package req

type UserPW struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,len=32"`
}
