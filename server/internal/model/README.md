# model

数据模型层

- 根目录：负责结构体-表映射，DB 操作接口定义，使用 Store 统一管理所有 DB 操作接口实例

- req：request model struct 定义

- resp：response model struct 定义

model struct 串联着各个层

开发规范：

1. `model` 目录下直接存放数据库 orm 模型文件，例：`model/book.go`


2. `model/req` 目录下存放请求参数结构体，使用 `binding` tag 完成参数校验绑定，例：`model/req/book.go`


3. `model/resp` 目录下存放响应结果结构体，例：`model/resp/book.go`