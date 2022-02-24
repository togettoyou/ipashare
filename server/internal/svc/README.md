# svc 即 service

业务逻辑层，负责处理复杂业务逻辑，处于 api 层和 dao 层之间。svc 只能通过 dao 层获取数据

开发规范：

1. `svc` 目录下直接存放业务处理文件，例：`svc/book.go`

开发教程：

1. 创建业务对象结构体，必须内嵌 `Service` 结构体

   ```go
   type Book struct {
   	Service
   }
   ```

2. 方法声明，推荐使用指针接收器，避免不必要的拷贝，业务层的错误必须使用 `pkg/e` 包实现

   ```go
   func (b *Book) GetList() ([]model.Book, error) {
   	b.log.Info("业务处理")
   	// 使用store调用dao层
   	books, err := b.store.Book.List()
   	if err != nil {
   		// 返回错误包含调用栈信息
   		return nil, e.NewWithStack(e.DBError, err)
   	}
   	return books, nil
   }
   ```
   

