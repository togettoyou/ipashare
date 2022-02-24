# router

负责 api 路由注册

开发规范：

1. `router` 目录下直接存放路由文件，例：`router/book.go`

开发教程：

1. 实例化 api 层结构体，传入 store，日志也是在这里定义模块划分的

   ```go
   book := v1beta1.Book{
       Base: api.New(store, log.New("book").L()),
   }
   ```

   

2. 路由注册

   ```go
   r.GET("", book.GetList)
   ```