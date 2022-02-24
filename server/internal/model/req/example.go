package req

type UriArgs struct {
	ID uint `json:"id" uri:"id" binding:"required,min=10"`
}

type QueryArgs struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

type FormArgs struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

type JSONBody struct {
	Email    string `json:"email" binding:"required,email" example:"admin@qq.com"`
	Username string `json:"username" binding:"required,checkUsername" example:"admin"`
}

type ErrArgs struct {
	ID uint `json:"id" uri:"id" binding:"required,oneof=1 2 3 4 5"`
}
