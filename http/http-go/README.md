### 序言

HTTP-GO(简称：ghttp) 主要是对golang 操作http 请求的一些简要封装，使用HTTP-GO 可以很方便的对http 的请求、响应等做操作。目前具备如下几种功能：

- 支持http 常用请求
- 支持query 及body 参数
- 参数支持interface 结构
- 支持灵活定制参数
- 支持灵活添加http 请求头
- 支持cookie
- 支持超时设置
- 支持响应数据转换为Struct
- 支持代理设置  
  具体使用方法可以参照Test 文件，里面列举了几种使用方法及情况。

### GET

HTTP-GO 默认为GET 请求，基本请求如下所示：

```
import ghttp "https://github.com/psoKnight/go-common/tree/main/http/http-go"

res, err := ghttp.Request{
	Url:   "http://127.0.0.1:8080",
}.Do()
```

### POST

HTTP-GO 使用POST 请求与golang 请求一致，基本请求如下所示：

```
import ghttp "https://github.com/psoKnight/go-common/tree/main/http/http-go"

res, err := ghttp.Request{
	Method: "POST",
	Url:   "http://127.0.0.1:8080",
}.Do()
```

### 请求参数

HTTP-GO 参数支持interface 类型，允许直接将Struct 作为参数赋值传递，同时支持tags 配置指定参数名称以及忽略参数空值，使用如下所示。

#### 使用Struct 作为参数

直接使用struct 作为参数时，默认元素小写后作为url 参数，下面例子请求后生成url 如request 所示：

```  
import ghttp "https://github.com/psoKnight/go-common/tree/main/http/http-go"

// request-> Get http://127.0.0.1:8080?name=xc&password=xc
type User struct {
	Name     string
	Password string
}
user := User{
	Name:     "xc",
	Password: "xc",
}
res, err := ghttp.Request{
	Url:   "http://127.0.0.1:8080",
	Query: user,
}.Do()
```

#### Tags 配置参数

HTTP-GO 支持将参数在结构体中配置为指定名称，关键字"-"表示此参数忽略不会拼接到url 中，"omitempty"关键词表示该字段为空时不做拼接，可参考下面例子生成的url。

``` 
import ghttp "https://github.com/psoKnight/go-common/tree/main/http/http-go"

// Get http://127.0.0.1:8080?name=xc
type User struct {
	Name     string `json:"name"`
	Password string `json:"-"`
	Sex      string `json:"sex,omitempty"`
}
user := User{
	Name:     "xc",
	Password: "xc",
	Sex:      "",
}
res, err := ghttp.Request{
	Url:   "http://127.0.0.1:8080",
	Query: user,
}.Do()
```

#### 结构体嵌套

HTTP-GO 支持结构体嵌套的方式拼接参数，这应该也是最为常见的一种方式，例子如下：

```json
``` 

### Header

HTTP-GO 支持Head 添加处理，如下所示：

```
import ghttp "https://github.com/psoKnight/go-common/tree/main/http/http-go"

type User struct {
	Name     string
	Password string
}
user := User{
	Name:     "xc",
	Password: "xc",
}
req := &ghttp.Request{
	Method:      "POST",
	Url:         "http://127.0.0.1:8080",
	Query:       user,
	ContentType: "application/json",
}
req.AddHeader("X-Custom", "haha")
res, err := req.Do()
```

### Cookie

```
import ghttp "https://github.com/psoKnight/go-common/tree/main/http/http-go"

res, err := ghttp.Request{
	Url:     "http://www.baidu.com",
	Timeout: 100 * time.Millisecond,
}.Do()
```

### 响应数据结构转换-Struct

```
import ghttp "https://github.com/psoKnight/go-common/tree/main/http/http-go"

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Sex      string `json:"sex"`
}
var user User
res, err := ghttp.Request{
	Url:     "http://127.0.0.1:8080",
	Timeout: 100 * time.Millisecond,
}.Do()
if err != nil {
	log.Panicln(err)
}
log.Println(res.Body.ToJson(&user))
```

### Proxy

```
import ghttp "https://github.com/psoKnight/go-common/tree/main/http/http-go"

res, err := ghttp.Request{
	Url:     "http://127.0.0.1:8080",
	Timeout: 100 * time.Millisecond,
	Proxy:   "http://127.0.0.1:8088",
}.Do()
```
