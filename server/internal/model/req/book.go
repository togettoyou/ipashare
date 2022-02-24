package req

type Book struct {
	Name string `json:"name" binding:"required,min=10"`
	Url  string `json:"url" binding:"required,url"`
}
