# api

api 接口层，负责参数校验解析，请求响应处理，不应处理复杂逻辑

开发规范：

1. 创建版本号，例：`v1beta1`
2. 版本号包下直接存放 api 文件，例：`v1beta1/book.go`

开发教程：

1. 创建对应 api 的结构体（通常为文件名首字母大写），必须内嵌 `api.Base` 结构体

   ```go
   type Book struct {
   	api.Base
   }
   ```


2. 方法声明，必须使用值接收器！！！

   ```go
   // bad 不能使用指针接收器，对于每一个请求都应当拷贝一份对象进行处理，避免并发情况下 MakeContext 和 MakeService 方法初始化异常
   func (g *Book) GetList(c *gin.Context) {
   	var bookSvc svc.Book
   	g.MakeContext(c).MakeService(&bookSvc.Service)
   	g.Log.Info("路由处理")
   	books, err := bookSvc.GetList()
   	if g.HasErr(err) {
   		return
   	}
   	g.OK(books)
   }
   
   // good 正确做法
   func (g Book) GetList(c *gin.Context) {
   	var bookSvc svc.Book
   	g.MakeContext(c).MakeService(&bookSvc.Service)
   	g.Log.Info("路由处理")
   	books, err := bookSvc.GetList()
   	if g.HasErr(err) {
   		return
   	}
   	g.OK(books)
   }
   ```