# go_utils

实用工具集合，按领域拆分多个包：中文数字转换、类型转换、字符串、文件与JSON、网络、校验、时间环境、加密随机、泛型集合等。

## 安装

```go
import (
    "go_utils/numberchinese"
    "go_utils/convert"
    // 其他包按需引用：
    // "go_utils/strutil"
    // "go_utils/filejson"
    // "go_utils/netutil"
    // "go_utils/validate"
    // "go_utils/timeenv"
    // "go_utils/cryptorand"
    // "go_utils/collections"
)
```

## 中文数字转换（Number → 中文）

提供整数与浮点的中文小写/大写读法，以及通用接口。

API：
- `ToChineseLowerInt(n int64) string`
- `ToChineseUpperInt(n int64) string`
- `ToChineseLowerFloat(f float64, decimalPlaces int) string`
- `ToChineseUpperFloat(f float64, decimalPlaces int) string`
- `ToChineseLowerNumber(v interface{}, decimalPlaces int) (string, error)`
- `ToChineseUpperNumber(v interface{}, decimalPlaces int) (string, error)`

示例：

```go
// 整数
numberchinese.ToChineseLowerInt(12345) // "一万二千三百四十五"
numberchinese.ToChineseUpperInt(10010) // "壹万零壹拾"

// 浮点（保留两位小数）
numberchinese.ToChineseLowerFloat(10.25, 2) // "十点二五"
numberchinese.ToChineseUpperFloat(1001.05, 2) // "壹仟零壹点零五"

// 通用，自动分派整数/浮点
s, err := numberchinese.ToChineseUpperNumber(123.456, 2) // "壹佰贰拾叁点四六"
```

说明：
- 四位一组（万、亿、兆），自动处理连续零。
- 10～19 的中文小写常用省略“一”，大写保留“壹”。
- 浮点小数位按 `decimalPlaces` 四舍五入。

## 人民币中文大写（元/角/分）

将金额转换为人民币中文大写，支持负数与“整”。

API：
- `ToChineseRMBUpper(amount float64) string`
- `ToChineseRMBUpperNumber(v interface{}) (string, error)`

示例：

```go
numberchinese.ToChineseRMBUpper(123.45)     // "壹佰贰拾叁元肆角伍分"
numberchinese.ToChineseRMBUpper(100.0)      // "壹佰元整"
numberchinese.ToChineseRMBUpper(-10.01)     // "负壹拾元零壹分"
numberchinese.ToChineseRMBUpperNumber("0") // "零元整"
```

说明：
- 元为整数部分，角/分为小数两位；第三位及以上按分位四舍五入。
- 整数部分为 0 时仍输出“零元”。小数为 0 时输出“整”。
- 负值前缀“负”。

## 类型转换工具（统一错误处理）

提供常见数据类型之间的转换方法，所有失败情况统一返回 `ConvertError`，包含来源类型、目标类型、原始值与错误说明。

API：
- `ToString(v interface{}) (string, error)`
- `ToInt(v interface{}) (int, error)`
- `ToInt64(v interface{}) (int64, error)`
- `ToUint64(v interface{}) (uint64, error)`
- `ToFloat64(v interface{}) (float64, error)`
- `ToBool(v interface{}) (bool, error)`

示例：

```go
// 字符串 → 数字
i, _ := convert.ToInt("123")         // 123
i64, _ := convert.ToInt64("-42")     // -42
u64, _ := convert.ToUint64("99")     // 99
f64, _ := convert.ToFloat64("3.14")  // 3.14

// 数字/布尔 → 字符串
s1, _ := convert.ToString(100)        // "100"
s2, _ := convert.ToString(true)       // "true"

// 任意 → 布尔
b1, _ := convert.ToBool("yes")        // true
b2, _ := convert.ToBool(0)            // false

// 错误处理统一：
if _, err := convert.ToUint64("-1"); err != nil {
    // err 为 *convert.ConvertError，包含 FromType/ToType/Value/Message
    fmt.Println(err)
}
```

说明：
- 字符串解析遵循 Go 的 `strconv` 规则，空字符串视为错误。
- `ToUint64` 不接受负数（包含负整型与负浮点）。
- 浮点转换采用截断取整原则，例如 `ToInt(3.9) == 3`。
- `ToBool` 支持常见真值/假值词：真值如 `"true"/"yes"/"on"/"1"/"是"/"真"/"开"`；假值如 `"false"/"no"/"off"/"0"/"否"/"假"/"关"`。
- 所有失败都返回 `ConvertError`，包含 `FromType/ToType/Value/Message`，便于统一日志和提示。

