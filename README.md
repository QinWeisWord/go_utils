# go_utils 常用静态方法工具库

一个面向日常开发的 Go 工具库，覆盖字符串、时间与环境、文件与 JSON、摘要与随机、集合泛型、网络与常见格式校验等功能。库内所有函数均含中文注释（参数、返回值、关键步骤），文件行数控制在 300 行以内，适合快速集成与二次开发。

## 版本要求

- Go 1.18+（使用了泛型）

## 安装与引用

本仓库已在当前目录初始化为模块：`go_utils`，包名为 `utils`。

- 本仓库内其他代码直接引用：

```go
import utils "go_utils"
```

- 在其他项目中以本地路径引用（推荐）：在目标项目 `go.mod` 中添加：

```go
require go_utils v0.0.0
replace go_utils => ../golang/go_utils
```

然后在代码中：

```go
import utils "go_utils"
```

> 说明：若未来托管到远程仓库（如 GitHub），可直接 `go get <repo>/go_utils` 并按包名 `utils` 引用。

## 快速开始

```go
package main

import (
    "fmt"
    utils "go_utils"
)

func main() {
    fmt.Println(utils.IsEmpty("  "))            // true
    fmt.Println(utils.FormatNow("2006-01-02"))  // 当天日期
    fmt.Println(utils.UUIDv4())                  // 随机UUID
    fmt.Println(utils.IsEmail("user@example.com")) // true
}
```

## 模块结构

- `strutil.go`：字符串处理
- `time_env_util.go`：时间与环境变量
- `file_json_util.go`：文件与 JSON
- `crypto_rand_util.go`：摘要/签名与随机工具
- `generic_collections_util.go`：集合泛型（切片/映射）
- `net_util.go`：本机 IPv4 获取
- `validate_util.go`：常见格式校验

## API 总览（部分）

- 字符串：`IsEmpty`、`Trim`、`ToUpper`、`ToLower`、`ContainsSubstr`、`ReplaceAll`、`Split`、`Join`、`Substring`（UTF-8安全）、`Reverse`、`PadLeft`、`PadRight`
- 时间与环境：`NowUnix`、`NowUnixMilli`、`FormatNow`、`FormatTime`、`ParseTime`、`AddDays`、`StartOfDay`、`EndOfDay`、`GetEnv`、`GetEnvDefault`
- 文件与 JSON：`Exists`、`IsDir`、`ReadFile`、`WriteFile`、`CopyFile`、`ListDir`、`ToJSON`、`ToPrettyJSON`、`FromJSON[T]`
- 摘要与随机：`MD5String`（非安全场景）、`SHA256String`、`SHA512String`、`HmacSHA256`、`RandomInt`、`RandomString`、`UUIDv4`
- 集合泛型：`Contains[T]`、`IndexOf[T]`、`Unique[T]`、`Map[T,R]`、`Filter[T]`、`Keys[K,V]`、`Values[K,V]`、`Merge[K,V]`、`GetOrDefault[K,V]`
- 网络：`GetLocalIPv4`
- 校验：`IsEmail`、`IsMobileCN`、`IsURL`、`IsIP`、`IsChineseIDCard`、`IsCNPatentApplicationNo`、`IsCNPatentNo`、`IsUnifiedSocialCreditCode`

## 常用示例

- 字符串

```go
utils.ContainsSubstr("hello", "he")      // true
utils.Substring("中文ABC", 0, 2)            // "中文"
utils.PadRight("ID", "0", 5)             // "ID000"
```

- 时间与环境

```go
utils.NowUnix()                            // 当前秒级时间戳
utils.FormatTime(time.Now(), "15:04:05")  // 当前时分秒
utils.GetEnvDefault("PORT", "8080")       // 环境变量或默认值
```

- 文件与 JSON

```go
_ = utils.WriteFile("out.txt", "content")
txt, _ := utils.ReadFile("out.txt")
js, _ := utils.ToPrettyJSON(map[string]int{"a":1})
type User struct { Name string }
u, _ := utils.FromJSON[User]("{\"Name\":\"Bob\"}")
```

- 摘要与随机

```go
utils.SHA256String("data")                // 64位十六进制字符串
utils.HmacSHA256("data", "secret")        // HMAC-SHA256
utils.RandomString(16)                     // 16位安全随机字符串
utils.UUIDv4()                             // 随机UUID v4
```

- 集合泛型

```go
utils.Contains([]int{1,2,3}, 2)            // true
utils.Unique([]string{"a","a","b"})     // ["a","b"]
utils.Map([]int{1,2,3}, func(x int) string { return fmt.Sprint(x) })
```

- 网络

```go
ips, _ := utils.GetLocalIPv4()             // 非回环IPv4列表
```

- 校验

```go
utils.IsEmail("user@example.com")         // true
utils.IsMobileCN("+86 13800138000")       // true
utils.IsURL("example.com")                // true（自动补全http）
utils.IsIP("2001:db8::1")                 // true（IPv6）
utils.IsChineseIDCard("11010519491231002X") // true/false 视具体号码
utils.IsCNPatentApplicationNo("201410123456.7") // true（常见格式）
utils.IsCNPatentNo("CN102345678A")        // true（常见格式）
utils.IsUnifiedSocialCreditCode("91350100M000100Y4") // true/false
```

## 注意事项

- `MD5String` 不适合安全场景（请使用 `SHA256String`/`HmacSHA256`）。
- URL 校验为通用场景设计，若存在企业内网域名或无后缀的主机名，请根据需要扩展。
- 中国身份证校验包含日期与校验位算法，不校验行政区划码字典。
- 中国专利申请号/专利号格式较多，此实现覆盖常见形式，如需更严谨校验，可提供样例加强规则。

## 构建

```bash
go build ./...
```

## 贡献与扩展

欢迎根据你的场景补充更多静态方法（如日志、正则校验、配置加载、HTTP 请求封装等）。保持中文注释与关键步骤说明的一致性即可。