## 进制转换

提供二进制/八进制/十六进制的格式化、解析与互转。

API：
- `convert.ToBinString(v interface{}) (string, error)`
- `convert.ToBinStringWithPrefix(v interface{}) (string, error)`
- `convert.ToOctString(v interface{}) (string, error)`
- `convert.ToOctStringWithPrefix(v interface{}) (string, error)`
- `convert.ToHexString(v interface{}, uppercase bool) (string, error)`
- `convert.ToHexStringWithPrefix(v interface{}, uppercase bool) (string, error)`
- `convert.ParseIntFromBase(s string, base int) (int64, error)`
- `convert.ParseUintFromBase(s string, base int) (uint64, error)`
- `convert.ConvertBaseString(s string, fromBase, toBase int, uppercase, withPrefix bool) (string, error)`

示例：

```go
// 数值转进制字符串
convert.ToBinString(10)                  // "1010"
convert.ToBinStringWithPrefix(-10)       // "-0b1010"
convert.ToOctString(64)                  // "100"
convert.ToHexString(255, true)           // "FF"
convert.ToHexStringWithPrefix(255, false)// "0xff"

// 解析与互转
convert.ParseIntFromBase("0b1010", 0)   // 10, nil（自动识别前缀）
convert.ParseUintFromBase("FF", 16)     // 255, nil
convert.ConvertBaseString("-255", 10, 16, false, true) // "-0xff"
```

说明：
- `ParseIntFromBase`/`ParseUintFromBase` 支持 `base=0` 自动检测（`0x/0o/0b`）。
- 负数格式化保留 `-` 前缀；无符号解析不支持负号。
- 十六进制可选择大小写；`WithPrefix` 函数添加 `0x/0o/0b` 前缀。

## 许可证

MIT
## 数组/切片 ↔ map 转换

提供数组/切片与 `map` 的互转方法（使用反射实现，顺序不保证）。

API：
- `convert.ToMapFromSlice(v interface{}) (map[int]interface{}, error)`：数组/切片转 `map`，键为索引。
- `convert.ToSliceFromMapValues(v interface{}) ([]interface{}, error)`：`map` 的值转切片。
- `convert.ToSliceFromMapKeys(v interface{}) ([]interface{}, error)`：`map` 的键转切片。

示例：

```go
arr := []int{10, 20, 30}
m, _ := convert.ToMapFromSlice(arr)               // map[int]interface{}{0:10,1:20,2:30}
vals, _ := convert.ToSliceFromMapValues(m)        // [10 20 30]（顺序不保证）
keys, _ := convert.ToSliceFromMapKeys(map[string]int{"a":1,"b":2}) // ["a" "b"]（顺序不保证）
```

说明：
- 仅支持数组或切片与 `map` 类型；空输入返回错误；类型不匹配返回 `ConvertError`。
- `map` 的遍历顺序不保证，生成的切片顺序可能不同。

## 结构体 ↔ map 转换

支持结构体与 `map[string]interface{}` 的互转，可指定标签名（如 `json`）。指针与嵌套结构体会递归处理。

API：
- `convert.ToMapFromStruct(v interface{}, tag string) (map[string]interface{}, error)`：结构体/结构体指针转 `map`。
- `convert.FillStructFromMap(m map[string]interface{}, out interface{}, tag string) error`：`map` 填充到结构体指针。

示例：

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`
    Age  uint8  `json:"age"`
}

u := User{ID: 1001, Name: "Alice", Age: 20}
m, _ := convert.ToMapFromStruct(u, "json") // map[string]interface{"id":1001,"name":"Alice","age":20}

var u2 User
_ = convert.FillStructFromMap(m, &u2, "json") // u2 恢复为原值

// 嵌套结构体/指针会递归处理
type Profile struct {
    User *User `json:"user"`
}
pf := Profile{User: &u}
pm, _ := convert.ToMapFromStruct(pf, "json") // map{"user": map{"id":1001,"name":"Alice","age":20}}

var pf2 Profile
_ = convert.FillStructFromMap(pm, &pf2, "json")
```

说明：
- 标签为 `"-"` 的字段会被跳过；空标签或无标签时使用字段名。
- 赋值会进行基础类型转换与范围校验（如 `int16/int32/uint8/uint16/uint32`）。
- 未导出字段不会被读写；`map` 遍历顺序不保证。